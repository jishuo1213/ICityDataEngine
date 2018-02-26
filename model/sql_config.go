package model

import (
	"strconv"
	//"database/sql"
	//"ICityDataEngine/repo"
)

type SqlParamConfig interface {
	GetDBDataSource() string
	GetSqlSentence() string
	GetDBType() string
	//QueryAndParseParams(parser func(rows ...*sql.Rows) error) error //查询并处理查询出来的参数
}

type MySqlConfig struct {
	UserName    string
	PassWord    string
	DBAddress   string
	Port        int
	DBName      string
	SqlSentence string
}

func (config *MySqlConfig) GetDBDataSource() string {
	return config.UserName + ":" + config.PassWord + "@tcp(" +
		config.DBAddress + ":" + strconv.Itoa(config.Port) + ")/" + config.DBName + "?charset=utf8"
}

func (config *MySqlConfig) GetSqlSentence() string {
	return config.SqlSentence
}

func (config *MySqlConfig) GetDBType() string {
	return "mysql"
}

//func (config *MySqlConfig) QueryAndParseParams(parser func(rows ...*sql.Rows) error) error {
//	return repo.QuerySqlParams(parser, config)
//}