package main

import (
	"ICityDataEngine/model"
	"encoding/json"
	"fmt"
)

func main() {
	//data := make(map[string]interface{})

	//data["res1"] = 1
	//data["res2"] = "aaa"
	res := model.HttpRes{100, nil}

	jsonRes, _ := json.Marshal(res)
	fmt.Println(string(jsonRes))
}
