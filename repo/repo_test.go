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
	"fmt"
	"reflect"
)

func TestMain(m *testing.M) {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "/tmp")
	flag.Set("v", "3")
	flag.Parse()

	ret := m.Run()
	os.Exit(ret)
}

type SqlString struct {
	a string
}

func Test_Query(t *testing.T) {
	config := model.MySqlConfig{UserName: "root", PassWord: "123456a?",
		DBAddress: "172.22.16.139", Port: 3306, DBName: "icity", SqlSentence: "SELECT custId, mobilePhone,idCard FROM cust_customer_check"}

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
	//var phone string
	//result := make([]interface{}, 0, 2)
	//value := "aa"
	//var value SqlString
	//value = "aaa"
	//value := SqlString{}
	//result = append(result, reflect.ValueOf(value).FieldByName(strings.Title("a")).Addr().Interface())
	//reflect.ValueOf().FieldByName().Addr().Interface()
	cols, _ := rows.Columns()
	for rows.Next() {
		//fmt.Println(rows.Columns())
		////rows.Scan()
		//columnsType, _ := rows.ColumnTypes()
		//for _, cType := range columnsType {
		//	fmt.Println(cType.Name())
		//	value := reflect.Zero(cType.ScanType())
		//	result = append(result, value)
		//}

		//rows.Scan(result...)
		//fmt.Println(result[0] == nil)
		//
		////sql.RawBytes{}
		//count++

		columns := make([]interface{}, len(cols))
		columnPointers := make([]interface{}, len(cols))
		for i := range columns {
			columnPointers[i] = &columns[i]
		}

		// Scan the result into the column pointers...
		if err := rows.Scan(columnPointers...); err != nil {
			fmt.Println(err)
			return
		}

		// Create our map, and retrieve the value for each column from the pointers slice,
		// storing it in the map with the name of the column as the key.
		m := make(map[string]interface{})
		for i, colName := range cols {
			val := columnPointers[i].(*interface{})
			fmt.Println(reflect.TypeOf(*val))
			m[colName] = *val
		}

		// Outputs: map[columnName:value columnName2:value2 columnName3:value3 ...]
		fmt.Println(m["mobilePhone"].(string))
	}
	//t.Error(count)
	log.Println(count)
}
