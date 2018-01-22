package scheduler

import (
	cron2 "github.com/robfig/cron"
	"ICityDataEngine/job"
)

var cron *cron2.Cron

func init() {
	cron = cron2.New()
	cron.Start()
}

func AddNewJob(job *job.DataEngineJob) error {
	cron.AddJob()
}
