package repo

import (
	"ICityDataEngine/model"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"ICityDataEngine/util"
)

func QueryMySqlParams(parser func(rows ...*sql.Rows) error, configs ... model.SqlParamConfig) (error) {
	rowsList := make([]*sql.Rows, 0, len(configs))
	dbList := make([]*sql.DB, 0, len(configs))
	defer func() {
		for _, db := range dbList {
			db.Close()
		}
		for _, rows := range rowsList {
			rows.Close()
		}
	}()
	for _, config := range configs {
		db, err := sql.Open(config.GetDBType(), config.GetDBDataSource())
		if err != nil {
			util.CheckPanicError(err)
		}
		dbList = append(dbList, db)
		//defer db.Close()
		rows, err := db.Query(config.GetSqlSentence())
		rowsList = append(rowsList, rows)
		//rows := db.QueryRow(config.SqlSentence)
		if err != nil {
			return err
		}
	}
	return parser(rowsList...)
}
