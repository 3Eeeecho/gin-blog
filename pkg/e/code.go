package e

// 定义全局常量，用于表示应用程序中的各种状态码
const (
	// HTTP 状态码
	SUCCESS        = 200 // 请求成功
	INVALID_PARAMS = 400 // 请求参数无效
	ERROR          = 500 // 服务器内部错误

	// 标签相关错误码
	ERROR_EXIST_TAG         = 10001 // 标签已存在
	ERROR_NOT_EXIST_TAG     = 10002 // 标签不存在
	ERROR_NOT_EXIST_ARTICLE = 10003 // 文章不存在

	// 认证和授权相关错误码
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001 // Token 验证失败
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002 // Token 已超时
	ERROR_AUTH_TOKEN               = 20003 // Token 生成或解析失败
	ERROR_AUTH                     = 20004 // 认证失败
)
