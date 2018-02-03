package repo

import (
	"ICityDataEngine/model"
	_ "github.com/go-sql-driver/mysql"
	"database/sql"
	"ICityDataEngine/util"
)

func Query(config model.MySqlConfig) (error){
	db, err := sql.Open("mysql", config.GetDBDataSource())
	if err != nil {
		util.CheckPanicError(err)
	}
	defer db.Close()
	//rows := db.QueryRow(config.SqlSentence)
	rows, err := db.Query(config.SqlSentence)
	if err != nil {
		return err
	}

	var phone string
	for rows.Next() {
		rows.Scan(&phone)
	}
	return nil
}
