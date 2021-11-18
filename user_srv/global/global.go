package global

import (
	"log"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func init() {
	dsn := "root:fshing@tcp(110.40.150.61:3306)/shop_user_srv?charset=utf8mb4&parseTime=True&loc=Local"

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold: time.Second, // 慢 SQL 阈值
			LogLevel:      logger.Info, // Log Level
			Colorful:      true,        // 禁用彩色打印
		})

	var err error
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			Logger: newLogger,
			//NamingStrategy: schema.NamingStrategy{
			//	SingularTable: true, // 不使用 gorm 创建的表名
			//},
		})

	if err != nil {
		log.Panic(err)
	}

}
