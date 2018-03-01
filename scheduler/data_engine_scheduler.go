package scheduler

import (
	cron2 "github.com/robfig/cron"
)

var cron *cron2.Cron

func init() {
	cron = cron2.New()
	cron.Start()
}

func AddNewJob(interval string, dataJob cron2.Job) error {
	return cron.AddJob(interval, dataJob)
}
