package repo

import (
	"ICityDataEngine/model"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"ICityDataEngine/util"
)



func QuerySqlParams(parser func(rows *sql.Rows) error, config model.SqlParamConfig) (error) {
	if config == nil {
		return parser(nil)
	}
	rowsList := make([]*sql.Rows, 0, 1)
	dbList := make([]*sql.DB, 0, 1)
	defer func() {
		for _, db := range dbList {
			db.Close()
		}
		for _, rows := range rowsList {
			rows.Close()
		}
	}()
	//for _, config := range configs {
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
	//}
	return parser(rowsList[0])
}
