package main

import (
	"api-game/database" // Import package ที่สร้างขึ้น
	"api-game/router"   // Import package ที่สร้างขึ้น
	"log"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	// 1. โหลดค่าจากไฟล์ .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("เกิดข้อผิดพลาดในการโหลดไฟล์ .env")
	}

	// 2. เรียกใช้ฟังก์ชัน InitDB จาก package database
	database.InitDB()

	// 3. เรียกใช้ฟังก์ชัน SetupRouter จาก package router
	r := router.SetupRouter()

	// 4. รัน Web Server ที่ Port 8080 โดยใช้ router ที่ตั้งค่าไว้
	log.Println("เซิร์ฟเวอร์กำลังทำงานที่พอร์ต 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}