package model

import "net/http"

type HttpVariables struct {
	Name      string
	Type      int //类型：包含两大类，固定值和非固定值 固定值0 数据库1 文件2
	Value     string
	DBMapping string
}

type HttpRequestConfig struct {
	Url         string
	Method      string
	ContentType string
	Headers     []HttpVariables
	Params      []HttpVariables
	SqlConfig   SqlParamConfig
}

func (config *HttpRequestConfig) GenerateRequest() error {
	http.Request{}
}
