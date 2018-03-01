package repo

import (
	"gopkg.in/mgo.v2"
	"IcityMessageBus/utils"
	"ICityDataEngine/constant"
	"ICityDataEngine/model"
)

func AddJob(job *model.HttpDataEngineJob) (error) {
	session, err := mgo.Dial(constant.MongoIp)
	defer session.Close()

	if err != nil {
		utils.CheckPanicError(err)
	}

	c := session.DB("ICityDataEngine").C("job")

	err = c.Insert(job)
	if err != nil {
		return err
	}
	return nil
}
