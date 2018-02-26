package util

import (
	"path/filepath"
	"os"
	"strings"
	"log"
	"encoding/json"
)

func GetCurrentDirectory() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		//log.Fatal(err)
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1)
}

func CheckPanicError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func ToJsonStr(data interface{}, defaultValue string) string {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return defaultValue
	}

	return string(jsonData)
}