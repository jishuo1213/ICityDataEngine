package model

import (
	"github.com/kataras/iris/core/errors"
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/logger"
	"ICityDataEngine/requester"
	"IcityMessageBus/cmsp"
	"ICityDataEngine/constant"
	"IcityMessageBus/utils"
	"IcityMessageBus/model"
	"ICityDataEngine/i"
	"strconv"
	"log"
	"sync"
)

type HttpDataEngineJob struct {
	//jobConfig *config.JobConfig
	Id             string            `bson:"_id"`
	Name           string            `bson:"friendly_name"`
	Interval       string            `bson:"interval"`
	ParallelNum    int               `bson:"parallel_num"`
	RequestConfig  i.IRequestConfig  `bson:"request_config"`
	ResponseConfig i.IResponseConfig `bson:"response_config"`
	LastRunTime    int64             `bson:"last_run_time"`
}

func (job *HttpDataEngineJob) Run() {
	log.Println("start run id = " + job.Id)
	statusParser := func(s interface{}) {
		switch s.(type) {
		case string:
			logger.Record(s)
			break
		case error:
			logger.Error(s)
			break
		default:
			break
		}
	}
	topicName, _ := utils.DigestMessage([]byte(job.Id))
	topicName = topicName[8:24]
	err := requester.GenerateRequest(job.RequestConfig, topicName, statusParser)
	defer cmsp.DisconnectCmsp(constant.CMSPIP, constant.CMSPPort, topicName)
	if err != nil {

		statusParser(err)
		//logger.Error(err)
		statusParser(job.Name + "生成请求失败,定时任务执行失败")
		//logger.Record()
		return
	}

	requestChan := make(chan *model.RequestInfo, 20)
	go readMessageLopper(topicName, requestChan, statusParser)
	wg := &sync.WaitGroup{}
	for index := 0; index < job.ParallelNum; index++ {
		wg.Add(1)
		go dealRequest(requestChan, statusParser, wg, job.ResponseConfig)
	}
	wg.Wait()

	log.Println("run======================end")
}

func (job *HttpDataEngineJob) TestRun(statusChan chan<- string) {
	queueId, _ := utils.DigestMessage([]byte(job.Id + "_test"))
	queueId = queueId[8:24]
	log.Println("TestRun:" + queueId)

	statusParser := func(s interface{}) {
		switch s.(type) {
		case string:
			statusChan <- s.(string)
			break
		case error:
			statusChan <- s.(error).Error()
			break
		default:
			break
		}
	}
	statusChan <- "开始测试:" + job.Name
	statusChan <- "开始生成请求------"
	err := requester.GenerateRequest(job.RequestConfig, queueId, statusParser)
	defer cmsp.DisconnectCmsp(constant.CMSPIP, constant.CMSPPort, queueId)
	if err != nil {
		//logger.Error(err)
		//logger.Record(job.Id + "生成请求失败,定时任务执行失败")
		statusChan <- "生成请求失败,错误:" + err.Error()
		return
	}
	statusChan <- "生成http请求成功,开始发送"
	requestChan := make(chan *model.RequestInfo, 20)
	go readMessageLopper(queueId, requestChan, statusParser)
	wg := &sync.WaitGroup{}
	for index := 0; index < job.ParallelNum; index++ {
		wg.Add(1)
		go dealRequest(requestChan, statusParser, wg, job.ResponseConfig)
	}
	wg.Wait()
	log.Println("testRun end")
}

func dealRequest(requestChan <-chan *model.RequestInfo, statusParser i.IDealRunStatus, wg *sync.WaitGroup, config i.IResponseConfig) {
	defer func() {
		log.Println("dealRequest end")
		wg.Done()
	}()
	statusParser("开始发送http请求------")
	for requestInfo := range requestChan {
		httpRequest, err := requestInfo.GenerateRequest()
		if err != nil {
			statusParser(errors.New("Generate request failed err :" + err.Error() + " url:" +
				requestInfo.Url + "body:" + requestInfo.Body))
			continue
		}

		code, body, err := requester.SendHttpRequest(httpRequest)
		//resp, err := httpClient.Do(httpRequest)

		if err != nil {
			statusParser(errors.New("send http request failed: request info:" +
				requestInfo.GetFormatString() + "error:" + err.Error()))
			continue
		}

		statusParser("执行:" + requestInfo.GetFormatString() + " 返回结果:")
		//body, err := ioutil.ReadAll(resp.Body)
		statusParser("状态码：" + strconv.Itoa(code) + "返回结果：" + string(body))
		if config != nil {
			if ok := config.IsSuccessResponse(code, body); ok {
				config.DealSuccessResponse(requestInfo, code, body)
			} else if ok = config.IsIgnoreResponse(code, body); ok {
				return
			} else {
				config.DealFailedRequest(requestInfo)
			}
		}
	}
}

func readMessageLopper(queueId string, requestChan chan<- *model.RequestInfo, statusParser i.IDealRunStatus) {
	statusParser("开始读取请求队列----")
	for {
		msg, err := cmsp.ReadMsgFromQueueNet(constant.CMSPIP, constant.CMSPPort, queueId)
		if err != nil {
			if _, ok := err.(cmsp.ConnectCmspError); ok {
				//logger.Error("connect cmsp failed ")
				statusParser(errors.New("connect cmsp failed "))
				close(requestChan)
				break
			} else {
				log.Println("read message end")
				close(requestChan)
				break
			}
		}

		var request model.RequestInfo
		err = utils.DecodeObject(&request, msg)
		if err != nil {
			//logger.Error("read msg from queue but can not decode msg = " + string(msg))
			statusParser(errors.New("read msg from queue but can not decode msg = " +
				string(msg) + " err:" + err.Error()))
			continue
		}
		requestChan <- &request
	}
}

func ParseConfig(jobConfig *simplejson.Json, id string) (*HttpDataEngineJob, error) {
	interval, err := jobConfig.Get("interval").String()
	if err != nil {
		return nil, errors.New("interval不存在或类型错误")
	}
	parallelNum, err := jobConfig.Get("parallel_num").Int()
	if err != nil {
		return nil, errors.New("parallel_num不存在或类型错误")
	}
	requestConfig, isHaveRequest := jobConfig.CheckGet("request")
	if !isHaveRequest {
		return nil, errors.New("request_config不存在或类型错误")
	}

	responseJson, isHaveResponse := jobConfig.CheckGet("response")

	var responseConfig i.IResponseConfig

	requestType, err := requestConfig.Get("type").String()

	if err != nil {
		return nil, errors.New("request_config中type不存在或类型错误")
	}
	var httpRequestConfig *httpRequestConfig
	switch requestType {
	case "http":
		//var err error
		httpRequestConfig, err = NewHttpRequestConfig(requestConfig, id)
		if err != nil {
			return nil, err
		}
		if isHaveResponse {
			responseConfig, err = NewHttpResponseConfig(responseJson, id)
			if err != nil {
				return nil, err
			}
		}
		break
	default:
		return nil, errors.New("暂不支持" + requestType + "类型服务")
	}

	return &HttpDataEngineJob{Interval: interval,
		ParallelNum: parallelNum, RequestConfig: httpRequestConfig, ResponseConfig: responseConfig}, nil
}
