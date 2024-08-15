package models

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	dsn := "root:123456789@tcp(10.12.135.79:3306)/test?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("无法连接数据库: %v", err)
	}
	log.Println("数据库连接成功")
	return db
}

var DB = Init()
