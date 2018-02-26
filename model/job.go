package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/kataras/iris/core/errors"
	"github.com/bitly/go-simplejson"
)

type HttpDataEngineJob struct {
	//jobConfig *config.JobConfig
	Id          bson.ObjectId `bson:"_id"`
	Interval    string        `bson:"interval"`
	ParallelNum int           `bson:"parallel_num"`
	//RequestConfig  map[string]interface{} `bson:"request_config"`
	RequestConfig  HttpRequestConfig      `bson:"request_config"`
	ResponseConfig map[string]interface{} `bson:"response_config"`
	LastRunTime    int64                  `bson:"last_run_time"`
}

func (job *HttpDataEngineJob) Run() {

	//requestType := string.(job.RequestConfig["type"])
	//fmt.Println(requestType)
}

func initJobConfig(config string) (*simplejson.Json, error) {
	js, err := simplejson.NewJson([]byte(config))
	if err != nil {
		return nil, err
	}
	return js, nil
}

func ParseConfig(config string) (*HttpDataEngineJob, error) {
	jobConfig, err := initJobConfig(config)
	if err != nil {
		return nil, errors.New("json格式解析失败")
	}
	interval, err := jobConfig.Get("interval").String()
	if err != nil {
		return nil, errors.New("interval不存在或类型错误")
	}
	parallelNum, err := jobConfig.Get("parallel_num").Int()
	if err != nil {
		return nil, errors.New("parallel_num不存在或类型错误")
	}
	requestConfig := jobConfig.Get("request")
	responseConfig, err := jobConfig.Get("response_config").Map()
	if err != nil {
		return nil, errors.New("response_config不存在或类型错误")
	}

	requestType, err := requestConfig.Get("type").String()

	if err != nil {
		return nil, errors.New("request_config中type不存在或类型错误")
	}
	var httpRequestConfig *HttpRequestConfig
	switch requestType {
	case "http":
		var err error
		httpRequestConfig, err = NewHttpRequestConfig(requestConfig)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("暂不支持" + requestType + "类型服务")
	}

	return &HttpDataEngineJob{Interval: interval,
		ParallelNum: parallelNum, RequestConfig: *httpRequestConfig, ResponseConfig: responseConfig}, nil
}
