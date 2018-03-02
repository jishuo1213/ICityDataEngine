package model

import (
	"strconv"
	"ICityDataEngine/i"
)

type MySqlConfig struct {
	UserName    string
	PassWord    string
	DBAddress   string
	Port        int
	DBName      string
	SqlSentence string
	//QueryParamRepo i.IMySqlParamRepo
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

func (config *MySqlConfig) QuerySqlParams(repo i.ISqlParamRepo, parser i.ISqlResultParser) error {
	return repo.QuerySqlParams(config, parser)
}
