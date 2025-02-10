package upload

import (
	"fmt"
	"mime/multipart"
	"os"
	"path"
	"strings"

	"github.com/3Eeeecho/go-gin-example/pkg/file"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/pkg/util"
)

func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

func GetImageName(name string) string {
	ext := path.Ext(name)
	fileName := strings.TrimSuffix(name, ext)
	fileName = util.EncodeMD5(fileName)

	return fileName + ext
}

func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

func CheckImageExt(fileName string) bool {
	ext := file.GetFileExt(fileName)
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.EqualFold(ext, allowExt) {
			return true
		}
	}
	return false
}

func CheckImageSize(f multipart.File) bool {
	size, err := file.GetFileSize(f)
	if err != nil {
		logging.Warn(err)
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize
}

func CheckImage(src string) error {
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	err = file.IsNotExistMkDir(path.Join(dir, src))
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	perm := file.CheckFilePermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
