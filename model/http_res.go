package model

import "encoding/json"

type HttpErrorCode int

const (
	SUCCESS           HttpErrorCode = iota + 100
	ERR_SERVER
	ERR_PARSE_CONFIG
	ERR_INSERT_JOB
	ERR_ADD_SCHEDULER
	ERR_FORMAT
)

type HttpRes struct {
	Code HttpErrorCode          `json:"code"`
	Msg  string                 `json:"msg,omitempty"`
	Data map[string]interface{} `json:"data,omitempty"`
}

type WSRes struct {
	HttpRes
	CmdId string `json:"cmd_id,omitempty"`
}

func NewWSRes(code HttpErrorCode, msg string, data map[string]interface{}, cmdId string) (*WSRes) {
	return &WSRes{HttpRes{code, msg, data}, cmdId}
}

func (httpRes *HttpRes) String() string {
	res, _ := json.Marshal(httpRes)
	return string(res)
}

func (httpRes *WSRes) String() string {
	res, _ := json.Marshal(httpRes)
	return string(res)
}
