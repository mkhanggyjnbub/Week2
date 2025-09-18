package Db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := "host=MKHANGDZ1ST user=root password=mkhang123abc dbname=taskmanagementsystem  port=5432 sslmode=disable TimeZone=Asia/Ho_Chi_Minh"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		go panic("Không thể kết nối database: " + err.Error())
	}

	fmt.Println("✅ Kết nối PostgreSQL thành công!")
	DB = database
}
