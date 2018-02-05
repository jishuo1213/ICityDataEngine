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

func AddNewJob(dataJob *job.DataEngineJob) error {
	return cron.AddJob(dataJob.Interval, dataJob)
}
