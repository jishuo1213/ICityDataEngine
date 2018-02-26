package repo

import (
	"gopkg.in/mgo.v2"
	"ICityDataEngine/util"
	"ICityDataEngine/model"
)

func AddJob(job *model.HttpDataEngineJob) (error) {
	//tx, err := db.Begin()
	//if err != nil {
	//	return 0, err
	//}
	//
	//stmt, err := tx.Prepare("INSERT INTO 'job' (job_id,interval,parallel_num,request_config,response_config) VALUES (?,?,?,?,?)")
	//if err != nil {
	//	return 0, err
	//}
	//defer stmt.Close()
	//
	//result, err := stmt.Exec(job.Id, job.Interval, job.ParallelNum, job.GetRequestConfig(), job.GetResponseConfig())
	//if err != nil {
	//	return 0, err
	//}
	//tx.Commit()
	//return result.RowsAffected()

	session, err := mgo.Dial("172.22.16.144:27017")
	defer session.Close()

	if err != nil {
		util.CheckPanicError(err)
	}

	c := session.DB("ICityDataEngine").C("job")

	err = c.Insert(job)
	if err != nil {
		return err
	}
	return nil
}
