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

type ISqlExecConfig interface {
	sqlConfig
	ExecSql()
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
	IsSuccessResponse(response []byte) (bool)
	IsIgnoreResponse(response []byte) (bool)
	DealSuccessResponse(body []byte) error
	DealFailedRequest(request IRequestInfo) error
}

type IResponseSaver interface {
	Save(response []byte) error
}

//type ISqlResultParser interface {
//	Parse(*sql.Rows) error
//}
