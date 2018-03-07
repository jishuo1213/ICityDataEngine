package model

import (
	"ICityDataEngine/i"
	"github.com/bitly/go-simplejson"
	"encoding/json"
	"github.com/kataras/iris/core/errors"
)

type CmspSaveConfig struct {
	//SaveTo string `bson:"save_to"`
	TopicName string `bson:"topic_name"`
}

type SaveKey struct {
	Name  string
	Type  string
	Value []interface{}
}

type RedisSaveConfig struct {
	Ip      string     `bson:"ip"`
	Port    int        `bson:"port"`
	SaveKey []*SaveKey `bson:"save_key"`
}

type DBSaveConfig struct {
	SqlConfig      i.ISqlSaveConfig
	InsertTemplate string
	InsertMapping  []*SaveKey
	BodyType       string
	InsertOrder    []string
}

func (config *DBSaveConfig) Save(response []byte) error {
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
				if _, ok := err.(ConnectDBError); ok {
					return err
				} else {
					//TODO:保存失败的需要处理
					continue
				}
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
	DataType          int    `bson:"data_type"`
	SuccessTemplate   string `bson:"template"`
	SaveConfig        i.IResponseSaver
	SuccessHttpConfig i.IRequestInfo
}

type ResponseConfig struct {
}

func (config *ResponseConfig) IsSuccessResponse(response []byte) (bool) {
	return false
}

func (config *ResponseConfig) IsIgnoreResponse(body []byte) (bool) {
	return false
}
func (config *ResponseConfig) DealSuccessResponse(body []byte) error {
	return nil
}
func (config *ResponseConfig) DealFailedRequest(request i.IRequestInfo) error {
	return nil
}

func NewHttpResponseConfig(responseConfig *simplejson.Json, id string) (*ResponseConfig, error) {
	return nil, nil

}
