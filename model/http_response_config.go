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
	Name string
	Type string
}

type RedisSaveConfig struct {
	Ip      string     `bson:"ip"`
	Port    int        `bson:"port"`
	SaveKey []*SaveKey `bson:"save_key"`
}

type DBSaveConfig struct {
	SqlConfig      i.ISqlExecConfig
	InsertTemplate string
	InsertMapping  []*SaveKey
	BodyType       string
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

		break
	}
	return nil
}

func getInsertValues(map[string]interface{}) ([]interface{}, error) {

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
