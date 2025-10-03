package middleware

import (
	"log"
	"net/http"
)

// AuthMiddleware เป็นตัวอย่าง middleware สำหรับตรวจสอบสิทธิ์
// ในโค้ดจริง ส่วนนี้อาจจะมีการตรวจสอบ JWT Token หรือ Session
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("Middleware: ตรวจสอบ Request...")

		// ตรรกะการตรวจสอบสิทธิ์จะอยู่ตรงนี้
		// isAuthenticated := true // สมมติว่าตรวจสอบแล้วผ่าน
		// if !isAuthenticated {
		// 	http.Error(w, "Unauthorized", http.StatusUnauthorized)
		// 	return
		// }

		// ถ้าผ่าน ก็ส่งต่อไปให้ handler ตัวถัดไปทำงาน
		next.ServeHTTP(w, r)
	}
}