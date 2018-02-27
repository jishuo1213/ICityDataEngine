package repo

import (
	"IcityMessageBus/cmsp"
	"database/sql"
	"bytes"
	"net/url"
	"mime/multipart"
	"encoding/json"
	"IcityMessageBus/utils"
	"ICityDataEngine/model"
	bModel "IcityMessageBus/model"
	"ICityDataEngine/constant"
	"errors"
)

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
	//var headerDbMapping map[string]string

	if len(config.Headers) > 0 {
		valueHeaders, _ = separateValues(config.Headers)
	}

	var valueParams map[string]string
	//var paramDbMapping map[string]string

	if len(config.Params) > 0 {
		valueParams, _ = separateValues(config.Params)
	}

	//if config.SqlConfig != nil {

	err = QuerySqlParams(func(rows *sql.Rows) error {
		//model.RequestInfo{}
		var requestHeaders map[string]string

		if valueHeaders != nil {
			requestHeaders = make(map[string]string)
			for key, value := range valueHeaders {
				requestHeaders[key] = value
			}
		}

		var body *bytes.Buffer
		var getUrl string

		if config.Method == constant.GET {
			data := url.Values{}
			if valueParams != nil {
				for key, value := range valueParams {
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
					for key, value := range valueParams {
						data.Add(key, value)
					}
				}
				body = bytes.NewBufferString(data.Encode())
				break
			case constant.BODY_FORM_TYPE:
				if valueParams != nil {
					body = bytes.NewBuffer(make([]byte, 0, 1024))
					w := multipart.NewWriter(body)
					defer w.Close()
					for key, value := range valueParams {
						fw, err := w.CreateFormField(key)
						if err != nil {
							return errors.New("生成请求失败1")
						}
						_, err = fw.Write([]byte(value))
						if err != nil {
							return errors.New("生成请求失败2")
						}
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

		if rows != nil {
			for rows.Next() {
				//rows[0].Scan()

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
