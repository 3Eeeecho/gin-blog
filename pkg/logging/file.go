package logging

import (
	"fmt"
	"os"
	"time"

	"github.com/3Eeeecho/go-gin-example/pkg/file"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
)

func getLogFilePath() string {
	return setting.AppSetting.RuntimeRootPath + setting.AppSetting.LogSavePath
}

func getLogFileName() string {
	return fmt.Sprintf("%s%s.%s",
		setting.AppSetting.LogSaveName,
		time.Now().Format(setting.AppSetting.TimeFormat),
		setting.AppSetting.LogFileExt,
	)
}

func openLogFile(filePath, fileName string) (*os.File, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("os.Getwd err: %v", err)
	}

	src := dir + "/" + filePath
	perm := file.CheckFilePermission(src)
	//没有权限，返回错误
	if perm {
		return nil, fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	err = file.IsNotExistMkDir(src)
	if err != nil {
		return nil, fmt.Errorf("file.IsNotExistMkDir src: %s, err: %v", src, err)
	}

	f, err := file.Open(src+fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to openfile :%v", err)
	}

	return f, nil
}
