package repo

import (
	"database/sql"
	"IcityMessageBus/utils"
	"log"
	"ICityDataEngine/i"
	_ "github.com/go-sql-driver/mysql"
)

type QueryMySqlParamsRepo struct {
}

func (repo *QueryMySqlParamsRepo) QuerySqlParams(config i.ISqlParamConfig, parser i.ISqlResultParser) error {
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
	return parser(rows, nil)
}
