package main

import (
	"log"
	"mxshop_srvs/user_srv/model"
	"os"
	"time"

	"gorm.io/gorm/schema"

	"gorm.io/gorm"

	"gorm.io/driver/mysql"
	"gorm.io/gorm/logger"
)

func main() {
	dsn := "root:fshing@tcp(localhost:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log Level
			Colorful:      true,        // 禁用彩色打印
		})

	db, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			Logger: newLogger,
			NamingStrategy: schema.NamingStrategy{
				SingularTable: true, // 不使用 gorm 创建的表名
			},
		})

	if err != nil {
		log.Panic(err)
	}

	_ = db.AutoMigrate(&model.User{})
}
