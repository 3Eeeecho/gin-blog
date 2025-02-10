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

// GetImageFullUrl 生成图片的完整访问 URL
// name: 图片文件名
// 返回: 完整的图片 URL，格式为 "ImagePrefixUrl/ImageSavePath/name"
func GetImageFullUrl(name string) string {
	return setting.AppSetting.ImagePrefixUrl + "/" + GetImagePath() + name
}

// GetImageName 生成唯一的图片文件名
// name: 原始文件名
// 返回: 唯一的图片文件名，格式为 "MD5(文件名) + 扩展名"
func GetImageName(name string) string {
	ext := path.Ext(name)                     // 获取文件扩展名
	fileName := strings.TrimSuffix(name, ext) // 去除扩展名
	fileName = util.EncodeMD5(fileName)       // 对文件名部分进行 MD5 哈希

	return fileName + ext // 拼接哈希值和扩展名
}

// GetImagePath 获取图片的存储路径
// 返回: 配置文件中设置的图片存储路径
func GetImagePath() string {
	return setting.AppSetting.ImageSavePath
}

// GetImageFullPath 获取图片的完整存储路径
// 返回: 运行时根路径 + 图片存储路径
func GetImageFullPath() string {
	return setting.AppSetting.RuntimeRootPath + GetImagePath()
}

// CheckImageExt 检查图片扩展名是否允许
// fileName: 图片文件名
// 返回: 如果扩展名允许返回 true，否则返回 false
func CheckImageExt(fileName string) bool {
	ext := file.GetFileExt(fileName) // 获取文件扩展名
	for _, allowExt := range setting.AppSetting.ImageAllowExts {
		if strings.EqualFold(ext, allowExt) { // 忽略大小写比较
			return true
		}
	}
	return false
}

// CheckImageSize 检查图片大小是否超过限制
// f: 上传的文件
// 返回: 如果文件大小未超过限制返回 true，否则返回 false
func CheckImageSize(f multipart.File) bool {
	size, err := file.GetFileSize(f) // 获取文件大小
	if err != nil {
		logging.Warn(err) // 记录警告日志
		return false
	}
	return size <= setting.AppSetting.ImageMaxSize // 比较文件大小和限制
}

// CheckImage 检查图片存储路径是否存在并验证权限
// src: 图片存储的相对路径
// 返回: 如果检查通过返回 nil，否则返回错误
func CheckImage(src string) error {
	dir, err := os.Getwd() // 获取当前工作目录
	if err != nil {
		return fmt.Errorf("os.Getwd err: %v", err)
	}

	// 检查目录是否存在，如果不存在则创建
	err = file.IsNotExistMkDir(path.Join(dir, src))
	if err != nil {
		return fmt.Errorf("file.IsNotExistMkDir err: %v", err)
	}

	// 检查目录权限
	perm := file.CheckFilePermission(src)
	if perm {
		return fmt.Errorf("file.CheckPermission Permission denied src: %s", src)
	}

	return nil
}
