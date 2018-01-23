package util

import (
	"path/filepath"
	"os"
	"strings"
	"log"
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