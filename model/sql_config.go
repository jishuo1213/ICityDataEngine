package model

import (
	"strconv"
	"database/sql"
	"IcityMessageBus/utils"
	"log"
	"ICityDataEngine/i"
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

func (config *MySqlConfig) GetSqlSentence() string {
	return config.SqlSentence
}

func (config *MySqlConfig) GetDBType() string {
	return "mysql"
}

func (config *MySqlConfig) QuerySqlParams(parser i.ISqlResultParser) error {
	db, err := sql.Open(config.GetDBType(), config.GetDBDataSource())
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	if err != nil {
		utils.CheckPanicError(err)
	}

	log.Println("query:" + config.GetSqlSentence())
	rows, err := db.Query(config.GetSqlSentence())
	defer func() {
		if rows != nil {
			rows.Close()
		}
	}()
	if err != nil {
		log.Println(err)
		return err
	}
	return parser.Parse(rows)
}
