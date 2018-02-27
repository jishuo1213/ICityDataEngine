package log

import (
	"ICityDataEngine/constant"
	"io"
	"os"
	"log"
)

var infoLogger *log.Logger
var errorLogger *log.Logger

func init() {
	var writer io.Writer
	if constant.DEBUG {
		writer = os.Stdout
	} else {
		fileWriter, err := os.Create("DataENgine.out")
		if err != nil {
			log.Println("init log failed")
			infoLogger = log.New(os.Stdout, "", log.LstdFlags|log.Llongfile)
			errorLogger = log.New(os.Stdout, "info", log.LstdFlags|log.Llongfile)
			return
		}
		writer = fileWriter
	}
	infoLogger = log.New(writer, "info", log.LstdFlags|log.Llongfile)
	errorLogger = log.New(writer, "info", log.LstdFlags|log.Llongfile)
}

func Info(v ...interface{}) {
	infoLogger.Println(v)
}

func Error(v ...interface{}) {
	errorLogger.Println(v)
}
