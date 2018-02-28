package model

import (
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/constant"
	"github.com/kataras/iris/core/errors"
	"encoding/json"
)

type HttpVariables struct {
	Name      string                 `bson:"name"`
	Type      constant.HttpParamFrom `bson:"type"` //类型：包含两大类，固定值和非固定值 固定值0 数据库1 文件2
	Value     string                 `bson:"value"`
	DBMapping string                 `bson:"mapping_name"`
	DataType  string                 `bson:"data_type"` //数据在json中的类型，有 int double string boolean四种 默认string
}

type HttpRequestConfig struct {
	Url         string               `bson:"url"`
	Method      constant.HttpMethod  `bson:"method"`
	ContentType constant.ContentType `bson:"content_type"`
	Headers     []*HttpVariables     `bson:"headers"`
	Params      []*HttpVariables     `bson:"params"`
	SqlConfig   SqlParamConfig       `bson:"sql_config"`
	Id          string               `bson:"-"`
}

func NewHttpRequestConfig(requestConfig *simplejson.Json, id string) (*HttpRequestConfig, error) {
	var config = &HttpRequestConfig{Id: id}

	requestUrl, err := requestConfig.Get("url").String()
	if err != nil {
		return nil, errors.New("request_config中url不存在或类型错误")
	}
	config.Url = requestUrl

	method, err := requestConfig.Get("method").Int()
	if err != nil {
		return nil, errors.New("request_config中method不存在或类型错误")
	}
	config.Method = method

	switch method {
	case constant.POST:
		bodyType, err := requestConfig.Get("body_type").Int()
		if err != nil {
			return nil, errors.New("request_config中body_type不存在或类型错误")
		}
		config.ContentType = bodyType
		switch bodyType {
		case constant.BODY_XFORM_TYPE:
			config.Headers = append(config.Headers, &HttpVariables{"Content-Type", constant.VALUE,
				"application/x-www-form-urlencoded", "", ""})
			break
		case constant.BODY_JSON_TYPE:
			config.Headers = append(config.Headers, &HttpVariables{"Content-Type", constant.VALUE,
				"application/json", "", ""})
			break
		case constant.BODY_FORM_TYPE:
			//----WebKitFormBoundary7MA4YWxkTrZu0gW
			//config.Headers = append(config.Headers, &HttpVariables{"Content-Type", constant.VALUE,
			//	"multipart/form-data; boundary=----" + config.Id, "", ""})
			break
		default:
			return nil, errors.New("不支持的请求体类型")
		}
		break
	case constant.GET:
		break
	default:
		return nil, errors.New("不支持的http的method类型")
	}

	paramsConfig, err := requestConfig.Get("variables_config").Map()

	if err == nil {
		sqlConfig, isDbConfig := paramsConfig["db_config"].(map[string]interface{})
		if params, ok := paramsConfig["variables"]; ok {
			paramsArray, success := params.([]interface{})
			if !success {
				return nil, errors.New("http请求参数设置json格式错误1")
			}
			for _, paramTemp := range paramsArray {
				param, success := paramTemp.(map[string]interface{})
				if !success {
					return nil, errors.New("http请求参数设置json格式错误2")
				}
				var httpParam HttpVariables
				paramName, success := param["name"].(string)
				if !success {
					return nil, errors.New("name不存在或格式错误,应该为string")
				}
				httpParam.Name = paramName

				//log.Println(reflect.TypeOf(param["from"]))
				paramFrom, err := param["from"].(json.Number).Int64()
				if err != nil {
					return nil, errors.New("from不存在或格式错误,应该为int")
				}
				if paramFrom == constant.DB && !isDbConfig {
					return nil, errors.New("参数来源数据库但未设置数据源")
				}
				httpParam.Type = int(paramFrom)
				switch paramFrom {
				case constant.DB:
					paramMappingName, success := param["mapping_name"].(string)
					if !success {
						return nil, errors.New("mapping_name不存在或格式错误,应该为string")
					}
					httpParam.DBMapping = paramMappingName
					break
				case constant.VALUE:
					paramValue, success := param["value"].(string)
					if !success {
						return nil, errors.New("value不存在或格式错误,应该为string")
					}
					httpParam.Value = paramValue
					break
				default:
					return nil, errors.New("参数来源" + paramName + "未知")
				}

				if config.ContentType == constant.BODY_JSON_TYPE {
					dataType, success := param["data_type"].(string)
					if success {
						httpParam.DataType = dataType
					} else {
						httpParam.DataType = "string"
					}
				}

				paramType, err := param["type"].(json.Number).Int64()
				if err != nil {
					return nil, errors.New("name不存在或格式错误,应该为int")
				}
				switch paramType {
				case constant.BODY:
					config.Params = append(config.Params, &httpParam)
					break
				case constant.HEADER:
					config.Headers = append(config.Headers, &httpParam)
					break
				default:
					return nil, errors.New("请问除了body(1)和header(2)你还想在哪传参数？")
					break
				}
			}
		}
		//config.SqlConfig = sqlConfig
		if isDbConfig {
			dbType := sqlConfig["db_type"].(string)
			switch dbType {
			case "mysql":
				userName, success := sqlConfig["user_name"].(string)
				if !success {
					return nil, errors.New("user_name不存在或格式错误,应该为string")
				}
				password, success := sqlConfig["password"].(string)
				if !success {
					return nil, errors.New("password不存在或格式错误,应该为string")
				}
				dbIP, success := sqlConfig["db_ip"].(string)
				if !success {
					return nil, errors.New("db_ip不存在或格式错误,应该为string")
				}
				dbPort, err := sqlConfig["port"].(json.Number).Int64()
				if err != nil {
					return nil, errors.New("port不存在或格式错误,应该为int")
				}
				dbName, success := sqlConfig["db_name"].(string)
				if !success {
					return nil, errors.New("db_name不存在或格式错误,应该为string")
				}
				sqlSentence, success := sqlConfig["sql"].(string)
				if !success {
					return nil, errors.New("sql不存在或格式错误,应该为string")
				}
				config.SqlConfig = &MySqlConfig{userName, password, dbIP, int(dbPort), dbName, sqlSentence}
				break
			default:
				return nil, errors.New("不支持的数据库类型")
			}
		} else {
			config.SqlConfig = nil
		}
	}

	return config, nil
}
