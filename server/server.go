package server

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"ICityDataEngine/log"
	"io/ioutil"
	"ICityDataEngine/config"
	"ICityDataEngine/scheduler"
	"ICityDataEngine/model"
)

func dataEngineHandle(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	switch ps.ByName("action") {
	case "add":
		jobConfig, err := ioutil.ReadAll(request.Body)
		if err != nil {
			res := model.HttpRes{Code: 101, Data: nil}
			writer.Write([]byte(res.String()))
			return
		}
		job, err := config.ParseConfig(string(jobConfig))
		if err != nil {
			res := model.HttpRes{Code: 101, Data: nil}
			writer.Write([]byte(res.String()))
		}

		err = scheduler.AddNewJob(job)
		if err != nil {
			res := model.HttpRes{Code: 102, Data: nil}
			writer.Write([]byte(res.String()))
		}



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

	http.ListenAndServe(":1215", router)
}
