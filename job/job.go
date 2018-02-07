package job

import (
	"gopkg.in/mgo.v2/bson"
	"ICityDataEngine/repo"
	"github.com/kataras/iris/core/errors"
	"fmt"
	"github.com/bitly/go-simplejson"
)

type DataEngineJob struct {
	//jobConfig *config.JobConfig
	Id             bson.ObjectId          `bson:"_id"`
	Interval       string                 `bson:"interval"`
	ParallelNum    int                    `bson:"parallel_num"`
	RequestConfig  map[string]interface{} `bson:"request_config"`
	ResponseConfig map[string]interface{} `bson:"response_config"`
	LastRunTime    int64                  `bson:"last_run_time"`
	paramParser    repo.ParamParser
}

func (job *DataEngineJob) Run() {

	requestType := string.(job.RequestConfig["type"])
	fmt.Println(requestType)
}

func initJobConfig(config string) (*simplejson.Json, error) {
	js, err := simplejson.NewJson([]byte(config))
	if err != nil {
		return nil, err
	}
	return js, nil
}

func ParseConfig(config string) (*DataEngineJob, error) {
	jobConfig, err := initJobConfig(config)
	if err != nil {
		return nil, err
	}
	//id, err := jobConfig.Get("id").String()
	interval, err := jobConfig.Get("interval").String()
	if err != nil {
		return nil, err
	}
	parallelNum, err := jobConfig.Get("parallel_num").Int()
	if err != nil {
		return nil, err
	}
	requestConfig, err := jobConfig.Get("request").Map()
	if err != nil {
		return nil, err
	}
	responseConfig, err := jobConfig.Get("response_config").Map()
	if err != nil {
		return nil, err
	}

	return &DataEngineJob{Interval: interval,
		ParallelNum: parallelNum, RequestConfig: requestConfig, ResponseConfig: responseConfig}, nil

}

func (job *DataEngineJob) InitJob() error {
	requestType, ok := job.RequestConfig["type"].(string)
	if !ok {
		return errors.New("type is not string")
	}
	switch requestType {
	case "http":

		break
	default:
		return errors.New("暂不支持" + requestType + "类型服务")
	}
}
