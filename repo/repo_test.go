package repo

import (
	"testing"
	_ "github.com/go-sql-driver/mysql"
	"ICityDataEngine/model"
	"flag"
	"os"
	"github.com/bitly/go-simplejson"
	"ICityDataEngine/requester"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

func TestGenerateRequest(t *testing.T) {
	requestConfigJson := `{
        "type": "http",
        "url": "http://www.icity24.cn/icity/as/app/getYearBill",
        "method": 2,
        "body_type": 3,
        "variables_config": {
            "variables": [
                {
                    "name": "mobile",
                    "data_to": 1,
                    "data_from": 1,
                    "mapping_name": "mobilePhone"
                },
                {
                    "name": "idCard",
                    "data_to": 1,
                    "data_from": 1,
                    "mapping_name": "idCard"
                },
                {
                    "name": "city_code",
                    "data_to": 1,
                    "data_from": 1,
                    "mapping_name": "custId",
					"data_type":"int"
                }
            ],
            "db_config": {
                "db_type": "mysql",
                "user_name": "root",
                "password": "123456a?",
                "db_ip": "172.22.16.139",
                "port": 3306,
                "db_name": "icity",
                "sql": "select custId,mobilePhone,idCard from cust_customer_check"
            }
        }
    }`
	json, err := simplejson.NewJson([]byte(requestConfigJson))
	if err != nil {
		t.Fatal(err)
	}
	requestConfig, err := model.NewHttpRequestConfig(json, "fadadfasdf")
	if err != nil {
		t.Fatal(err)
	}
	t.Log("NewHttpRequestConfig success")
	err = requester.GenerateRequest(requestConfig)
	if err != nil {
		t.Fatal(err)
	}
}

//func _Test_Query(t *testing.T) {
//	config := model.MySqlConfig{UserName: "root", PassWord: "123456a?",
//		DBAddress: "172.22.16.139", Port: 3306, DBName: "icity", SqlSentence: "SELECT * FROM cust_customer_check where mobilePhone = '15069085636'"}
//	t.Log(config.GetDBDataSource())
//	db, err := sql.Open("mysql", config.GetDBDataSource())
//	if err != nil {
//		utils.CheckPanicError(err)
//	}
//	defer db.Close()
//	//rows := db.QueryRow(config.SqlSentence)
//	rows, err := db.Query(config.SqlSentence)
//	if err != nil {
//		t.Error(err)
//	}
//
//	count := 0
//	log.Println("===========")
//	cols, _ := rows.Columns()
//	log.Println(cols)
//	rc := newMapStringScan(cols)
//
//	for rows.Next() {
//		err := rc.Update(rows)
//		if err != nil {
//			fmt.Println(err)
//			return
//		}
//
//		cv := rc.Get()
//		log.Printf("%#v\n\n", cv)
//
//	}
//	log.Println(count)
//}

//type mapStringScan struct {
//	// cp are the column pointers
//	cp []interface{}
//	// row contains the final result
//	row      map[string]string
//	colCount int
//	colNames []string
//}
