package server

import (
	"net/http"
	"golang.org/x/net/websocket"
	"log"
	"ICityDataEngine/model"
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/repo"
	"ICityDataEngine/logger"
)

func TestJobHandler(ws *websocket.Conn) {
	log.Println("===============================")
	msg := make([]byte, 512)
	_, err := ws.Read(msg)
	if err != nil {
		log.Println(err)
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
		go func() {
			for status := range testStatusChan {
				res := model.NewWSRes(model.SUCCESS, status, nil, cmdId)
				ws.Write([]byte(res.String()))
			}
		}()
		job.TestRun(testStatusChan)
		close(testStatusChan)
		//testStatusChan <- "测试完成----------"
		res := model.NewWSRes(model.SUCCESS, "测试完成----------", nil, cmdId)
		ws.Write([]byte(res.String()))
		break
	}
}

func StartWebSocketServer() {
	http.Handle("/icity/data/engine/ws", websocket.Handler(TestJobHandler))
	logger.Record("start web socket server---------")
	if err := http.ListenAndServe(":1217", nil); err != nil {
		log.Fatal(err)
	}
}
