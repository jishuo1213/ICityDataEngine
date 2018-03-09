package model

import (
	"ICityDataEngine/i"
	"github.com/bitly/go-simplejson"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
	"IcityMessageBus/cmsp"
	"ICityDataEngine/constant"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
	"IcityMessageBus/model"
	"reflect"
	"ICityDataEngine/requester"
)

type CmspSaveConfig struct {
	//SaveTo string `bson:"save_to"`
	TopicName string `json:"topic_name"`
}

func (config *CmspSaveConfig) Save(response []byte, info i.IRequestInfo) error {
	return cmsp.PutMsgIntoQueueNet(constant.CMSPIP, constant.CMSPPort, config.TopicName, response)
}

type SaveKey struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value []interface{}
}

type RedisSaveConfig struct {
	Ip          string     `json:"ip"`
	Port        int        `json:"port"`
	RedisDB     int        `json:"redis_db"`
	ExpiredTime string     `json:"expired_time"`
	SaveKey     []*SaveKey `json:"key"`
	//Request     i.IRequestInfo `json:"-"`
}

func (config *RedisSaveConfig) GetSaveKeys(request i.IRequestInfo) string {
	var key string
	for _, saveKey := range config.SaveKey {
		switch saveKey.Type {
		case "value":
			key += saveKey.Name
			break
		case "params":
			if value, ok := request.GetHeaders()[saveKey.Name]; ok {
				key += value
			} else if value, ok = request.GetBodyParams()[saveKey.Name]; ok {
				key += value
			}
			break
		}
	}
	return key
}

type RedisKeyEmptyError struct {
	error string
}

func (err RedisKeyEmptyError) Error() string {
	return err.error
}

type ConnectRedisError struct {
	error string
}

func (err ConnectRedisError) Error() string {
	return err.error
}

func (config *RedisSaveConfig) Save(response []byte, request i.IRequestInfo) error {
	key := config.GetSaveKeys(request)
	if len(key) == 0 {
		return RedisKeyEmptyError{"redis 存储key为空"}
	}

	options := make([]redis.DialOption, 0, 5)

	options = append(options, redis.DialConnectTimeout(30*time.Second))
	options = append(options, redis.DialReadTimeout(30*time.Second))
	options = append(options, redis.DialWriteTimeout(30*time.Second))
	if config.RedisDB >= 0 {
		options = append(options, redis.DialDatabase(config.RedisDB))
	}

	c, err := redis.Dial("tcp", config.Ip+":"+strconv.Itoa(config.Port), options ...)
	if err != nil {
		return ConnectRedisError{"连接redis失败" + err.Error()}
	}
	defer c.Close()
	n, err := c.Do("SET", key, response)
	if err != nil {
		return err
	}
	if n == 1 {
		if len(config.ExpiredTime) > 0 {
			n, err := c.Do("EXPIRE", key, config.ExpiredTime)
			if err != nil {
				return err
			}
			if n == 1 {
				return nil
			} else {
				return errors.New("设置redis过期时间失败")
			}
		} else {
			return nil
		}
	} else {
		return errors.New("保存到redis失败")
	}
}

type DBSaveConfig struct {
	SqlConfig      i.ISqlSaveConfig
	InsertTemplate string
	InsertMapping  []*SaveKey
	BodyType       string
	InsertOrder    []string
}

func (config *DBSaveConfig) Save(response []byte, request i.IRequestInfo) error {
	switch config.BodyType {
	case "jsonObject":
		responseMap := make(map[string]interface{})
		err := json.Unmarshal(response, &responseMap)
		if err != nil {
			return errors.New("返回结果格式不正确,应该是jsonObject")
		}
		templateMap := make(map[string]interface{})
		err = json.Unmarshal([]byte(config.InsertTemplate), &templateMap)
		if err != nil {
			return errors.New("template格式不正确,应该是jsonObject")
		}

		saveKeyMap := make(map[string]*SaveKey)
		for _, key := range config.InsertMapping {
			saveKeyMap[key.Name] = key
		}

		values, keys, err := getInsertValues(responseMap, templateMap, saveKeyMap)
		if err != nil {
			return err
		}

		if len(values)%len(config.InsertMapping) != 0 {
			return errors.New("结果中查询的值个数应该是数据库字段值的整数倍")
		}

		count := len(values) / len(config.InsertMapping)
		for index, key := range keys {
			saveKeyMap[key].Value = append(saveKeyMap[key].Value, values[index])
		}

		for index := 0; index < count; index++ {
			insertValues := make([]interface{}, 0, len(config.InsertOrder))
			for _, name := range config.InsertOrder {
				insertValues = append(insertValues, saveKeyMap[name].Value[index])
				err := config.SqlConfig.ExecInsert(insertValues)
				//if _, ok := err.(ConnectDBError); ok {
				//	return err
				//} else {
				//continue
				return err
				//}
			}
		}
		break
	case "xml":
		break
	}
	return nil
}

func getInsertValues(response map[string]interface{}, valueTemplate map[string]interface{}, saveKeyMap map[string]*SaveKey) ([]interface{}, []string, error) {
	res := make([]interface{}, 0, 10)
	keys := make([]string, 0, 10)
	for key, value := range valueTemplate {
		if saveKey, ok := saveKeyMap[key]; ok {
			switch saveKey.Type {
			case "string":
				value, ok := response[key].(string)
				if ok {
					res = append(res, value)
					keys = append(keys, key)
				} else {
					return nil, nil, errors.New(key + "type is not string")
				}
				break
			case "int":
				value, ok := response[key].(int64)
				if ok {
					res = append(res, value)
					keys = append(keys, key)
				} else {
					return nil, nil, errors.New(key + "type is not int")
				}
				break
			case "double":
				value, ok := response[key].(float64)
				if ok {
					res = append(res, value)
					keys = append(keys, key)
				} else {
					return nil, nil, errors.New(key + "type is not double")
				}
				break
			}
		} else {
			switch value.(type) {
			case map[string]interface{}:
				inResponse, ok := response[key].(map[string]interface{})
				if ok {
					insertValues, insertKeys, err := getInsertValues(inResponse, value.(map[string]interface{}), saveKeyMap)
					if err != nil {
						return nil, nil, err
					}
					res = append(res, insertValues...)
					keys = append(keys, insertKeys...)
				} else {
					return nil, nil, errors.New("取值模板和返回值类型不同:" + key)
				}
				break
			case []interface{}:
				inResponse, ok := response[key].([]interface{})
				if ok {
					insertValues, insertKeys, err := getInsertArrayValues(inResponse, value.([]interface{}), saveKeyMap)
					if err != nil {
						return nil, nil, err
					}
					res = append(res, insertValues...)
					keys = append(keys, insertKeys...)
				} else {
					return nil, nil, errors.New("取值模板和返回值类型不同:" + key)
				}
				break
			default:
				return nil, nil, errors.New("结果中找不到" + key)
			}
		}
	}
	return res, keys, nil
}

func getInsertArrayValues(response []interface{}, template []interface{}, saveKeyMap map[string]*SaveKey) ([]interface{}, []string, error) {
	res := make([]interface{}, 0, 10)
	keys := make([]string, 0, 10)
	switch template[0].(type) {
	case map[string]interface{}:
		for _, value := range response {
			mapValue, ok := value.(map[string]interface{})
			if ok {
				insertValues, insertKeys, err := getInsertValues(mapValue, template[0].(map[string]interface{}), saveKeyMap)
				if err != nil {
					return nil, nil, err
				}
				res = append(res, insertValues...)
				keys = append(keys, insertKeys...)
			} else {
				return nil, nil, errors.New("取值模板和返回值类型不同:")
			}
		}
		break
	case []interface{}:
		for _, value := range response {
			arrayValue, ok := value.([]interface{})
			if ok {
				insertValues, insertKeys, err := getInsertArrayValues(arrayValue, template[0].([]interface{}), saveKeyMap)
				if err != nil {
					return nil, nil, err
				}
				res = append(res, insertValues...)
				keys = append(keys, insertKeys...)
			} else {
				return nil, nil, errors.New("取值模板和返回值类型不同:")
			}
		}
		break
	default:
		return nil, nil, errors.New("取值模板和返回值类型不同:")
		//break
	}
	return res, keys, nil
}

type SuccessResponseConfig struct {
	Code              int    `bson:"success_response_code"`
	DataType          string `bson:"data_type"`
	SuccessTemplate   string `bson:"template"`
	AfterAction       []string
	SaveConfig        i.IResponseSaver
	SuccessHttpConfig i.IRequestInfo
}

func (config *SuccessResponseConfig) IsSuccessResponse(code int, response []byte) bool {

	switch config.DataType {
	case "jsonObject":

		template := make(map[string]interface{})
		err := json.Unmarshal([]byte(config.SuccessTemplate), &template)
		if err != nil {
			return false
		}
		body := make(map[string]interface{})
		err = json.Unmarshal(response, body)
		if err != nil {
			return false
		}
		return findSameInJsonBody(template, body, true)
	default:
		return false
	}
}

type SaveBodyError struct {
	info string
}

func (err SaveBodyError) Error() string {
	return err.info
}

type SendHttpError struct {
	info string
}

func (err SendHttpError) Error() string {
	return err.info
}

func (config *SuccessResponseConfig) DealSuccessResponse(request i.IRequestInfo, body []byte) error {
	for _, action := range config.AfterAction {
		switch action {
		case "http":
			request, err := config.SuccessHttpConfig.GenerateRequest()
			if err != nil {
				return SendHttpError{"生成请求失败"}
			}
			_, _, err = requester.SendHttpRequest(request)
			//return err
			if err != nil {
				return SendHttpError{"发送请求失败"}
			}
		case "save":
			err := config.SaveConfig.Save(body, request)
			if err != nil {
				return SaveBodyError{err.Error()}
			}
		}
	}
	return nil
}

func (config *IgnoreResponseConfig) IsIgnoreResponse(code int, response []byte) (bool) {

	switch config.DataType {
	case "jsonObject":
		body := make(map[string]interface{})
		err := json.Unmarshal(response, body)
		if err != nil {
			return false
		}
		for _, strTemplate := range config.Template {
			template := make(map[string]interface{})
			err = json.Unmarshal([]byte(strTemplate), &template)
			if err != nil {
				return false
			}
			if findSameInJsonBody(template, body, true) {
				return true
			}
		}
		return false
	default:
		return false
	}
}

func findSameInJsonBody(template map[string]interface{}, body map[string]interface{}, res bool) bool {
	//res := true
	if !res {
		return res
	}
	for key, value := range template {
		switch value.(type) {
		case map[string]interface{}:
			bodyValue, ok := body[key].(map[string]interface{})
			//if !ok {
			//	return false
			//}
			res = findSameInJsonBody(value.(map[string]interface{}), bodyValue, ok)
			//if !res {
			//	return res
			//}
			break
		case []interface{}:
			bodyValue, ok := body[key].([]interface{})
			//if !ok {
			//	return false
			//}
			res = findSameInArrayBody(value.([]interface{}), bodyValue, ok)
			//if !res {
			//	return res
			//}
			break
		default:
			res = res && reflect.TypeOf(value) == reflect.TypeOf(body[key]) && body[key] == value
			if !res {
				return res
			}
		}
	}
	return res
}

func findSameInArrayBody(template []interface{}, body []interface{}, res bool) bool {
	//res := true
	for index, v := range template {
		switch v.(type) {
		case []interface{}:
			bodyValue, ok := body[index].([]interface{})
			//if !ok {
			//	return false
			//}
			res = findSameInArrayBody(v.([]interface{}), bodyValue, ok)
			//if !res {
			//	return res
			//}
			break
		case map[string]interface{}:
			bodyValue, ok := body[index].(map[string]interface{})
			//if !ok {
			//	return false
			//}
			res = findSameInJsonBody(v.(map[string]interface{}), bodyValue, ok)
			//if !res {
			//	return res
			//}
			break
		default:
			res = reflect.TypeOf(v) == reflect.TypeOf(body[index]) && body[index] == v
			if !res {
				return res
			}
		}
	}
	return res
}

type FailedResponseConfig struct {
	Redo          bool `json:"redo"`
	MaxRetryTimes int  `json:"max_retry_times"`
	RetryInterval int  `json:"retry_interval"`
}

type IgnoreResponseConfig struct {
	ResponseCode int      `json:"response_code"`
	DataType     string   `json:"data_type"`
	Template     []string `json:"template"`
}

type ResponseConfig struct {
	*SuccessResponseConfig
	*FailedResponseConfig
	*IgnoreResponseConfig
}

//func (config *ResponseConfig) IsSuccessResponse(response []byte) (bool) {
//
//	return false
//}

func (config *ResponseConfig) DealFailedRequest(request i.IRequestInfo) error {
	return nil
}

func NewHttpResponseConfig(responseConfig *simplejson.Json, id string) (*ResponseConfig, error) {
	var successResConfig *SuccessResponseConfig
	var failedResConfig *FailedResponseConfig
	var ignoreResConfig *IgnoreResponseConfig
	if successResponseJson, ok := responseConfig.CheckGet("success_response_config"); ok {
		successResConfig = new(SuccessResponseConfig)
		code, err := successResponseJson.Get("response_code").Int()
		if err != nil {
			return nil, errors.New("success_response_config中response_code不存在或类型错误")
		}
		successResConfig.Code = code

		dataType, err := successResponseJson.Get("data_type").String()
		if err != nil {
			return nil, errors.New("success_response_config中data_type不存在或类型错误")
		}
		successResConfig.DataType = dataType

		successTemplate, err := successResponseJson.Get("template").String()
		if err != nil {
			return nil, errors.New("success_response_config中template不存在或类型错误")
		}
		successResConfig.SuccessTemplate = successTemplate

		if doAfter, ok := successResponseJson.CheckGet("do_after"); ok {
			doAfterArray, err := doAfter.StringArray()
			if err != nil {
				return nil, errors.New("success_response_config中do_after不存在或类型错误")
			}
			successResConfig.AfterAction = doAfterArray

			for _, action := range doAfterArray {
				switch action {
				case "save":
					if saveConfig, ok := successResponseJson.CheckGet("save_config"); ok {
						saveTo, err := saveConfig.Get("save_to").String()
						if err != nil {
							return nil, errors.New("save_config中save_to不存在或类型错误")
						}
						switch saveTo {
						case "redis":
							redisSaveConfig := RedisSaveConfig{}
							redisSaveJson, _ := saveConfig.MarshalJSON()
							err = json.Unmarshal(redisSaveJson, &redisSaveConfig)
							if err != nil {
								return nil, err
							}
							successResConfig.SaveConfig = &redisSaveConfig
							break
						case "db":
							dbType, err := saveConfig.Get("db_type").String()
							if err != nil {
								return nil, errors.New("save_config中db_type不存在或类型错误")
							}

							switch dbType {
							case "mysql":
								mySqlConfig := MySqlSaveConfig{}
								sqlSaveJson, _ := saveConfig.MarshalJSON()
								err = json.Unmarshal(sqlSaveJson, &mySqlConfig)
								if err != nil {
									return nil, err
								}
								dbSaveConfig := DBSaveConfig{}
								dbSaveConfig.BodyType = dataType

								dbSaveConfig.SqlConfig = &mySqlConfig

								dbSaveConfig.InsertTemplate, err = saveConfig.Get("insert_template").String()
								if err != nil {
									return nil, errors.New("save_config中insert_template不存在或类型错误")
								}

								dbSaveConfig.InsertOrder, err = saveConfig.Get("insert_order").StringArray()
								if err != nil {
									return nil, errors.New("save_config中insert_order不存在或类型错误")
								}

								mappingArray, err := saveConfig.Get("insert_mapping").MarshalJSON()
								if err != nil {
									return nil, errors.New("save_config中insert_mapping不存在或类型错误")
								}
								saveMapping := make([]*SaveKey, 0, 4)
								err = json.Unmarshal(mappingArray, &saveMapping)
								if err != nil {
									return nil, err
								}
								dbSaveConfig.InsertMapping = saveMapping

								successResConfig.SaveConfig = &dbSaveConfig
								//dbSaveConfig.InsertOrder
								break
							default:
								return nil, errors.New("不支持的数据库类型")
							}

							break
						case "cmsp":
							cmspSaveConfig := &CmspSaveConfig{}
							cmspSaveJson, _ := saveConfig.MarshalJSON()
							err = json.Unmarshal(cmspSaveJson, cmspSaveConfig)
							if err != nil {
								return nil, err
							}
							if len(cmspSaveConfig.TopicName) == 0 {
								return nil, errors.New("cmsp的topicname为空")
							}

							successResConfig.SaveConfig = cmspSaveConfig
							break
						default:
							return nil, errors.New("不知耻的存储类型")
						}
					} else {
						return nil, errors.New("请问你想把结果存哪去？在你深深的脑海里吗?")
					}
					break
				case "http":
					if httpConfig, ok := successResponseJson.CheckGet("http_config"); ok {
						requestType, err := httpConfig.Get("type").String()
						if err != nil {
							return nil, errors.New("http_config中的type不存在或类型错误")
						}
						switch requestType {
						case "http":
							url, err := httpConfig.Get("url").String()
							if err != nil {
								return nil, errors.New("http_config中的url不存在或类型错误")
							}
							method, err := httpConfig.Get("method").String()
							if err != nil {
								return nil, errors.New("http_config中的method不存在或类型错误")
							}
							var headers map[string]string
							headersJson, ok := httpConfig.CheckGet("headers")
							if ok {
								headers = make(map[string]string)
								headerData, _ := headersJson.MarshalJSON()
								err = json.Unmarshal(headerData, &headers)
								if err != nil {
									return nil, err
								}
							}
							//if err != nil {
							//return nil, errors.New("http_config中的headers不存在或类型错误")
							//}
							if method == "POST" {
								body, err := httpConfig.Get("body").String()
								if err != nil {
									return nil, errors.New("http_config中的body不存在或类型错误")
								}
								successResConfig.SuccessHttpConfig = &model.RequestInfo{Method: method, Url: url, Headers: headers, Body: body}
							} else {
								successResConfig.SuccessHttpConfig = &model.RequestInfo{Method: method, Url: url, Headers: headers}
							}

							break
						default:
							return nil, errors.New("http_config不支持的请求类型")
						}

					} else {
						return nil, errors.New("请问你怎么发http请求,传音吗?")
					}
					break
				}
			}

		} else {
			successResConfig = nil
		}
	} else if ignoreJson, ok := responseConfig.CheckGet("ignore_response_config"); ok {
		ignoreResConfig = &IgnoreResponseConfig{}
		ignoreData, err := ignoreJson.MarshalJSON()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(ignoreData, ignoreResConfig)
		if err != nil {
			return nil, err
		}
	} else if failedJson, ok := responseConfig.CheckGet("failed_response_config"); ok {
		failedResConfig = &FailedResponseConfig{}
		failedData, err := failedJson.MarshalJSON()
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(failedData, failedResConfig)
		if err != nil {
			return nil, err
		}
	}
	return &ResponseConfig{successResConfig, failedResConfig, ignoreResConfig}, nil
}
