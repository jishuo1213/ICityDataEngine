package model

import (
	"strconv"
)

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
