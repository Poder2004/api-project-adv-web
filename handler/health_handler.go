package handler

import (
	"api-game/database" // Import package database ที่เราสร้าง
	"api-game/model"    // Import package model ที่เราสร้าง
	"encoding/json"
	"log"
	"net/http"
)

// HealthCheckHandler คือฟังก์ชันที่จัดการ request ที่เข้ามาที่ /health
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// ใช้ตัวแปร DB จาก package database
	err := database.DB.Ping()

	var response model.HealthResponse

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response = model.HealthResponse{Status: "error", Database: "disconnected"}
		log.Printf("Health check ล้มเหลว: ไม่สามารถ ping ฐานข้อมูลได้: %v", err)
	} else {
		w.WriteHeader(http.StatusOK)
		response = model.HealthResponse{Status: "ok", Database: "connected"}
	}

	json.NewEncoder(w).Encode(response)
} 