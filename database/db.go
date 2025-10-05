package database

import (
	"fmt"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func SetupDatabaseConnection() (*gorm.DB, error) {
	dsn := "66011212129:191047@tcp(202.28.34.210:3309)/db66011212129?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Printf("❌ ไม่สามารถเชื่อมต่อฐานข้อมูลได้: %v", err)
		return nil, err
	}

	fmt.Println("✅ เชื่อมต่อฐานข้อมูลสำเร็จผ่าน GORM!")
	return db, nil
}
