package file

import (
	"errors"
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path"
)

func GetFileSize(f multipart.File) (int, error) {
	content, err := io.ReadAll(f)
	return len(content), err
}

func GetFileExt(filename string) string {
	return path.Ext(filename)
}

func CheckFileNotExist(src string) bool {
	_, err := os.Stat(src)
	return errors.Is(err, fs.ErrNotExist)
}

func CheckFilePermission(src string) bool {
	_, err := os.Stat(src)
	return errors.Is(err, fs.ErrPermission)
}

func MkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}
	return nil
}

func IsNotExistMkDir(src string) error {
	if CheckFileNotExist(src) {
		return MkDir(src)
	}
	return nil
}

func Open(name string, flag int, perm os.FileMode) (*os.File, error) {
	f, err := os.OpenFile(name, flag, perm)
	if err != nil {
		return nil, err
	}

	return f, nil
}
