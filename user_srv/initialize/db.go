package initialize

import (
	"fmt"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"mxshop_srvs/user_srv/global"
	"os"
	"time"
)

func InitDB() {
	config := global.ServerConfig.MysqlInfo
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User,config.Password,config.Host,config.Port,config.Name,
		)

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log Level
			Colorful:      true,        // 禁用彩色打印
		})

	var err error
	global.DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			Logger: newLogger,
			//NamingStrategy: schema.NamingStrategy{
			//	SingularTable: true, // 不使用 gorm 创建的表名
			//},
		})

	if err != nil {
		zap.S().Panic(err)
	}

}
