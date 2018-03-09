package model

import (
	"strconv"
	"ICityDataEngine/i"
	"database/sql"
	"IcityMessageBus/utils"
	"log"
	_ "github.com/go-sql-driver/mysql"
	"sync"
)

var lock *sync.RWMutex

func init() {
	lock = &sync.RWMutex{}
}

type ConnectDBError struct {
}

func (ConnectDBError) Error() string {
	return "连接数据库失败"
}

type MySqlConfig struct {
	UserName    string `json:"user_name"`
	PassWord    string `json:"password"`
	DBAddress   string `json:"db_ip"`
	Port        int    `json:"port"`
	DBName      string `json:"db_name"`
	SqlSentence string `json:"sql"`
	//QueryParamRepo i.IMySqlParamRepo
}

type MySqlSaveConfig struct {
	MySqlConfig
	dbConnection *sql.Tx
	db           *sql.DB
	insertCount  int
}

func (config *MySqlSaveConfig) ExecInsert(values ...interface{}) (error) {
	err := config.initConnection()
	if err != nil {
		return err
	}
	_, err = config.dbConnection.Exec(config.GetSqlSentence(), values...)
	if err != nil {
		return err
	}
	return nil
}

func (config *MySqlSaveConfig) Commit() (error) {
	defer func() {
		if config.db != nil {
			config.db.Close()
		}
		config.dbConnection = nil
		config.db = nil
	}()
	err := config.dbConnection.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (config *MySqlSaveConfig) initConnection() (error) {

	if config.dbConnection == nil && config.db == nil {
		lock.Lock()
		if config.dbConnection != nil {
			return nil
		}
		db, err := sql.Open(config.GetDBType(), config.GetDBDataSource())

		//defer func() {
		//	if db != nil {
		//		db.Close()
		//	}
		//}()

		if err != nil {
			return ConnectDBError{}
		}
		config.db = db

		config.dbConnection, err = db.Begin()
		if err != nil {
			config.db = nil
			return ConnectDBError{}
		}
		lock.Unlock()
		return nil
	}
	return nil
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
