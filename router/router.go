package routers

import (
	handlers "api-game/handler"
	handlersadmin "api-game/handler/admin"
	handlersauser "api-game/handler/user"
	"api-game/middleware"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SetupRouter ตั้งค่า routes ของแอป
func SetupRouter(r *gin.Engine, db *gorm.DB) {

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello game")
	})

	r.Static("/uploads", "./uploads")

	// ---------- Public Routes ----------
	r.POST("/register", func(c *gin.Context) {
		handlers.RegisterHandler(c, db)
	})

	r.POST("/login", func(c *gin.Context) {
		handlers.LoginHandler(c, db)
	})

	r.GET("/api/categories", func(c *gin.Context) {
		handlers.GetAllCategories(c, db)
	})
	r.GET("/api/games/:id", func(c *gin.Context) {
			handlers.GetGameByIDHandler(c, db)
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
		protected.GET("/search", func(c *gin.Context) {
			handlers.SearchHandler(c, db)
		})

		// ใช้ POST /api/coupons/:did/claim
		protected.POST("/coupons/:did/claim", func(c *gin.Context) {
			handlersauser.ClaimCouponHandler(c, db)
		})
		// ✅ เพิ่ม Route นี้เข้าไปใหม่
		protected.GET("/my-coupons", func(c *gin.Context) {
			handlersauser.GetMyCouponsHandler(c, db)
		})
		//  เติมเงิน
		protected.POST("/wallet", func(c *gin.Context) {
			handlers.AddWalletHandler(c, db)
		})
		
	}

	admin := r.Group("/admin")
	admin.Use(middleware.AuthMiddleware())
	{
		admin.POST("/addgames", func(c *gin.Context) {
			handlersadmin.CreateGame(c, db)
		})
		admin.GET("/games", func(c *gin.Context) {
			handlersadmin.GetAllGamesHandler(c, db)
		})
		admin.POST("/coupons", func(c *gin.Context) {
			handlersadmin.CreateCouponHandler(c, db)
		})
		admin.GET("/allcoupons", func(c *gin.Context) {
			handlersadmin.GetAllCouponsHandler(c, db)
		})
		admin.GET("/alluser", func(c *gin.Context) {
			handlersadmin.GetAllUsersHandler(c, db)
		})

		admin.PUT("/games/:id", func(c *gin.Context) {
			handlersadmin.UpdateGameHandler(c, db)
		})

		admin.DELETE("/games/:id", func(c *gin.Context) {
			handlersadmin.DeleteGameHandler(c, db)
		})
	}
}
