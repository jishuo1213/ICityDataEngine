package i

import (
	"ICityDataEngine/constant"
	"database/sql"
)

type ISqlParamConfig interface {
	GetDBDataSource() string
	GetSqlSentence() string
	GetDBType() string
	QuerySqlParams(parser ISqlResultParser) error
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

type ISqlResultParser interface {

	Parse(*sql.Rows) error
}
