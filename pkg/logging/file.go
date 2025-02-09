package logging

import (
	"fmt"
	"log"
	"os"
	"time"
)

var (
	LogSavePath = "runtime/logs/"
	LogSaveName = "log"
	LogFileExt  = "log"
	TimeFormat  = "20060101"
)

func getLogFilePath() string {
	return LogSavePath
}

func getLogFileFullPath() string {
	return fmt.Sprintf("%s%s.%s", getLogFilePath(), time.Now().Format(TimeFormat), LogFileExt)
}

func mkDir() {
	err := os.MkdirAll(getLogFilePath(), os.ModePerm)
	if err != nil {
		panic(err)
	}
}

func openLogFile(filepath string) *os.File {
	_, err := os.Stat(filepath)
	switch {
	case os.IsNotExist(err):
		mkDir()
	case os.IsPermission(err):
		log.Fatalf("Permission :%v", err)
	}

	handle, err := os.OpenFile(filepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("Fail to OpenFile :%v", err)
	}

	return handle
}
