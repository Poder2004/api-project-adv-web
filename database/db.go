package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

// DB คือตัวแปร global สำหรับ connection pool ที่จะให้ package อื่นเรียกใช้
var DB *sql.DB

// InitDB คือฟังก์ชันสำหรับเริ่มต้นการเชื่อมต่อฐานข้อมูล
func InitDB() {
	// สร้าง Data Source Name (DSN) string
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	// เชื่อมต่อฐานข้อมูลและเก็บ connection pool ไว้ในตัวแปร DB
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("ไม่สามารถเปิดการเชื่อมต่อฐานข้อมูลได้: %v", err)
	}

	// ตรวจสอบว่าการเชื่อมต่อสำเร็จจริง
	err = DB.Ping()
	if err != nil {
		log.Fatalf("ไม่สามารถเชื่อมต่อฐานข้อมูลได้: %v", err)
	}

	fmt.Println("เชื่อมต่อฐานข้อมูลสำเร็จ!")
}
