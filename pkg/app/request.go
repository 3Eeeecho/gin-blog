package app

import (
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/astaxie/beego/validation"
)

func MakrErrors(errors []*validation.Error) {
	for _, err := range errors {
		logging.Info(err.Key, err.Message)
	}
}
