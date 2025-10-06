package main

import (
	"api-game/database"
	routers "api-game/router"
	"log"
	"os"

	"github.com/gin-contrib/cors" // import ตัวนี้
	"github.com/gin-gonic/gin"
)

func main() {
	// ตั้งค่าการเชื่อมต่อฐานข้อมูล
	db, err := database.SetupDatabaseConnection()
	if err != nil {
		panic("Failed to connect to the database")
	}

	// สร้าง Gin router
	r := gin.Default()

	// --- เริ่มส่วนที่แก้ไข CORS ---
	// กำหนดค่าคอนฟิกของ CORS ด้วยตัวเอง
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"}, // <<-- จุดสำคัญ! เพิ่ม Authorization ตรงนี้
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}

	// ใช้ middleware ของ CORS ที่เราสร้างขึ้นมา
	r.Use(cors.New(config))
	// --- จบส่วนที่แก้ไข CORS ---

	// เรียกใช้ routes จาก package routers
	routers.SetupRouter(r, db)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // default เวลา run local
	}
	log.Printf("🚀 Server running on port %s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatal("❌ Failed to start server: ", err)
	}
}
