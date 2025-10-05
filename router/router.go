package routers

import (
	handlers "api-game/Handler"
	"api-game/middleware"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter ตั้งค่า routes ของแอป
func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.Static("/uploads", "./uploads")
	// ---------- Public Routes ----------
	r.POST("/register", func(c *gin.Context) {
		handlers.RegisterHandler(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.LoginHandler(c, db)
	})

	// Protected routes (ต้อง login)
	// สร้าง Group สำหรับเส้นทางที่ต้องผ่าน middleware
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		// เพิ่มเส้นทาง PUT /profile ที่นี่
		protected.PUT("/profile", func(c *gin.Context) {
			handlers.EditProfileHandler(c, db)
		})

	}
}
