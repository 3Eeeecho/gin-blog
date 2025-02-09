package setting

import (
	"log"
	"time"

	"github.com/go-ini/ini"
)

var (
	Cfg *ini.File

	RunMode string

	HTTPPort     int
	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	PageSize  int
	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini': %v", err)
	}
	LoadBase()
	LoadServer()
	LoadApp()
}

// LoadBase 加载基础配置
func LoadBase() {
	// 从默认节（无节名部分）中获取 RUN_MODE 配置项
	// 如果不存在或为空，则使用默认值 "debug"
	RunMode = Cfg.Section("").Key("RunMode").MustString("debug")
}

// LoadServer 加载服务器相关配置
func LoadServer() {
	// 获取名为 "server" 的配置节
	sec, err := Cfg.GetSection("server")
	if err != nil {
		// 如果获取失败，记录错误并退出程序
		log.Fatalf("Fail to get section 'server': %v", err)
	}

	// 从 server 节中获取 HTTP_PORT 配置项
	// 如果不存在或为空，则使用默认值 8000
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)

	// 从 server 节中获取 READ_TIMEOUT 配置项
	// 如果不存在或为空，则使用默认值 60，并将其转换为 time.Duration 类型（单位为秒）
	ReadTimeout = time.Duration(sec.Key("READ_TIMEOUT").MustInt(60)) * time.Second

	// 从 server 节中获取 WRITE_TIMEOUT 配置项
	// 如果不存在或为空，则使用默认值 60，并将其转换为 time.Duration 类型（单位为秒）
	WriteTimeout = time.Duration(sec.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

// LoadApp 加载应用程序相关配置
func LoadApp() {
	// 获取名为 "app" 的配置节
	sec, err := Cfg.GetSection("app")
	if err != nil {
		// 如果获取失败，记录错误并退出程序
		log.Fatalf("Fail to get section 'app': %v", err)
	}

	// 从 app 节中获取 JWT_SECRET 配置项
	// 如果不存在或为空，则使用默认值 "!@)*#)!@U#@*!@!)"
	JwtSecret = sec.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")

	// 从 app 节中获取 PAGE_SIZE 配置项
	// 如果不存在或为空，则使用默认值 10
	PageSize = sec.Key("PAGE_SIZE").MustInt(10)
}
