package repo

import (
	"testing"
	"database/sql"
	"ICityDataEngine/util"
	_ "github.com/go-sql-driver/mysql"
	"ICityDataEngine/model"
	"flag"
	"os"
	"log"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

func Test_Query(t *testing.T) {
	config := model.MySqlConfig{UserName: "root", PassWord: "123456a?",
		DBAddress: "172.22.16.139", Port: 3306, DBName: "icity", SqlSentence: "SELECT phone FROM cust_customer_action"}

	t.Log(config.GetDBDataSource())
	db, err := sql.Open("mysql", config.GetDBDataSource())
	if err != nil {
		util.CheckPanicError(err)
	}
	defer db.Close()
	//rows := db.QueryRow(config.SqlSentence)
	rows, err := db.Query(config.SqlSentence)
	if err != nil {
		t.Error(err)
	}

	count := 0
	log.Println("===========")
	var phone string
	for rows.Next() {
		rows.Scan(&phone)
		count++
	}
	//t.Error(count)
	log.Println(count)
}
