package server

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"ICityDataEngine/log"
	"io/ioutil"
	"ICityDataEngine/scheduler"
	"ICityDataEngine/model"
	"ICityDataEngine/repo"
	"gopkg.in/mgo.v2/bson"
	"fmt"
)

func dataEngineHandle(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	switch ps.ByName("action") {
	case "add":
		jobConfig, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Error(err)
			res := model.HttpRes{Code: model.ERR_SERVER, Msg: err.Error()}
			writer.Write([]byte(res.String()))
			return
		}
		jobId := bson.NewObjectId()
		engineJob, err := model.ParseConfig(string(jobConfig), jobId.String())
		if err != nil {
			log.Error(err)
			res := model.HttpRes{Code: model.ERR_PARSE_CONFIG, Msg: err.Error()}
			writer.Write([]byte(res.String()))
			return
		}
		engineJob.Id = jobId
		err = repo.AddJob(engineJob)
		if err != nil {
			log.Error(err)
			res := model.HttpRes{Code: model.ERR_INSERT_JOB, Msg: err.Error()}
			writer.Write([]byte(res.String()))
			return
		}
		err = scheduler.AddNewJob(engineJob)
		if err != nil {
			log.Error(err)
			res := model.HttpRes{Code: model.ERR_ADD_SCHEDULER, Msg: err.Error()}
			writer.Write([]byte(res.String()))
			return
		}

		data := make(map[string]interface{})
		data["job_id"] = engineJob.Id

		res := model.HttpRes{Code: model.SUCCESS, Data: data}

		writer.Write([]byte(res.String()))
		break
	case "delete":
		break
	}
}

func Start() {
	router := httprouter.New()
	router.POST("/icity/data/engine/:action", httprouter.Handle(func(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
		log.Info("request-----start")
		dataEngineHandle(writer, request, ps)
		log.Info("request-----end")
	}))

	fmt.Println("start server======================")
	http.ListenAndServe(":1215", router)
}
