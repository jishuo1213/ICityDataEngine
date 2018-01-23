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
	if os.IsExist(err) {
		db, err = sql.Open("sqlite3", )
	} else {
		db, err = sql.Open("sqlite3", )
	}
	//if err != nil {
	//	log.Fatal(err)
	//}
	util.CheckPanicError(err)
}

func AddJob(job *job.DataEngineJob) (int, error) {

}
