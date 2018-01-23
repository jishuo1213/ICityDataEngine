package job

type DataEngineJob struct {
	//jobConfig *config.JobConfig
	Id             string
	Interval       string
	ParallelNum    int
	RequestConfig  map[string]interface{}
	ResponseConfig map[string]interface{}
}

func (job *DataEngineJob) Run() {

}
