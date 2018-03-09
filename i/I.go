package i

import (
	"ICityDataEngine/constant"
	"database/sql"
	"net/http"
)

type ISqlParamRepo interface {
	QuerySqlParams(ISqlParamConfig, ISqlResultParser) error
}

type ISqlResultParser func(*sql.Rows, IDealRunStatus) error

type IDealRunStatus func(interface{})

type IRequestInfo interface {
	GenerateRequest() (*http.Request, error)
	//GetSaveKeys(saveKeys []*model.SaveKey) string
	GetHeaders() map[string]string
	GetBodyParams() map[string]string
}

type sqlConfig interface {
	GetDBDataSource() string
	GetSqlSentence() string
	GetDBType() string
}

type ISqlParamConfig interface {
	sqlConfig
	QuerySqlParams(ISqlResultParser, IDealRunStatus) error
}

type ISqlSaveConfig interface {
	sqlConfig
	ExecInsert(...interface{}) error
	Commit() (error)
}

type IRequestConfig interface {
	GetId() string
	GetMethod() constant.HttpMethod
	GetContentType() constant.ContentType
	GetUrl() string
	GetSqlConfig() ISqlParamConfig
	InitValueHeaders() (map[string]string, map[string]string)
	InitValueParams() (map[string]string, map[string]string)
	GenerateJsonBody(map[string]string) (map[string]interface{}, error)
}

type IResponseConfig interface {
	IsSuccessResponse(code int, response []byte) (bool)
	IsIgnoreResponse(code int, response []byte) (bool)
	DealSuccessResponse(request IRequestInfo, body []byte) error
	DealFailedRequest(request IRequestInfo) error
}

type IResponseSaver interface {
	Save(response []byte, request IRequestInfo) error
}

//type ISqlResultParser interface {
//	Parse(*sql.Rows) error
//}
