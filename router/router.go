package routers

import (
	handlers "api-game/handler" // 👈 แก้จาก "Handler" เป็น "handler" (h ตัวเล็ก)
	"api-game/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter ตั้งค่า routes ของแอป
func SetupRouter(r *gin.Engine, db *gorm.DB) {
	r.Static("/uploads", "./uploads")

	// --- ส่วนของ CORS Config (ถูกต้องแล้ว ไม่ต้องแก้ไข) ---
	config := cors.Config{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}

	r.Use(cors.New(config))
	// --- สิ้นสุดส่วนของ CORS Config ---

	// ---------- Public Routes ----------
	r.POST("/register", func(c *gin.Context) {
		handlers.RegisterHandler(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.LoginHandler(c, db)
	})

	// Protected routes
	protected := r.Group("/api")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.PUT("/updateprofile", func(c *gin.Context) {
			handlers.EditProfileHandler(c, db)
		})
		protected.GET("/profile", func(c *gin.Context) {
			handlers.GetProfileHandler(c, db)
		})
	}
}
