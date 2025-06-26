package v1

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/pkg/app"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/export"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/pkg/util"
	"github.com/3Eeeecho/go-gin-example/service/tag_service"
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
// @Success 200 {object} app.Response "返回标签列表和总数"
// @Router /api/v1/tags [get]
func GetTags(c *gin.Context) {
	g := app.Gin{C: c}
	name := c.Query("name")
	state := -1
	if arg := c.Query("state"); arg != "" {
		state = com.StrTo(arg).MustInt()
	}

	tagService := tag_service.Tag{
		Name:     name,
		State:    state,
		PageNum:  util.GetPage(c),
		PageSize: setting.AppSetting.PageSize,
	}

	tags, err := tagService.GetAll()
	if err != nil {
		logging.Info(err)
		g.Response(http.StatusInternalServerError, e.ERROR_GET_TAGS_FAIL, nil)
		return
	}

	count, err := tagService.Count()
	if err != nil {
		logging.Info(err)
		g.Response(http.StatusInternalServerError, e.ERROR_COUNT_TAG_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, map[string]interface{}{
		"lists": tags,
		"total": count,
	})
}

type AddTagForm struct {
	Name      string `form:"name" valid:"Required;MaxSize(100)"`
	CreatedBy string `form:"created_by" valid:"Required;MaxSize(100)"`
	State     int    `form:"state" valid:"Range(0,1)"`
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
// @Success 200 {object} app.Response "返回成功信息"
// @Failure 400 {object} app.Response "标签不存在"
// @Router /api/v1/tags [post]
func AddTag(c *gin.Context) {
	var (
		form AddTagForm
		g    = app.Gin{C: c}
	)

	httpCode, errCode := app.BindAndValue(c, &form)
	if errCode != e.SUCCESS {
		logging.Info(errCode)
		g.Response(httpCode, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{
		Name:      form.Name,
		CreatedBy: form.CreatedBy,
		State:     form.State,
	}

	exist, err := tagService.ExistByName()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if exist {
		g.Response(http.StatusOK, e.ERROR_EXIST_TAG, nil)
		return
	}

	err = tagService.Add()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_ADD_TAG_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)
}

type EditTagForm struct {
	ID         int    `form:"id" valid:"Required;Min(1)"`
	Name       string `form:"name" valid:"Required;MaxSize(100)"`
	ModifiedBy string `form:"modified_by" valid:"Required;MaxSize(100)"`
	State      int    `form:"state" valid:"Range(0,1)"`
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
// @Success 200 {object} app.Response "返回成功信息"
// @Failure 400 {object} app.Response "标签不存在"
// @Router /api/v1/tags/{id} [put]
func EditTag(c *gin.Context) {
	var (
		form = EditTagForm{ID: com.StrTo(c.Param("id")).MustInt()}
		g    = app.Gin{C: c}
	)

	httpCode, errCode := app.BindAndValue(c, &form)
	if errCode != e.SUCCESS {
		g.Response(httpCode, errCode, nil)
		return
	}

	tagService := tag_service.Tag{
		ID:         form.ID,
		Name:       form.Name,
		ModifiedBy: form.ModifiedBy,
		State:      form.State,
	}

	exist, err := tagService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		g.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Edit()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EDIT_TAG_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)
}

// DeleteTag 删除文章标签
// @Summary 删除文章标签
// @Description 删除指定标签
// @Tags 标签
// @Accept  json
// @Produce json
// @Param id path int true "标签ID"  // 必填参数，标签ID
// @Success 200 {object} app.Response "返回成功信息"
// @Failure 400 {object} app.Response "标签不存在"
// @Router /api/v1/tags/{id} [delete]
func DeleteTag(c *gin.Context) {
	id := com.StrTo(c.Param("id")).MustInt()

	g := app.Gin{C: c}
	valid := validation.Validation{}
	valid.Required(id, "id").Message("ID不能为空")

	if valid.HasErrors() {
		app.MakrErrors(valid.Errors)
		g.Response(http.StatusInternalServerError, e.INVALID_PARAMS, nil)
		return
	}

	tagService := tag_service.Tag{ID: id}
	exist, err := tagService.ExistByID()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_EXIST_TAG_FAIL, nil)
		return
	}

	if !exist {
		g.Response(http.StatusOK, e.ERROR_NOT_EXIST_TAG, nil)
		return
	}

	err = tagService.Delete()
	if err != nil {
		g.Response(http.StatusInternalServerError, e.ERROR_DELETE_TAG_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)
}

// ExportTag 导出标签数据
// @Summary 导出标签信息
// @Description 生成 Excel 文件并返回下载地址
// @Tags 标签管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param name formData string false "标签名称（可选）"
// @Param state formData int false "标签状态（可选），1=启用，0=禁用"
// @Success 200 {object} map[string]string "导出成功"
// @Failure 500 {object} app.Response "导出失败"
// @Router /api/tags/export [post]
func ExportTag(c *gin.Context) {
	g := app.Gin{C: c}
	name := c.PostForm("name")

	state := -1
	if arg := c.PostForm("state"); arg != "" {
		if state != -1 {
			state = com.StrTo(arg).MustInt()
		}
	}

	tagService := tag_service.Tag{
		Name:  name,
		State: state,
	}

	filename, err := tagService.Export()
	if err != nil {
		logging.Info(err)
		g.Response(http.StatusInternalServerError, e.ERROR_EXPORT_TAG_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"export_url":      export.GetExcelFullUrl(filename),
		"export_save_url": export.GetExcelPath() + filename,
	})
}

// ImportTag 导入标签数据
// @Summary 导入标签信息
// @Description 导入 Excel 文件并存储在本地
// @Tags 标签管理
// @Accept application/x-www-form-urlencoded
// @Produce json
// @Param file formData file true "文件"
// @Success 200 {object} map[string]string "导入成功"
// @Failure 500 {object} app.Response "导入失败"
// @Router /api/tags/import [post]
func ImportTag(c *gin.Context) {
	g := app.Gin{C: c}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		logging.Warn(err)
		g.Response(http.StatusInternalServerError, e.ERROR, nil)
		return
	}

	tagService := tag_service.Tag{}
	err = tagService.Import(file)
	if err != nil {
		logging.Warn(err)
		g.Response(http.StatusInternalServerError, e.ERROR_IMPORT_TAG_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, nil)
}
