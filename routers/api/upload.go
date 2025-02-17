package api

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/pkg/app"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/upload"
	"github.com/gin-gonic/gin"
)

// @Summary Import Image
// @Produce  json
// @Param image formData file true "Image File"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /api/v1/tags/import [post]
func UpLoadImage(c *gin.Context) {
	g := app.Gin{C: c}

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		g.Response(http.StatusOK, e.ERROR, nil)
		return
	}

	if image == nil {
		g.Response(http.StatusBadRequest, e.INVALID_PARAMS, nil)
		return
	}

	imageName := upload.GetImageName(image.Filename)
	fullPath := upload.GetImageFullPath()
	savePath := upload.GetImagePath()
	src := fullPath + imageName

	if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
		g.Response(http.StatusBadRequest, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	err = upload.CheckImage(fullPath)
	if err != nil {
		logging.Warn(err)
		g.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	if err = c.SaveUploadedFile(image, src); err != nil {
		logging.Warn(err)
		g.Response(http.StatusInternalServerError, e.ERROR_UPLOAD_CHECK_IMAGE_FAIL, nil)
		return
	}

	g.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"image_url":      upload.GetImageFullUrl(imageName),
		"image_save_url": savePath + imageName,
	})
}
