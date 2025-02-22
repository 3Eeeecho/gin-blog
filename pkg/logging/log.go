package logging

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Level int

const (
	DEBUG Level = iota
	INFO
	WARNING
	ERROR
	FATAL
)

var (
	F *os.File

	DefaultPrefix      = ""
	DefaultCallerDepth = 2
	logger             *log.Logger
	logPrefix          = ""
	levelFlags         = []string{"DEBUG", "INFO", "WARN", "ERROR", "FATAL"}
)

func SetUp() {

	filepath := getLogFilePath()
	fileName := getLogFileName()
	F, err := openLogFile(filepath, fileName)
	if err != nil {
		log.Fatalln(err)
	}
	logger = log.New(F, DefaultPrefix, log.LstdFlags)
}

func Debug(v ...interface{}) {
	setPrefix(DEBUG)
	logger.Println(v...)
	F.Sync()
}

func Info(v ...interface{}) {
	setPrefix(INFO)
	logger.Println(v...)
	F.Sync()
}

func Warn(v ...interface{}) {
	setPrefix(WARNING)
	logger.Println(v...)
	F.Sync()
}

func Error(v ...interface{}) {
	setPrefix(ERROR)
	logger.Println(v...)
	F.Sync()
}

func Fatal(v ...interface{}) {
	setPrefix(FATAL)
	logger.Println(v...)
	F.Sync()
}

func setPrefix(level Level) {
	_, file, line, ok := runtime.Caller(DefaultCallerDepth)
	if ok {
		logPrefix = fmt.Sprintf("[%s][%s:%d]", levelFlags[level], filepath.Base(file), line)
	} else {
		logPrefix = fmt.Sprintf("[%s]", levelFlags[level])
	}
	logger.SetPrefix(logPrefix)
}
