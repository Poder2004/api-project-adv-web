package routers

import (
	handlers "api-game/Handler"
	// "api-game/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter ตั้งค่า routes ของแอป
func SetupRouter(r *gin.Engine, db *gorm.DB) {

	// ---------- Public Routes ----------
	r.POST("/register", func(c *gin.Context) {
		handlers.RegisterHandler(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.LoginHandler(c, db)
	})

	// // ---------- Protected Routes (ต้องมี token) ----------
	// api := r.Group("/api")
	// api.Use(middleware.AuthMiddleware()) // ใช้ middleware ตรวจ token ก่อนเข้า
	// {
	// 	api.GET("/profile", func(c *gin.Context) {
	// 		handlers.ProfileHandler(c, db) // ตัวอย่าง endpoint ดึงข้อมูลผู้ใช้
	// 	})

	// 	api.GET("/wallet", func(c *gin.Context) {
	// 		handlers.WalletHandler(c, db) // endpoint สมมติ
	// 	})
	// }
}
