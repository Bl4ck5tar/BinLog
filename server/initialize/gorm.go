package initialize

import (
	"BinLog/server/global"
	"gorm.io/driver/mysql"
	"os"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

//InitGorm 初始化并返回一个使用MYSQL 配置的GORM 数据库连接
func InitGorm() *gorm.DB {
	mysqlConfig := global.Config.Mysql

	//使用给定的DSN(Data Source Name) 和日志级别打开MYSQ 数据库连接
	db, err := gorm.Open(mysql.Open(mysqlConfig.Dsn()), &gorm.Config{
		Logger: logger.Default.LogMode(mysqlConfig.LogLevel()),
	})
	if err != nil {
		global.Log.Error("Failed to connect to MySQL:", zap.Error(err))
		os.Exit(1)
	}
	//获取底层的 SQL 数据库连接对象
	sqlDB, _ := db.DB()
	//设置数据库连接池中的最大空闲连接数
	sqlDB.SetMaxIdleConns(mysqlConfig.MaxIdleConns)
	//设置数据库的最大打开连接数
	sqlDB.SetMaxOpenConns(mysqlConfig.MaxOpenConns)

	return db
}