package model

import (
	"testing"
	"encoding/json"
	"flag"
	"os"
	"fmt"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

func TestGetInsertValues(t *testing.T) {

	response := `{
    "state": "1",
    "message": "调用成功",
    "code": "0000",
    "error": "",
    "result": [
        [
            {
                "statusType": "socialPension",
                "level": 2,
                "cityCode": "370100",
                "topTitle": "社保查询",
                "banner": "1",
                "description": "查询个人养老储存额信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/jnSocialSecurityNew/index.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/user_pension_iv.png",
                "displayValue": "12940.99",
                "appId": 319,
                "isShare": "0",
                "name": "养老储存额"
            },
            {
                "statusType": "socialMedical",
                "level": 2,
                "cityCode": "370100",
                "topTitle": "社保查询",
                "banner": "1",
                "description": "查询个人医保余额信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/jnSocialSecurityNew/index.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/user_medicare_iv.png",
                "displayValue": "1590.17",
                "appId": 319,
                "isShare": "0",
                "name": "医保余额"
            },
            {
                "statusType": "accumulationFund",
                "level": 2,
                "cityCode": "370100",
                "banner": "1",
                "description": "查询个人公积金余额信息",
                "gotoUrl": "https://www.icity24.cn/icity/apps/areas/jinan/accumulation-fund/index.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/user_accumulation_iv.png",
                "displayValue": "1666.93",
                "appId": 275,
                "isShare": "0",
                "name": "公积金余额"
            },
            {
                "statusType": "incomeTax",
                "level": 2,
                "cityCode": "370100",
                "topTitle": "个税查询",
                "banner": "1",
                "description": "个税信息查询",
                "gotoUrl": "https://www.icity24.cn/icity/apps/areas/jinan/income-tax/index.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/user_tax_iv.png",
                "displayValue": "10011.75",
                "appId": 699,
                "isShare": "0",
                "name": "个税查询"
            }
        ],
        [
            {
                "statusType": "violation",
                "level": 2,
                "cityCode": "370100",
                "banner": "4",
                "description": "用户可以通过违章查询应用，快速获取车辆的违章信息，打开应用添加车辆，输入车牌号、车架号、发动机号，即可查询该车辆的违章信息，也可以通过切换车辆，查询多个车辆的违章信息。",
                "gotoUrl": "http://www.icity24.cn/icity/apps/trafficViolation/result.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/violation_inquiry.png",
                "displayValue": "<font color=\"#ff0000\">5个违章未处理</font>",
                "appId": 159,
                "isShare": "1",
                "name": "违章查询"
            },
            {
                "statusType": "license",
                "level": 2,
                "cityCode": "370100",
                "banner": "4",
                "description": "记分快速查",
                "gotoUrl": "https://www.icity24.cn/icity/apps/areas/jinan/driving-license-integral/index.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/20171218093828311642.png",
                "displayValue": "去绑定驾驶证",
                "appId": 1012,
                "isShare": "0",
                "name": "驾驶证记分"
            }
        ],
        [
            {
                "statusType": "waterFee",
                "level": 2,
                "cityCode": "370100",
                "banner": "3",
                "description": "水费查缴信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/ecpayNew/payWater.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/water_bill.png",
                "displayValue": "去绑定缴费户号",
                "appId": 270,
                "isShare": "0",
                "name": "水费"
            },
            {
                "statusType": "electricFee",
                "level": 2,
                "cityCode": "370100",
                "banner": "3",
                "description": "电费查缴信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/ecpayNew/payPower.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/electric_charge.png",
                "displayValue": "去绑定缴费户号",
                "appId": 271,
                "isShare": "0",
                "name": "电费"
            },
            {
                "statusType": "gasFee",
                "level": 2,
                "cityCode": "370100",
                "banner": "3",
                "description": "燃气费查缴信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/ecpayNew/payGas.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/gas_bill.png",
                "displayValue": "去绑定缴费户号",
                "appId": 272,
                "isShare": "0",
                "name": "燃气费"
            },
            {
                "statusType": "heatingFee",
                "level": 2,
                "cityCode": "370100",
                "banner": "3",
                "description": "暖气费查缴信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/ecpayNew/payHeat.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/heating_costs.png",
                "displayValue": "去绑定缴费户号",
                "appId": 273,
                "isShare": "0",
                "name": "暖气费"
            }
        ],
        [
            {
                "statusType": "woodpecker",
                "level": 2,
                "cityCode": "370100",
                "banner": "2",
                "description": "“啄木鸟在行动”重点关注大气污染违法违规问题和政府部门不作为、乱作为问题。市民可以通过点击我要反映，进行拍照或相册上传，反映身边的污染环境、违法违规现象，提交成功后，可以在历史记录中查看已提交的问题及有关部门的答复。",
                "gotoUrl": "http://www.icity24.cn/icity/apps/woodpecker/index.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/wood_keeper.png",
                "displayValue": "为了青山绿水",
                "appId": 259,
                "isShare": "0",
                "name": "啄木鸟行动"
            },
            {
                "statusType": "hotLine",
                "level": 2,
                "cityCode": "370100",
                "banner": "2",
                "description": "市民热线为用户提供了向12345发布事项咨询的途径，用户可以通过填写标题、内容、图片，并选择是否需要回复，即可发布事项，提交成功后系统会返回事项编号和查询码，用户可在回复查询中查询该事项的进度。",
                "gotoUrl": "http://www.icity24.cn/icity/apps/12345/stand_alone_form.html",
                "imgUrl": "http://icity24.infile.cloud.inspur.com/Image/Life/hot_line_12345.png",
                "displayValue": "服务全天候",
                "appId": 266,
                "isShare": "0",
                "name": "12345热线"
            }
        ]
    ],
    "server": {
        "build": "10",
        "time": "2018-03-07 19:05:16",
        "version": "3.0.0"
    },
    "markInfo": {
        "isPop": 0,
        "mark": 0
    }
}`

	responseMap := make(map[string]interface{})
	json.Unmarshal([]byte(response), &responseMap)

	template := `{
    "result": [
        [
            {
                "cityCode": "370100",
                "description": "查询个人养老储存额信息",
                "gotoUrl": "http://www.icity24.cn/icity/apps/jnSocialSecurityNew/index.html",
                "displayValue": "12940.99"
            }
        ]
    ]
}`
	templateMap := make(map[string]interface{})
	json.Unmarshal([]byte(template), &templateMap)

	saveKeysMaps := make(map[string]*SaveKey)

	saveKeysMaps["cityCode"] = &SaveKey{"cityCode", "string", nil}
	saveKeysMaps["gotoUrl"] = &SaveKey{"gotoUrl", "string", nil}
	saveKeysMaps["description"] = &SaveKey{"description", "string", nil}
	saveKeysMaps["displayValue"] = &SaveKey{"displayValue", "string", nil}

	//res, keys, err := getInsertValues(responseMap, templateMap, saveKeysMaps)
	fmt.Println(findSameInJsonBody(templateMap, responseMap, true))
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(strconv.Itoa(len(res)) + "================")
	//log.Println(res)

	//log.Println(keys)

	//for i, v := range keys {
	//	fmt.Print(v)
	//	fmt.Print(":")
	//	fmt.Print(res[i])
	//	fmt.Print("\n")
	//}
}
