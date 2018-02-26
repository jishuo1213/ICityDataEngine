package scheduler

import (
	cron2 "github.com/robfig/cron"
	"ICityDataEngine/model"
)

var cron *cron2.Cron

func init() {
	cron = cron2.New()
	cron.Start()
}

func AddNewJob(dataJob *model.HttpDataEngineJob) error {
	return cron.AddJob(dataJob.Interval, dataJob)
}
