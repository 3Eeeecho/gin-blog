package app

import (
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/gin-gonic/gin"
)

type Gin struct {
	C *gin.Context
}

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (g *Gin) Response(httpStatus, errCode int, data interface{}) {
	g.C.JSON(httpStatus, &Response{
		Code: errCode,
		Msg:  e.GetMsg(errCode),
		Data: data,
	})
}
