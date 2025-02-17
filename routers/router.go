package routers

import (
	_ "github.com/3Eeeecho/go-gin-example/docs"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/3Eeeecho/go-gin-example/pkg/upload"
	"github.com/3Eeeecho/go-gin-example/routers/api"
	v1 "github.com/3Eeeecho/go-gin-example/routers/api/v1"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"
)

// InitRouter 初始化Gin路由
// @Description 配置并返回一个Gin路由实例，注册API路由、Swagger文档和中间件
// @Return *gin.Engine Gin引擎实例
func InitRouter() *gin.Engine {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	gin.SetMode(setting.ServerSetting.RunMode)

	r.Static("/upload/images", upload.GetImageFullPath())

	r.GET("/auth", api.GetAuth)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.POST("/upload", api.UpLoadImage)

	apiv1 := r.Group("/api/v1")
	//apiv1.Use(jwt.JWT())
	{
		//获取标签列表
		apiv1.GET("/tags", v1.GetTags)
		//新建标签
		apiv1.POST("/tags", v1.AddTag)
		//更新指定标签
		apiv1.PUT("/tags/:id", v1.EditTag)
		//删除指定标签
		apiv1.DELETE("/tags/:id", v1.DeleteTag)

		//获取文章列表
		apiv1.GET("/articles", v1.GetArticles)
		//获取指定文章
		apiv1.GET("/articles/:id", v1.GetArticle)
		//新建文章
		apiv1.POST("/articles", v1.AddArticle)
		//更新指定文章
		apiv1.PUT("/articles/:id", v1.EditArticle)
		//删除指定文章
		apiv1.DELETE("/articles/:id", v1.DeleteArticle)
	}

	return r
}
