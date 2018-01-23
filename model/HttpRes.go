package model

import "encoding/json"

type HttpRes struct {
	Code int `json:"code"`
	//Msg  string `json:"msg"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func (httpRes *HttpRes) String() string {
	res, _ := json.Marshal(httpRes)
	return string(res)
}
