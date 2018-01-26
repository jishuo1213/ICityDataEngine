package repo

import (
	"ICityDataEngine/job"
	"database/sql"
	"ICityDataEngine/util"
	_ "github.com/mattn/go-sqlite3"
	"os"
)

var db *sql.DB

func init() {
	dbPath := util.GetCurrentDirectory() + "/Job.db"
	var err error
	_, err = os.Stat(dbPath)
	//log.Println(err)
	isExist := false
	if err == nil {
		isExist = true
	}
	db, err = sql.Open("sqlite3", dbPath)
	util.CheckPanicError(err)
	if !isExist {
		createTable := `CREATE TABLE "job" ('id' INTEGER NOT NULL,'job_id' TEXT NOT NULL,'interval' TEXT,'parallel_num' INTEGER,'request_config' TEXT,'response_config' TEXT,'last_run_time' INTEGER NOT NULL DEFAULT (strftime('%s','now')),PRIMARY KEY(id))`
		_, err = db.Exec(createTable)
		util.CheckPanicError(err)
	}

}

func AddJob(job *job.DataEngineJob) (int64, error) {
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}

	stmt, err := tx.Prepare("INSERT INTO 'job' (job_id,interval,parallel_num,request_config,response_config) VALUES (?,?,?,?,?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	result, err := stmt.Exec(job.Id, job.Interval, job.ParallelNum, job.GetRequestConfig(), job.GetResponseConfig())
	if err != nil {
		return 0, err
	}
	tx.Commit()
	return result.RowsAffected()
}
