package main

import (
	"api-game/database"
	routers "api-game/router"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors" 
)

func main() {
	// ตั้งค่าการเชื่อมต่อฐานข้อมูล
	db, err := database.SetupDatabaseConnection()
	if err != nil {
		panic("Failed to connect to the database")
	}

	// สร้าง Gin router
	r := gin.Default()
	
	r.Use(cors.Default())

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
