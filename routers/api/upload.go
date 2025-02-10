package api

import (
	"net/http"

	"github.com/3Eeeecho/go-gin-example/models"
	"github.com/3Eeeecho/go-gin-example/pkg/e"
	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/upload"
	"github.com/gin-gonic/gin"
)

func UpLoadImage(c *gin.Context) {
	code := e.SUCCESS
	data := make(map[string]string)

	file, image, err := c.Request.FormFile("image")
	if err != nil {
		logging.Warn(err)
		code = e.ERROR
		response := models.NewResponse(code, e.GetMsg(code), data)
		c.JSON(http.StatusOK, response)
	}

	if image == nil {
		code = e.INVALID_PARAMS
	} else {
		imageName := upload.GetImageName(image.Filename)
		fullPath := upload.GetImageFullPath()
		savePath := upload.GetImagePath()

		src := fullPath + imageName
		if !upload.CheckImageExt(imageName) || !upload.CheckImageSize(file) {
			code = e.ERROR_UPLOAD_CHECK_IMAGE_FORMAT
		} else {
			err := upload.CheckImage(fullPath)
			if err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_CHECK_IMAGE_FAIL
			} else if err := c.SaveUploadedFile(image, src); err != nil {
				logging.Warn(err)
				code = e.ERROR_UPLOAD_SAVE_IMAGE_FAIL
			} else {
				data["image_url"] = upload.GetImageFullUrl(imageName)
				data["image_save_url"] = savePath + imageName
			}
		}
	}
	response := models.NewResponse(code, e.GetMsg(code), data)
	c.JSON(http.StatusOK, response)
}
