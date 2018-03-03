package model

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/kataras/iris/core/errors"
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/logger"
	"ICityDataEngine/requester"
	"IcityMessageBus/cmsp"
	"ICityDataEngine/constant"
	"IcityMessageBus/utils"
	"IcityMessageBus/model"
	"net/http"
	"time"
	"io/ioutil"
)

var httpClient *http.Client

func init() {
	httpClient = &http.Client{
		Timeout: time.Second * 20,
	}
}

type HttpDataEngineJob struct {
	//jobConfig *config.JobConfig
	Id             bson.ObjectId          `bson:"_id"`
	Name           string                 `bson:"friendly_name"`
	Interval       string                 `bson:"interval"`
	ParallelNum    int                    `bson:"parallel_num"`
	RequestConfig  httpRequestConfig      `bson:"request_config"`
	ResponseConfig map[string]interface{} `bson:"response_config"`
	LastRunTime    int64                  `bson:"last_run_time"`
}

func (job *HttpDataEngineJob) Run() {
	err := requester.GenerateRequest(&job.RequestConfig, job.Id.String())
	if err != nil {
		logger.Error(err)
		logger.Record(job.Id + "生成请求失败,定时任务执行失败")
		return
	}

	requestChan := make(chan *model.RequestInfo, 20)
	go readMessageLopper(job.Id.String(), requestChan)
	for index := 0; index < job.ParallelNum; index++ {
		go dealRequest(requestChan)
	}
}

func (job *HttpDataEngineJob) TestRun() {
	err := requester.GenerateRequest(&job.RequestConfig, job.Id.String()+"_test")
	if err != nil {
		logger.Error(err)
		logger.Record(job.Id + "生成请求失败,定时任务执行失败")
		return
	}

}

func dealRequest(requestChan <-chan *model.RequestInfo) {
	for requestInfo := range requestChan {
		httpRequest, err := requestInfo.GenerateRequest()
		if err != nil {
			logger.Error("Generate request failed err :" + err.Error() + " url:" +
				requestInfo.Url + "body:" + requestInfo.Body)
			continue
		}

		resp, err := httpClient.Do(httpRequest)
		if err != nil {
			logger.Error("send http request failed: request info:" +
				requestInfo.GetFormatString() + "error:" + err.Error())
			continue
		}

		logger.Record(requestInfo.GetFormatString() + " res:")
		logger.Record(resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		logger.Record(string(body))
	}
}

//func sendHttpRequest(request *http.Request) (int, *http.Header, []byte, error) {
//	//client := &http.Client{}
//	resp, err := httpClient.Do(request)
//	defer func() {
//		if resp != nil {
//			resp.Body.Close()
//		}
//	}()
//	if err != nil {
//		log.Print(err)
//		return 500, nil, nil, err
//	} else {
//		body, err := ioutil.ReadAll(resp.Body)
//		log.Print(string(body))
//		if err != nil {
//			log.Print(err)
//			return resp.StatusCode, nil, nil, err
//		}
//		return resp.StatusCode, &resp.Header, body, nil
//	}
//}

func readMessageLopper(queueId string, requestChan chan<- *model.RequestInfo) {
	for {
		msg, err := cmsp.ReadMsgFromQueueNet(constant.CMSPIP, constant.CMSPPort, queueId)
		if err != nil {
			if _, ok := err.(ConnectCmspError); ok {
				logger.Error("connect cmsp failed ")
				close(requestChan)
				break
			} else {
				close(requestChan)
				break
			}
		}

		var request model.RequestInfo
		err = utils.DecodeObject(&request, msg)
		if err != nil {
			logger.Error("read msg from queue but can not decode msg = " + string(msg))
			break
		}
		requestChan <- &request
	}
}

func initJobConfig(config []byte) (*simplejson.Json, error) {
	js, err := simplejson.NewJson(config)
	if err != nil {
		return nil, err
	}
	return js, nil
}

func ParseConfig(config []byte, id string) (*HttpDataEngineJob, error) {
	jobConfig, err := initJobConfig(config)
	if err != nil {
		return nil, errors.New("json格式解析失败")
	}
	interval, err := jobConfig.Get("interval").String()
	if err != nil {
		return nil, errors.New("interval不存在或类型错误")
	}
	parallelNum, err := jobConfig.Get("parallel_num").Int()
	if err != nil {
		return nil, errors.New("parallel_num不存在或类型错误")
	}
	requestConfig := jobConfig.Get("request")
	responseConfig, err := jobConfig.Get("response_config").Map()
	if err != nil {
		return nil, errors.New("response_config不存在或类型错误")
	}

	requestType, err := requestConfig.Get("type").String()

	if err != nil {
		return nil, errors.New("request_config中type不存在或类型错误")
	}
	var httpRequestConfig *httpRequestConfig
	switch requestType {
	case "http":
		var err error
		httpRequestConfig, err = NewHttpRequestConfig(requestConfig, id)
		if err != nil {
			return nil, err
		}
		break
	default:
		return nil, errors.New("暂不支持" + requestType + "类型服务")
	}

	return &HttpDataEngineJob{Interval: interval,
		ParallelNum: parallelNum, RequestConfig: *httpRequestConfig, ResponseConfig: responseConfig}, nil
}
