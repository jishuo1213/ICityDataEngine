package server

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"io/ioutil"
	"ICityDataEngine/model"
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/repo"
)

func TestJobHandler(ws *websocket.Conn) {
	msg, err := ioutil.ReadAll(ws)
	if err != nil {
		res := model.NewWSRes(model.ERR_SERVER, err.Error(), nil, "")
		ws.Write([]byte(res.String()))
		ws.Close()
		return
	}

	data, err := simplejson.NewJson(msg)
	if err != nil {
		res := model.NewWSRes(model.ERR_FORMAT, err.Error(), nil, "")
		ws.Write([]byte(res.String()))
		ws.Close()
		return
	}

	cmdId, err := data.Get("cmd_id").String()
	if err != nil {
		res := model.NewWSRes(model.ERR_FORMAT, err.Error(), nil, "")
		ws.Write([]byte(res.String()))
		return
	}

	action, err := data.Get("action").String()
	if err != nil {
		res := model.NewWSRes(model.ERR_FORMAT, err.Error(), nil, "")
		ws.Write([]byte(res.String()))
		return
	}

	switch action {
	case "connect":
		res := model.NewWSRes(model.SUCCESS, "connect success", nil, cmdId)
		ws.Write([]byte(res.String()))
		break
	case "test":
		jobId, err := data.Get("job_id").String()
		if err != nil {
			res := model.NewWSRes(model.ERR_FORMAT, err.Error(), nil, cmdId)
			ws.Write([]byte(res.String()))
			return
		}
		job, err := repo.QueryJobById(jobId)
		if err != nil {
			res := model.NewWSRes(model.ERR_SERVER, "query failed:"+err.Error(), nil, cmdId)
			ws.Write([]byte(res.String()))
			return
		}
		testStatusChan := make(chan string, 20)
		job.TestRun(testStatusChan)
		for status := range testStatusChan {
			res := model.NewWSRes(model.SUCCESS, status, nil, cmdId)
			ws.Write([]byte(res.String()))
		}
		return
		break
	case "disconnect":
		break
	}
}

func StartWebSocketServer() {
	http.Handle("/icity/data/engine/ws", websocket.Handler(TestJobHandler))
	if err := http.ListenAndServe(":1217", nil); err != nil {
		log.Fatal(err)
	}
}
