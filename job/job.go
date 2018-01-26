package job

import (
	"ICityDataEngine/util"
)

type DataEngineJob struct {
	//jobConfig *config.JobConfig
	Id             string
	Interval       string
	ParallelNum    int
	RequestConfig  map[string]interface{}
	ResponseConfig map[string]interface{}
}

func (job *DataEngineJob) Run() {
	requestType := string.(job.RequestConfig["type"])
	switch requestType {
	case "http":
		url := string.(job.RequestConfig["url"])
		method := string.(job.RequestConfig["method"])

		break
	}
}

func (job *DataEngineJob) GetRequestConfig() (string) {
	return util.ToJsonStr(job.RequestConfig, "{}")
}

func (job *DataEngineJob) GetResponseConfig() (string) {
	//responseConfig, err := json.Marshal(job.RequestConfig)
	return util.ToJsonStr(job.ResponseConfig, "{}")
}
