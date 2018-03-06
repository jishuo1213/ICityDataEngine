package repo

import (
	"gopkg.in/mgo.v2"
	"IcityMessageBus/utils"
	"ICityDataEngine/constant"
	"ICityDataEngine/model"
	"gopkg.in/mgo.v2/bson"
	"log"
	"github.com/bitly/go-simplejson"
	"encoding/json"
)

func AddJob(job *simplejson.Json) (error) {
	mapJob, err := job.Map()
	if err != nil {
		return err
	}

	session, err := mgo.Dial(constant.MongoIp)
	defer session.Close()

	if err != nil {
		utils.CheckPanicError(err)
	}

	c := session.DB("ICityDataEngine").C("job")

	err = c.Insert(mapJob)
	if err != nil {
		return err
	}
	return nil
}

func QueryJobById(id string) (*model.HttpDataEngineJob, error) {
	log.Println("QueryJobById" + id)
	session, err := mgo.Dial(constant.MongoIp)
	defer session.Close()
	c := session.DB("ICityDataEngine").C("job")
	result := c.Find(bson.D{{"_id", bson.ObjectIdHex(id)}})
	//job := model.HttpDataEngineJob{}
	resultMap := make(map[string]interface{})
	err = result.One(&resultMap)
	if err != nil {
		return nil, err
	}
	jsonJobData, err := json.Marshal(resultMap)
	if err != nil {
		return nil, err
	}
	jsonJob, err := simplejson.NewJson(jsonJobData)
	if err != nil {
		return nil, err
	}
	return model.ParseConfig(jsonJob, id)
}
