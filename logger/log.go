package logger

import (
	"os"
	"log"
	"ICityDataEngine/constant"
)

var errorLogger *log.Logger
var fileLogger *log.Logger

func init() {

	normalLogFileWriter, err := os.Create("DataEngine.out")
	if err != nil {
		log.Println("init log failed")
		log.Fatal(err)
		return
	}

	errLogFileWriter, err := os.Create("DataEngineError.out")
	if err != nil {
		log.Println("init error log failed")
		log.Fatal(err)
		return
	}

	//infoLogger = log.New(os.Stdout, "info", log.LstdFlags|log.Llongfile)
	errorLogger = log.New(errLogFileWriter, "info", log.LstdFlags|log.Llongfile)
	fileLogger = log.New(normalLogFileWriter, "info", log.LstdFlags|log.Llongfile)
}

func Error(v ...interface{}) {
	if constant.DEBUG {
		log.Println(v)
	} else {
		errorLogger.Println(v)
	}
}

func Record(v ...interface{}) {
	if constant.DEBUG {
		log.Println(v)
	} else {
		fileLogger.Println(v)
	}
}
