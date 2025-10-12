// file: database/database.go

package database

import (
	"fmt"
	"log"
	"os"   //  <-- 1. Import เพิ่ม
	"time" //  <-- 2. Import เพิ่ม

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger" //  <-- 3. Import เพิ่ม
)

func SetupDatabaseConnection() (*gorm.DB, error) {
	dsn := "66011212129:191047@tcp(202.28.34.210:3309)/db66011212129?charset=utf8mb4&parseTime=True&loc=Local"

	// --- 👇 [เพิ่มโค้ดส่วนนี้เข้าไป] ---
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // กำหนดให้ log แสดงผลใน Console
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // แจ้งเตือนถ้า query ช้ากว่า 0.2 วินาที
			LogLevel:                  logger.Info,            //  <-- ตั้งเป็น Info เพื่อให้แสดง SQL ทุกคำสั่ง
			IgnoreRecordNotFoundError: true,                   // ไม่ต้อง log error "record not found"
			Colorful:                  true,                   // แสดง log แบบมีสีสันให้อ่านง่าย
		},
	)
	// --- 👆 [สิ้นสุดส่วนที่เพิ่ม] ---

	// --- 👇 [แก้ไขตรงนี้ โดยเพิ่ม Logger เข้าไปใน Config] ---
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: newLogger, //  <-- ✅ เพิ่มบรรทัดนี้
	})

	if err != nil {
		log.Printf("❌ ไม่สามารถเชื่อมต่อฐานข้อมูลได้: %v", err)
		return nil, err
	}

	fmt.Println("✅ เชื่อมต่อฐานข้อมูลสำเร็จผ่าน GORM!")
	return db, nil
}
