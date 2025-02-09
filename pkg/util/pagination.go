package util

import (
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/gin-gonic/gin"
	"github.com/unknwon/com"
)

// GetPage 函数用于从请求中获取分页参数，并计算当前页的偏移量
func GetPage(c *gin.Context) int {
	result := 0 // 默认偏移量为 0

	// 从请求的查询参数中获取 "page" 参数，并将其转换为整数
	// c.Query("page") 获取 URL 中的查询参数，例如：/articles?page=2
	// com.StrTo 将字符串转换为其他类型，这里转换为整数
	page, _ := com.StrTo(c.Query("page")).Int()

	// 如果 page 大于 0，则计算偏移量
	// 偏移量的计算公式为：(page - 1) * 每页显示的条数
	if page > 0 {
		result = (page - 1) * setting.PageSize
	}

	// 返回计算后的偏移量
	return result
}
