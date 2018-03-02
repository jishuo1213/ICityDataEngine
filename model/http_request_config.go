package model

import (
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/constant"
	"github.com/kataras/iris/core/errors"
	"encoding/json"
	"strconv"
	"log"
	"ICityDataEngine/i"
)

type httpVariables struct {
	Name      string                 `bson:"name"`
	Type      constant.HttpParamFrom `bson:"type"` //类型：包含两大类，固定值和非固定值 固定值0 数据库1 文件2
	Value     string                 `bson:"value"`
	DBMapping string                 `bson:"mapping_name"`
	DataType  string                 `bson:"data_type"` //数据在json中的类型，有 int double string boolean四种 默认string
}

type httpRequestConfig struct {
	Url         string               `bson:"url"`
	Method      constant.HttpMethod  `bson:"method"`
	ContentType constant.ContentType `bson:"content_type"`
	Headers     []*httpVariables     `bson:"headers"`
	Params      []*httpVariables     `bson:"params"`
	SqlConfig   i.ISqlParamConfig    `bson:"sql_config"`
	SqlRepo     i.ISqlParamRepo      `bson:"-"`
	Id          string               `bson:"-"`
}

func (config *httpRequestConfig) GetId() string {
	return config.Id
}

func (config *httpRequestConfig) GetMethod() constant.HttpMethod {
	return config.Method
}

func (config *httpRequestConfig) GetContentType() constant.ContentType {
	return config.ContentType
}
func (config *httpRequestConfig) GetUrl() string {
	return config.Url
}
func (config *httpRequestConfig) GetSqlConfig() i.ISqlParamConfig {
	return config.SqlConfig
}

func (config *httpRequestConfig) InitValueHeaders() (map[string]string, map[string]string) {
	return initHttpVariables(config.Headers)
}

func (config *httpRequestConfig) InitValueParams() (map[string]string, map[string]string) {
	return initHttpVariables(config.Params)
}

func (config *httpRequestConfig) GenerateJsonBody(values map[string]string) (map[string]interface{}, error) {
	if config.ContentType != constant.BodyJsonType {
		return nil, errors.New("别乱调方法")
	}
	res := make(map[string]interface{})
	for _, variable := range config.Params {
		if variable.Type != constant.DB {
			continue
		}
		if value, ok := values[variable.Name]; ok {
			switch variable.DataType {
			case "int":
				v, err := strconv.Atoi(value)
				if err != nil {
					return nil, errors.New(variable.Name + "类型转换失败," + value + err.Error())
				}
				res[variable.Name] = v
				break
			case "double":
				v, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, errors.New(variable.Name + "类型转换失败," + value + err.Error())
				}
				res[variable.Name] = v
				break
			case "boolean":
				v, err := strconv.ParseBool(value)
				if err != nil {
					return nil, errors.New(variable.Name + "类型转换失败," + value + err.Error())
				}
				res[variable.Name] = v
				break
			case "string":
				res[variable.Name] = value
				break
			}
		} else {
			return nil, errors.New(variable.Name + "在数据库中没查出来")
		}

	}
	return res, nil
}

func initHttpVariables(initArg []*httpVariables) (map[string]string, map[string]string) {
	if len(initArg) > 0 {
		if valueHeaders, DbMapping := separateValues(initArg); valueHeaders != nil {
			requestHeaders := make(map[string]string)
			for key, value := range valueHeaders {
				requestHeaders[key] = value
			}
			return requestHeaders, DbMapping
		} else {
			return nil, DbMapping
		}
	}
	return nil, nil
}

func separateValues(params []*httpVariables) (map[string]string, map[string]string) {
	values := make(map[string]string)
	dbMapping := make(map[string]string)
	for _, header := range params {
		if header.Type == constant.Value {
			values[header.Name] = header.Value
		} else if header.Type == constant.DB {
			dbMapping[header.Name] = header.DBMapping
		}
	}
	if len(values) == 0 {
		values = nil
	}
	if len(dbMapping) == 0 {
		dbMapping = nil
	}

	return values, dbMapping
}

func NewHttpRequestConfig(requestConfig *simplejson.Json, id string) (*httpRequestConfig, error) {
	var config = &httpRequestConfig{Id: id}

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
	case constant.Post:
		bodyType, err := requestConfig.Get("body_type").Int()
		if err != nil {
			return nil, errors.New("request_config中body_type不存在或类型错误")
		}
		config.ContentType = bodyType
		switch bodyType {
		case constant.BodyXFormType:
			config.Headers = append(config.Headers, &httpVariables{"Content-Type", constant.Value,
				"application/x-www-form-urlencoded", "", ""})
			break
		case constant.BodyJsonType:
			config.Headers = append(config.Headers, &httpVariables{"Content-Type", constant.Value,
				"application/json", "", ""})
			break
		case constant.BodyFormType:
			//----WebKitFormBoundary7MA4YWxkTrZu0gW
			//config.Headers = append(config.Headers, &httpVariables{"Content-Type", constant.Value,
			//	"multipart/form-data; boundary=----" + config.Id, "", ""})
			break
		default:
			return nil, errors.New("不支持的请求体类型")
		}
		break
	case constant.Get:
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
				var httpParam httpVariables
				paramName, success := param["name"].(string)
				if !success {
					return nil, errors.New("name不存在或格式错误,应该为string")
				}
				httpParam.Name = paramName

				//log.Println(reflect.TypeOf(param["from"]))
				paramFrom, err := param["data_from"].(json.Number).Int64()
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
				case constant.Value:
					paramValue, success := param["value"].(string)
					if !success {
						return nil, errors.New("value不存在或格式错误,应该为string")
					}
					httpParam.Value = paramValue
					break
				default:
					return nil, errors.New("参数来源" + paramName + "未知")
				}

				if config.ContentType == constant.BodyJsonType {
					dataType, success := param["data_type"].(string)
					if success {
						httpParam.DataType = dataType
					} else {
						httpParam.DataType = "string"
					}
				}

				paramType, err := param["data_to"].(json.Number).Int64()
				if err != nil {
					return nil, errors.New("name不存在或格式错误,应该为int")
				}
				log.Println(httpParam.Name+":", paramType)
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
				//config.SqlRepo = &repo.QueryMySqlParamsRepo{}

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
