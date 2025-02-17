package v1

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/pkg/app"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/pkg/util"
	"github.com/3Eeeecho/go-gin-example/service/article_service"
	"github.com/3Eeeecho/go-gin-example/service/tag_service"
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

	if valid.HasErrors() {
		app.MakrErrors(valid.Errors)
		g.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}

	exists, err := articleService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}

	if !exists {
		g.Response(http.StatusInternalServerError, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	article, err := articleService.Get()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_GET_ARTICLE_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, article)
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
// @Router /api/v1/articles [get]
func GetArticles(c *gin.Context) {
	g := app.Gin{C: c}
	valid := validation.Validation{}
	state := -1

	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
		valid.Range(state, 0, 1, "state").Message("状态只允许0或1")
	}

	var tagId int = -1
	if arg := c.Query("tag_id"); arg != "" {
		tagId = com.StrTo(arg).MustInt()
		valid.Min(tagId, 1, "tag_id").Message("标签ID必须大于0")
	}

	if valid.HasErrors() {
		app.MakrErrors(valid.Errors)
		g.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{
		TagID:    tagId,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	total, err := articleService.Count()
	if err != nil {
		g.Response(http.StatusOK, e.ERROR_COUNT_ARTICLE_FAIL, nil)
		return
	}

	articles, err := articleService.GetAll()
	if err != nil {
		g.Response(http.StatusOK, e.ERROR_GET_ARTICLES_FAIL, nil)
		return
	}

	data := make(map[string]interface{})
	data["lists"] = articles
	data["total"] = total

	g.Response(http.StatusOK, e.SUCCESS, data)
}

type AddArticleForm struct {
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	CreatedBy     string `form:"created_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Required;Range(0,1)"`
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
// @Param state query int false "状态"  // 文章状态（0: 草稿, 1: 已发布）
// @Success 200 {object} app.Response "成功返回数据"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 404 {object} app.Response "标签不存在"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /api/v1/articles [post]
func AddArticle(c *gin.Context) {
	var (
		form AddArticleForm
		g    = app.Gin{C: c}
	)

	httpCode, errCode := app.BindAndValue(c, &form)
	if errCode != e.SUCCESS {
		g.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err := tagService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}
	if !exists {
		g.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
		CreatedBy:     form.CreatedBy,
	}
	if err := articleService.Add(); err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_ADD_ARTICLE_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditArticleForm struct {
	ID            int    `form:"id" valid:"Required;Min(1)"`
	TagID         int    `form:"tag_id" valid:"Required;Min(1)"`
	Title         string `form:"title" valid:"Required;MaxSize(100)"`
	Desc          string `form:"desc" valid:"Required;MaxSize(255)"`
	Content       string `form:"content" valid:"Required;MaxSize(65535)"`
	ModifiedBy    string `form:"modified_by" valid:"Required;MaxSize(100)"`
	CoverImageUrl string `form:"cover_image_url" valid:"Required;MaxSize(255)"`
	State         int    `form:"state" valid:"Required;Range(0,1)"`
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
// @Param state query int false "状态"  // 文章状态（0: 草稿, 1: 已发布）
// @Success 200 {object} app.Response "返回成功信息"
// @Failure 400 {object} app.Response "参数验证失败"
// @Failure 404 {object} app.Response "文章不存在"
// @Failure 500 {object} app.Response "服务器错误"
// @Router /api/v1/articles/{id} [put]
func EditArticle(c *gin.Context) {
	var (
		form = EditArticleForm{ID: com.StrTo(c.Param("id")).MustInt()}
		g    = app.Gin{C: c}
	)

	httpCode, errCode := app.BindAndValue(c, form)
	if errCode != e.SUCCESS {
		g.Response(httpCode, errCode, nil)
		return
	}

	articleService := article_service.Article{
		TagID:         form.TagID,
		Title:         form.Title,
		Desc:          form.Desc,
		Content:       form.Content,
		ModifiedBy:    form.ModifiedBy,
		CoverImageUrl: form.CoverImageUrl,
		State:         form.State,
	}

	exists, err := articleService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		g.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	tagService := tag_service.Tag{ID: form.TagID}
	exists, err = tagService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exists {
		g.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = articleService.Edit()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EDIT_ARTICLE_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)

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
// @Router /api/v1/articles/{id} [delete]
func DeleteArticle(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()
	g := app.Gin{C: c}

	valid := validation.Validation{}
	valid.Min(id, 1, "id").Message("ID必须大于0")

	if valid.HasErrors() {
		app.MakrErrors(valid.Errors)
		g.Response(http.StatusOK, e.INVALID_PARAMS, nil)
		return
	}

	articleService := article_service.Article{ID: id}
	exists, err := articleService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_CHECK_EXIST_ARTICLE_FAIL, nil)
		return
	}
	if !exists {
		g.Response(http.StatusOK, e.ERROR_NOT_EXIST_ARTICLE, nil)
		return
	}

	err = articleService.Delete()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_DELETE_ARTICLE_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)
}
