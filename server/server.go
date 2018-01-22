package server

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"ICityDataEngine/log"
)

func dataEngineHandle(writer http.ResponseWriter, request *http.Request, ps httprouter.Params) {
	switch ps.ByName("action") {
	case "add":

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
