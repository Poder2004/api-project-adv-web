package routers

import (
	handlers "api-game/handler"
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
