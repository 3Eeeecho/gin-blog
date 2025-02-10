package models

import (
	"fmt"

	"github.com/3Eeeecho/go-gin-example/pkg/logging"
	"github.com/3Eeeecho/go-gin-example/pkg/setting"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// 定义全局变量 db，用于存储数据库连接
var db *gorm.DB

// Model 是一个基础模型，包含所有模型共有的字段
type Model struct {
	ID         int `gorm:"primary_key" json:"id" `            // 主键 ID
	CreatedOn  int `json:"created_on" gorm:"autoCreateTime"`  // 创建时间
	ModifiedOn int `json:"modified_on" gorm:"autoUpdateTime"` // 修改时间
	DeletedOn  int `json:"delete_on" gorm:"softDelete"`
}

func SetUp() {
	var (
		err                                               error
		dbType, dbName, user, password, host, tablePrefix string
	)

	// 读取数据库配置
	dbType = setting.DatabaseSetting.Type             // 数据库类型（如 mysql）
	dbName = setting.DatabaseSetting.Name             // 数据库名称
	user = setting.DatabaseSetting.User               // 数据库用户名
	password = setting.DatabaseSetting.Password       // 数据库密码
	host = setting.DatabaseSetting.Host               // 数据库主机地址
	tablePrefix = setting.DatabaseSetting.TablePrefix // 表前缀

	// 使用 GORM 打开数据库连接
	db, err = gorm.Open(dbType, fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		user,
		password,
		host,
		dbName))

	if err != nil {
		logging.Info(err) // 如果连接失败，记录错误日志
	}

	// 设置 GORM 的表名处理函数
	// 默认表名会加上配置中的表前缀
	gorm.DefaultTableNameHandler = func(db *gorm.DB, defaultTableName string) string {
		return tablePrefix + defaultTableName
	}

	// 设置 GORM 使用单数表名（默认情况下，GORM 会将结构体名称转换为复数形式作为表名）
	db.SingularTable(true)

	// 开启 GORM 的调试模式，打印所有执行的 SQL 语句
	db.LogMode(true)

	// 设置数据库连接池的最大空闲连接数
	db.DB().SetMaxIdleConns(10)

	// 设置数据库连接池的最大打开连接数
	db.DB().SetMaxOpenConns(100)
}

// CloseDB 关闭数据库连接
func CloseDB() {
	defer db.Close()
}
