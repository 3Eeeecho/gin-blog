package api

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/pkg/app"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/util"
	"github.com/3Eeeecho/go-gin-example/service/auth_service"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

// GetAuth 获取授权（登录）
// @Summary 获取授权 Token
// @Description 通过用户名和密码进行验证，成功后返回一个 Token，供后续请求验证使用。
// @Tags 认证
// @Accept  json
// @Produce json
// @Param username query string true "用户名"  // 用户名，必填
// @Param password query string true "密码"  // 密码，必填
// @Success 200 {object} app.Response "返回成功信息，包含 Token"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 401 {object} app.Response "认证失败，用户名或密码错误"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	username := c.Query("username")
	password := c.Query("password")

	g := app.Gin{C: c}
	valid := validation.Validation{}

	a := auth{Username: username, Password: password}
	ok, _ := valid.Valid(&a)

	if !ok {
		app.MakrErrors(valid.Errors)
		g.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	authService := auth_service.Auth{Username: username, Password: password}
	isExist, err := authService.Check()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_AUTH_CHECK_TOKEN_FAIL, nil)
		return
	}

	if !isExist {
		g.Response(http.StatusUnauthorized, e.ERROR_AUTH, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_AUTH_TOKEN, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}
