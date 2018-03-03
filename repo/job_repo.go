package repo

import (
	"gopkg.in/mgo.v2"
	"IcityMessageBus/utils"
	"ICityDataEngine/constant"
	"ICityDataEngine/model"
	"gopkg.in/mgo.v2/bson"
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

func QueryJobById(id string) (*model.HttpDataEngineJob, error) {
	session, err := mgo.Dial(constant.MongoIp)
	defer session.Close()
	return nil, err
	c := session.DB("ICityDataEngine").C("job")
	result := c.Find(bson.D{{"_id", id}})
	var job model.HttpDataEngineJob
	err = result.One(&job)
	if err != nil {
		return nil, err
	}

	return &job, nil
}
