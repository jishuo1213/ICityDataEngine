package repo

import (
	"IcityMessageBus/cmsp"
	"database/sql"
	"IcityMessageBus/utils"
	"ICityDataEngine/model"
	bModel "IcityMessageBus/model"
	"ICityDataEngine/constant"
	"errors"
	"net/url"
	"bytes"
	"mime/multipart"
	"encoding/json"
	"strconv"
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
	for i := 0; i < lenCN; i++ {
		s.cp[i] = new(sql.RawBytes)
	}
	return s
}

func (s *mapStringScan) Update(rows *sql.Rows) error {
	if err := rows.Scan(s.cp...); err != nil {
		return err
	}

	for i := 0; i < s.colCount; i++ {
		if rb, ok := s.cp[i].(*sql.RawBytes); ok {
			s.row[s.colNames[i]] = string(*rb)
			*rb = nil // reset pointer to discard current value to avoid a bug
		} else {
			return errors.New("Cannot convert column " + s.colNames[i] + " to type *sql.RawBytes")
		}
	}
	return nil
}

func (s *mapStringScan) Get() map[string]string {
	return s.row
}

func GenerateRequest(config *model.HttpRequestConfig) error {
	//http.Request{}
	//switch expr {
	//
	//}

	err := cmsp.DeleteQueueNet(constant.CMSP_IP, constant.CMSP_PORT, config.Id)
	defer cmsp.DisconnectCmsp(constant.CMSP_IP, constant.CMSP_PORT, config.Id)
	if err != nil {
		return errors.New("连接cmsp失败:" + err.Error())
	}

	var valueHeaders map[string]string
	var headerDbMapping map[string]string

	if len(config.Headers) > 0 {
		valueHeaders, headerDbMapping = separateValues(config.Headers)
	}

	var valueParams map[string]string
	var paramDbMapping map[string]string

	if len(config.Params) > 0 {
		valueParams, paramDbMapping = separateValues(config.Params)
	}

	//if config.SqlConfig != nil {

	err = QuerySqlParams(func(rows *sql.Rows) error {
		//model.RequestInfo{}
		var requestHeaders map[string]string
		var requestParams map[string]string

		if valueHeaders != nil {
			requestHeaders = make(map[string]string)
			for key, value := range valueHeaders {
				requestHeaders[key] = value
			}
		}

		if valueParams != nil {
			requestParams = make(map[string]string)
			for key, value := range valueParams {
				requestParams[key] = value
			}
		}
		if rows != nil {
			columns, err := rows.Columns()
			if err != nil {
				return errors.New("get sql query result columns failed")
			}
			mapScan := newMapStringScan(columns)
			multipartWriters := make([]*multipart.Writer, 0, 10)
			defer func() {
				for _, w := range multipartWriters {
					w.Close()
				}
			}()
			for i := 0; rows.Next(); i++ {
				//rows[0].Scan()
				err := mapScan.Update(rows)
				if err != nil {
					return err
				}
				result := mapScan.Get()

				if headerDbMapping != nil {
					for name, dbName := range headerDbMapping {
						requestHeaders[name] = result[dbName]
					}
				}

				if paramDbMapping != nil {
					for name, dbName := range headerDbMapping {
						requestParams[name] = result[dbName]
					}
				}

				var body *bytes.Buffer
				var getUrl string

				if config.Method == constant.GET {
					data := url.Values{}
					if valueParams != nil {
						for key, value := range requestParams {
							data.Add(key, value)
						}
					}
					getUrl = "?" + data.Encode()
					config.Url += getUrl
				} else if config.Method == constant.POST {
					switch config.ContentType {
					case constant.BODY_XFORM_TYPE:
						data := url.Values{}
						if valueParams != nil {
							for key, value := range requestParams {
								data.Add(key, value)
							}
						}
						body = bytes.NewBufferString(data.Encode())
						break
					case constant.BODY_FORM_TYPE:
						if valueParams != nil {
							body = bytes.NewBuffer(make([]byte, 0, 1024))
							w := multipart.NewWriter(body)
							boundary := "--" + config.Id + strconv.Itoa(i)
							w.SetBoundary(boundary)
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
					case constant.BODY_JSON_TYPE:
						if valueParams != nil {
							valueBytes, err := json.Marshal(valueParams)
							if err != nil {
								return errors.New("生成请求失败3")
							}
							body = bytes.NewBuffer(valueBytes)
						}
						break
					}
				}

			}
		} else {
			requestInfo := bModel.RequestInfo{Method: getMethod(config.Method), Url: config.Url, Headers: requestHeaders, Body: body.String(), Id: ""}
			requestBytes, err := utils.EncodeObject(requestInfo)
			if err != nil {
				return errors.New("序列化请求失败")
			}
			err = cmsp.PutMsgIntoQueueNet(constant.CMSP_IP, constant.CMSP_PORT, config.Id, requestBytes)
			if err != nil {
				return errors.New("请求入队失败" + err.Error())
			}
		}
		return nil
	}, config.SqlConfig)
	if err != nil {
		return err
	}
	return nil
}

func getMethod(method constant.HttpMethod) string {
	switch method {
	case constant.GET:
		return "GET"
	case constant.POST:
		return "POST"
	}
	return ""
}

func separateValues(params []*model.HttpVariables) (map[string]string, map[string]string) {
	values := make(map[string]string)
	dbMapping := make(map[string]string)
	for _, header := range params {
		if header.Type == constant.VALUE {
			values[header.Name] = header.Value
		} else if header.Type == constant.DB {
			dbMapping[header.Name] = header.DBMapping
		}
	}
	return values, dbMapping
}
