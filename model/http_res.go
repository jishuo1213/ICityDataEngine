package model

import "encoding/json"

type HttpErrorCode int

const (
	SUCCESS          HttpErrorCode = 100
	ERR_SERVER
	ERR_PARSE_CONFIG
	ERR_INSERT_JOB
	ERR_ADD_SCHEDULER
)

type HttpRes struct {
	Code HttpErrorCode          `json:"code"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

func (httpRes *HttpRes) String() string {
	res, _ := json.Marshal(httpRes)
	return string(res)
}
