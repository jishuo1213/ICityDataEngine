package job

import (
	"gopkg.in/mgo.v2/bson"
	"ICityDataEngine/repo"
	"github.com/kataras/iris/core/errors"
	"fmt"
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
	//switch requestType {
	//case "http":
	//	url := string.(job.RequestConfig["url"])
	//	method := string.(job.RequestConfig["method"])
	//
	//	break
	//}
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
