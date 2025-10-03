package router

import (
	"api-game/handler" // Import package handler
	"net/http"
)

// SetupRouter ทำหน้าที่ตั้งค่าและคืนค่า router ทั้งหมดของแอปพลิเคชัน
func SetupRouter() *http.ServeMux {
	// ใช้ ServeMux ของ standard library (สามารถเปลี่ยนเป็น Gin, Echo, Chi ได้ในอนาคต)
	mux := http.NewServeMux()

	// กำหนดเส้นทาง /health ให้ไปเรียกใช้ HealthCheckHandler
	mux.HandleFunc("/health", handler.HealthCheckHandler)

	// --- เพิ่ม routes อื่นๆ ได้ที่นี่ ---
	// mux.HandleFunc("/users", handler.CreateUserHandler)
	// mux.HandleFunc("/users/{id}", handler.GetUserHandler)

	return mux
}