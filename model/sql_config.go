package model

import (
	"strconv"
	"ICityDataEngine/i"
	"database/sql"
	"IcityMessageBus/utils"
	"log"
	_ "github.com/go-sql-driver/mysql"
)

type MySqlConfig struct {
	UserName    string `bson:"user_name"`
	PassWord    string `bson:"password"`
	DBAddress   string `bson:"db_ip"`
	Port        int    `bson:"port"`
	DBName      string `bson:"db_name"`
	SqlSentence string `bson:"sql"`
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

func (config *MySqlConfig) QuerySqlParams(parser i.ISqlResultParser, statusParser i.IDealRunStatus) error {
	//return repo.QuerySqlParams(config, parser)
	statusParser("开始连接数据库----")
	db, err := sql.Open(config.GetDBType(), config.GetDBDataSource())
	defer func() {
		if db != nil {
			db.Close()
		}
	}()
	if err != nil {
		utils.CheckPanicError(err)
	}

	statusParser("执行sql语句----" + config.GetSqlSentence())

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
	return parser(rows, statusParser)
}
