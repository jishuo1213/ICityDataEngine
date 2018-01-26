package main

import (
	"ICityDataEngine/server"
)

func main() {
	//data := make(map[string]interface{})
	//
	//data["res1"] = 1
	//data["res2"] = "aaa"
	////res := model.HttpRes{100, nil}
	//
	//jsonRes, _ := json.Marshal(data)
	//fmt.Println(string(jsonRes))
	server.Start()
}
