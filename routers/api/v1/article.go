package v1

import (
	"fmt"
	"net/http"

	"github.com/3Eeeecho/go-gin-example/models"
	"github.com/3Eeeecho/go-gin-example/pkg/app"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/pkg/util"
	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetArticle 获取单篇文章的详细信息
// @Summary 获取单篇文章的详细信息
// @Description 根据文章ID获取文章数据
// @Tags 文章
// @Accept  json
// @Produce json
// @Param id path int true "文章ID"  // 必填参数，文章的ID
// @Success 200 {object} app.Response "返回文章信息"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 404 {object} app.Response "文章不存在"
// @Router /api/v1/articles/{id} [get]
func GetArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	g := app.Gin{C: c}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	var data interface{}
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			code = e.SUCCESS
			data = models.GetArticle(id)
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			fmt.Println("Logging error:", err.Key, err.Message) // 调试信息
			logging.Info(err.Key, err.Message)
		}
	}

	g.Response(http.StatusOK, code, data)
}

// GetArticles 获取文章列表
// @Summary 获取文章列表
// @Description 根据请求参数（如状态、标签ID）返回文章列表数据和总数
// @Tags 文章
// @Accept  json
// @Produce json
// @Param state query int false "文章状态"  // 可选参数，0: 草稿，1: 已发布
// @Param tag_id query int false "标签ID"  // 可选参数，标签ID，必须大于0
// @Success 200 {object} app.Response "返回文章列表和总数"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /articles [get]
func GetArticles(c *gin.Context) {
	data := make(map[string]interface{})
	maps := make(map[string]interface{})
	g := app.Gin{C: c}
	valid := validation.Validation{}
	state := -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
		maps["state"] = state
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
		maps["tag_id"] = tagId
	}

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		code = e.SUCCESS
		data["lists"] = models.GetArticles(util.GetPage(c), setting.AppSetting.PageSize, maps)
		data["total"] = models.GetArticleTotal(maps)
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	g.Response(http.StatusOK, code, data)
}

// AddArticle 新增文章
// @Summary 新增一篇文章
// @Description 通过传入文章的相关信息（标签ID、标题、简述、内容、创建人、状态）来新增一篇文章。
// @Tags 文章
// @Accept  json
// @Produce json
// @Param tag_id query int true "标签ID"  // 标签ID，必填，必须大于0
// @Param title query string true "标题"  // 文章标题，必填
// @Param desc query string true "简述"  // 文章简述，必填
// @Param content query string true "内容"  // 文章内容，必填
// @Param created_by query string true "创建人"  // 创建人的名称，必填
// @Param state query int false "状态"  // 文章状态（0: 草稿, 1: 已发布），可选，默认0
// @Success 200 {object} app.Response "成功返回数据"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 404 {object} app.Response "标签不存在"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /articles [post]
func AddArticle(c *gin.Context) {
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	createdBy := c.Query("created_by")
	state := com.StrTo(c.DefaultQuery("state", "0")).MustInt()
	cover_image_url := c.Query("cover_image_url")

	g := app.Gin{C: c}
	valid := validation.Validation{}
	valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	valid.Required(title, "title").Message("标题不能为空")
	valid.Required(desc, "desc").Message("简述不能为空")
	valid.Required(content, "content").Message("内容不能为空")
	valid.Required(createdBy, "created_by").Message("创建人不能为空")
	valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	valid.Required(cover_image_url, "cover_image_url").Message("图片不能为空")
	valid.Max(cover_image_url, 255, "cover_image_url").Message("url最长不能超过255字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistTagByID(tagId) {
			data := make(map[string]interface{})

			data["tag_id"] = tagId
			data["title"] = title
			data["desc"] = desc
			data["content"] = content
			data["created_by"] = createdBy
			data["state"] = state

			models.AddArticle(data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_TAG
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	g.Response(http.StatusOK, code, nil)
}

// EditArticle 修改文章
// @Summary 修改文章
// @Description 通过文章ID和更新的参数修改文章信息（如标签ID、标题、简述、内容、修改人、状态）
// @Tags 文章
// @Accept  json
// @Produce json
// @Param id path int true "文章ID"  // 文章ID，必填，必须大于0
// @Param tag_id query int false "标签ID"  // 标签ID，可选，必须大于0
// @Param title query string false "标题"  // 文章标题，可选
// @Param desc query string false "简述"  // 文章简述，可选
// @Param content query string false "内容"  // 文章内容，可选
// @Param modified_by query string true "修改人"  // 修改人的名称，必填
// @Param state query int false "状态"  // 文章状态（0: 草稿, 1: 已发布），可选，默认不修改
// @Success 200 {object} app.Response "返回成功信息"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 404 {object} app.Response "文章不存在"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /articles/{id} [put]
func EditArticle(c *gin.Context) {
	g := app.Gin{C: c}
	valid := validation.Validation{}

	id := com.StrTo(c.Param("id")).MustInt()
	tagId := com.StrTo(c.Query("tag_id")).MustInt()
	title := c.Query("title")
	desc := c.Query("desc")
	content := c.Query("content")
	modifiedBy := c.Query("modified_by")
	cover_image_url := c.Query("cover_image_url")

	state := -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo("state").MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	valid.Min(id, 1, "id").Message("ID必须大于0")
	valid.Min(tagId, 1, "tag_id").Message("TagID必须大于0")
	valid.MaxSize(title, 100, "title").Message("标题最长为100字符")
	valid.MaxSize(desc, 255, "desc").Message("简述最长为255字符")
	valid.MaxSize(content, 65535, "content").Message("内容最长为65535字符")
	valid.Required(modifiedBy, "modified_by").Message("修改人不能为空")
	valid.MaxSize(modifiedBy, 100, "modified_by").Message("修改人最长为100字符")
	valid.Required(cover_image_url, "cover_image_url").Message("图片不能为空")
	valid.Max(cover_image_url, 255, "cover_image_url").Message("url最长不能超过255字符")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			data := make(map[string]interface{})
			if tagId > 0 {
				data["tag_id"] = tagId
			}
			if title != "" {
				data["title"] = title
			}
			if desc != "" {
				data["desc"] = desc
			}
			if content != "" {
				data["content"] = content
			}

			data["modified_by"] = modifiedBy
			models.EditArticle(id, data)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	g.Response(http.StatusOK, code, nil)

}

// DeleteArticle 删除文章
// @Summary 删除文章
// @Description 通过文章ID删除指定的文章
// @Tags 文章
// @Accept  json
// @Produce json
// @Param id path int true "文章ID"  // 文章ID，必填，必须大于0
// @Success 200 {object} app.Response "返回成功信息"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 404 {object} app.Response "文章不存在"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	g := app.Gin{C: c}
	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	code := e.INVALID_PARAMS
	if !valid.HasErrors() {
		if models.ExistArticleByID(id) {
			models.DeleteArticle(id)
			code = e.SUCCESS
		} else {
			code = e.ERROR_NOT_EXIST_ARTICLE
		}
	} else {
		for _, err := range valid.Errors {
			logging.Info(err.Key, err.Message)
		}
	}

	g.Response(http.StatusOK, code, nil)
}
