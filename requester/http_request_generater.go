package requester

import (
	"IcityMessageBus/cmsp"
	"database/sql"
	"IcityMessageBus/utils"
	bModel "IcityMessageBus/model"
	"ICityDataEngine/constant"
	"errors"
	"net/url"
	"bytes"
	"mime/multipart"
	"encoding/json"
	"strconv"
	"ICityDataEngine/i"
)

type mapStringScan struct {
	// cp are the column pointers
	cp []interface{}
	// row contains the final result
	row      map[string]string
	colCount int
	colNames []string
}

func newMapStringScan(columnNames []string) *mapStringScan {
	lenCN := len(columnNames)
	s := &mapStringScan{
		cp:       make([]interface{}, lenCN),
		row:      make(map[string]string, lenCN),
		colCount: lenCN,
		colNames: columnNames,
	}
	for index := 0; index < lenCN; index++ {
		s.cp[index] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for index := 0; index < s.colCount; index++ {
		if rb, ok := s.cp[index].(*sql.RawBytes); ok {
			s.row[s.colNames[index]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return errors.New("Cannot convert column " + s.colNames[index] + " to type *sql.RawBytes")
		}
	}
	return nil
}

func (s *mapStringScan) Get() map[string]string {
	return s.row
}

type sqlResultParser struct {
	requestConfig i.IRequestConfig
	queueId       string
}

func (parser *sqlResultParser) Parse(rows *sql.Rows) error {

	requestConfig := parser.requestConfig

	requestHeaders, headerDbMapping := requestConfig.InitValueHeaders()
	requestParams, paramDbMapping := requestConfig.InitValueParams()
	isNoRows := false

	if rows == nil {
		isNoRows = true
	}

	//if rows != nil {
	var mapScan *mapStringScan
	if rows != nil {
		columns, err := rows.Columns()
		if err != nil {
			return errors.New("get sql query result columns failed")
		}
		mapScan = newMapStringScan(columns)
	}
	multipartWriters := make([]*multipart.Writer, 0, 10)

	for count := 0; isNoRows || rows.Next(); count++ {
		//rows[0].Scan()
		if mapScan != nil {
			err := mapScan.Update(rows)
			if err != nil {
				return err
			}

			requestHeaders = addDBArguments(mapScan, headerDbMapping, requestHeaders)
			requestParams = addDBArguments(mapScan, paramDbMapping, requestParams)
		}

		var body *bytes.Buffer
		var requestUrl string

		if requestConfig.GetMethod() == constant.Get {
			data := url.Values{}
			if requestParams != nil {
				for key, value := range requestParams {
					data.Add(key, value)
				}
			}
			requestUrl = "?" + data.Encode()
			//config.Url += getUrl
			requestUrl = requestConfig.GetUrl() + requestUrl
			//body = bytes.NewBufferString("")
		} else if requestConfig.GetMethod() == constant.Post {
			requestUrl = requestConfig.GetUrl()
			switch requestConfig.GetContentType() {
			case constant.BodyXFormType:
				data := url.Values{}
				if requestParams != nil {
					for key, value := range requestParams {
						data.Add(key, value)
					}
				}
				body = bytes.NewBufferString(data.Encode())
				break
			case constant.BodyFormType:
				if requestParams != nil {
					body = bytes.NewBuffer(make([]byte, 0, 1024))
					w := multipart.NewWriter(body)
					boundary := "--" + requestConfig.GetId() + strconv.Itoa(count)
					w.SetBoundary(boundary)
					if requestHeaders == nil {
						requestHeaders = make(map[string]string)
					}
					requestHeaders["Content-Type"] = "multipart/form-data; boundary=" + boundary

					for key, value := range requestParams {
						fw, err := w.CreateFormField(key)
						if err != nil {
							return errors.New("生成请求失败1")
						}
						_, err = fw.Write([]byte(value))
						if err != nil {
							return errors.New("生成请求失败2")
						}
					}
					err := w.Close()
					if err != nil {
						multipartWriters = append(multipartWriters, w)
					}
				}
				break
			case constant.BodyJsonType:
				if requestParams != nil {

					jsonBody, err := requestConfig.GenerateJsonBody(requestParams)
					if err != nil {
						return err
					}
					valueBytes, err := json.Marshal(jsonBody)
					if err != nil {
						return errors.New("生成请求失败3")
					}
					body = bytes.NewBuffer(valueBytes)
				}
				break
			}
		}
		var bodyStr string
		if body != nil {
			bodyStr = body.String()
		}
		requestInfo := bModel.RequestInfo{Method: getMethod(requestConfig.GetMethod()), Url: requestUrl, Headers: requestHeaders, Body: bodyStr, Id: ""}
		requestBytes, err := utils.EncodeObject(requestInfo)
		if err != nil {
			return errors.New("序列化请求失败")
		}
		err = cmsp.PutMsgIntoQueueNet(constant.CMSPIP, constant.CMSPPort, parser.queueId, requestBytes)
		if err != nil {
			return errors.New("请求入队失败" + err.Error())
		}

		if rows == nil {
			break
		}
	}
	defer func() {
		for _, w := range multipartWriters {
			w.Close()
		}
	}()
	return nil
}

func GenerateRequest(wrapConfig i.IRequestConfig, queueId string) error {

	err := cmsp.DeleteQueueNet(constant.CMSPIP, constant.CMSPPort, queueId)
	defer cmsp.DisconnectCmsp(constant.CMSPIP, constant.CMSPPort, queueId)
	if err != nil {
		return errors.New("连接cmsp失败:" + err.Error())
	}

	parser := &sqlResultParser{requestConfig: wrapConfig, queueId: queueId}

	if wrapConfig.GetSqlConfig() == nil {
		err = parser.Parse(nil)
		return err
	}

	err = wrapConfig.GetSqlConfig().QuerySqlParams(parser)

	if err != nil {
		return err
	}
	return nil
}

func addDBArguments(mapScan *mapStringScan, DbMapping map[string]string, request map[string]string) (map[string]string) {
	result := mapScan.Get()

	if DbMapping != nil {
		if request == nil {
			request = make(map[string]string)
		}
		for name, dbName := range DbMapping {
			request[name] = result[dbName]
		}
		return request
	}
	return request
}

func getMethod(method constant.HttpMethod) string {
	switch method {
	case constant.Get:
		return "GET"
	case constant.Post:
		return "POST"
	}
	return ""
}
