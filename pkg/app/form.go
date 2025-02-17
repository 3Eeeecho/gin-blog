package app

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

func BindAndValue(c *gin.Context, form interface{}) (int, int) {
	err := c.Bind(form)
	if err != nil {
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		return http.StatusInternalServerError, e.ERROR
	}

	if !check {
		MakrErrors(valid.Errors)
		return http.StatusBadRequest, e.INVALID_PARAMS
	}

	return http.StatusOK, e.SUCCESS
}
