package v1

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/models"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetTags 获取标签列表
// @Summary 获取标签列表
// @Description 根据请求的参数（如标签名、状态）获取标签数据
// @Tags 标签
// @Accept  json
// @Produce json
// @Param name query string false "标签名称"  // 可选参数，按标签名称进行过滤
// @Param state query int false "标签状态"  // 可选参数，按状态过滤，0: 禁用，1: 启用
// @Success 200 {object} models.Response "返回标签列表和总数"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	maps := make(map[string]interface{})
	if name != "" {
		maps["name"] = name
	}
	if state != -1 {
		maps["state"] = state
	}

	lists := models.GetTags(util.GetPage(c), setting.AppSetting.PageSize, maps)
	total := models.GetTagTotal(maps)

	data := make(map[string]interface{})
	data["lists"] = lists
	data["total"] = total

	response := models.Success(data)
	c.JSON(http.StatusOK, response)
}

// AddTag 新增文章标签
// @Summary 新增文章标签
// @Description 创建新的标签
// @Tags 标签
// @Accept  json
// @Produce json
// @Param name query string true "标签名称"  // 必填参数，标签名称
// @Param state query int false "标签状态"  // 可选参数，0: 禁用，1: 启用
// @Param created_by query int false "创建人"  // 可选参数，创建人的用户名
// @Success 200 {object} models.Response "返回成功信息"
// @Failure 400 {object} models.Response "参数验证失败"
// @Failure 409 {object} models.Response "标签已存在"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	name := c.Query("name")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	createdBy := c.Query("created_by")

	valid := validation.Validation{}
	valid.Required(name, "name").Message("消息不能为空")
	valid.MaxSize(name, 100, "name").Message("名称最长为100个字符")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.MaxSize(createdBy, 100, "created_by").Message("创建人最长为100字符")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if !models.ExistTagByName(name) {
			code = e.SUCCESS
			models.AddTag(name, state, createdBy)
		} else {
			code = e.ERROR_EXIST_TAG
		}
	}
	response := models.NewResponse(code, e.GetMsg(code), make(map[string]string))
	c.JSON(http.StatusOK, response)
}

// EditTag 修改文章标签
// @Summary 修改文章标签
// @Description 编辑已有标签的信息
// @Tags 标签
// @Accept  json
// @Produce json
// @Param id path int true "标签ID"  // 必填参数，标签ID
// @Param name query string false "标签名称"  // 可选参数，标签名称
// @Param state query int false "标签状态"  // 可选参数，0: 禁用，1: 启用
// @Param modified_by query string true "修改人"  // 必填参数，修改人
// @Success 200 {object} models.Response "返回成功信息"
// @Failure 400 {object} models.Response "参数验证失败"
// @Failure 404 {object} models.Response "标签不存在"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	name := c.Query("name")
	modifiedBy := c.Query("modified_by")

	valid := validation.Validation{}

	var state int = -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Required(id, "id").Message("ID不能为空")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.MaxSize(name, 100, "name").Message("名称最长为100字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			data := make(map[string]interface{})
			data["modified_by"] = modifiedBy
			if name != "" {
				data["name"] = name
			}
			if state != -1 {
				data["state"] = state
			}

			models.EditTag(id, data)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}

	response := models.NewResponse(code, e.GetMsg(code), make(map[string]string))
	c.JSON(http.StatusOK, response)
}

// DeleteTag 删除文章标签
// @Summary 删除文章标签
// @Description 删除指定标签
// @Tags 标签
// @Accept  json
// @Produce json
// @Param id path int true "标签ID"  // 必填参数，标签ID
// @Success 200 {object} models.Response "返回成功信息"
// @Failure 400 {object} models.Response "参数验证失败"
// @Failure 404 {object} models.Response "标签不存在"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不能为空")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		if models.ExistTagByID(id) {
			models.DeleteTag(id)
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	}
	response := models.NewResponse(code, e.GetMsg(code), make(map[string]string))
	c.JSON(http.StatusOK, response)
}
