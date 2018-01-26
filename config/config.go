package config

import (
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/job"
)

//
//const BODY_JSON_TYPE int = 1
//const BODY_X_FORM_DATA int = 2
//
//const PARAM_DB ParamsSourceType = 1
//const PARAM_FILE ParamsSourceType = 2
//const PARAM_LIST ParamsSourceType = 3
//
////const ParamsSourceType  =
//
//type ParamsSourceType int
//
//type jobConfig struct {
//	interval    string
//	parallelNum int
//}
//
////表示各种参数的设置,参数是从哪里获取的
//type paramsConfig struct {
//	name       string
//	sourceType ParamsSourceType
//}
//
//type ParamsGetter interface {
//	getParamData() ([]string, error)
//}
//
//type dbParamsConfig struct {
//	paramsConfig
//	dbType    string
//	userName  string
//	password  string
//	address   string
//	dbName    string
//	tableName string
//	columns   string
//}
//
//func (config *dbParamsConfig) getParamSql() string {
//
//}
//
//type httpRequestConfig struct {
//	url         string
//	method      string
//	bodyType    int
//	headersName map[string]paramsConfig
//	paramsName  map[string]paramsConfig
//}
//
//type JobConfig interface {
//}

//type JobConfig map[string]interface{}

func initJobConfig(config string) (*simplejson.Json, error) {
	js, err := simplejson.NewJson([]byte(config))
	if err != nil {
		//log.Error(err)
		return nil, err
	}
	//var jobConfig JobConfig
	//jobConfig, err = js.Map()
	//if err != nil {
	//	return nil, err
	//}

	return js, nil
}

func ParseConfig(config string) (*job.DataEngineJob, error) {
	jobConfig, err := initJobConfig(config)
	if err != nil {
		return nil, err
	}
	//id, err := jobConfig.Get("id").String()
	interval, err := jobConfig.Get("interval").String()
	parallelNum, err := jobConfig.Get("parallel_num").Int()
	requestConfig, err := jobConfig.Get("request").Map()
	responseConfig, err := jobConfig.Get("response_config").Map()

	return &job.DataEngineJob{Interval: interval,
		ParallelNum: parallelNum, RequestConfig: requestConfig, ResponseConfig: responseConfig}, nil

}
