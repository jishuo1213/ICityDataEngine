package job

import (
	"ICityDataEngine/util"
	"gopkg.in/mgo.v2/bson"
)

type DataEngineJob struct {
	//jobConfig *config.JobConfig
	Id             bson.ObjectId          `bson:"_id"`
	Interval       string                 `bson:"interval"`
	ParallelNum    int                    `bson:"parallel_num"`
	RequestConfig  map[string]interface{} `bson:"request_config"`
	ResponseConfig map[string]interface{} `bson:"response_config"`
}

func (job *DataEngineJob) Run() {
	//requestType := string.(job.RequestConfig["type"])
	//switch requestType {
	//case "http":
	//	url := string.(job.RequestConfig["url"])
	//	method := string.(job.RequestConfig["method"])
	//
	//	break
	//}
}

func (job *DataEngineJob) GetRequestConfig() (string) {
	return util.ToJsonStr(job.RequestConfig, "{}")
}

func (job *DataEngineJob) GetResponseConfig() (string) {
	//responseConfig, err := json.Marshal(job.RequestConfig)
	return util.ToJsonStr(job.ResponseConfig, "{}")
}
