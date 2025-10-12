package handlersadmin

import (
	"api-game/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func GetAllUsersHandler(c *gin.Context, db *gorm.DB) {
	var users []model.User

	// --- 👇 [ส่วนที่แก้ไข] ---
	// 1. เพิ่ม .Where("role = ?", "member") เพื่อกรองข้อมูล
	if err := db.Where("role = ?", "member").Order("user_id desc").Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch users"})
		return
	}
	// --- 👆 [สิ้นสุดส่วนที่แก้ไข] ---

	// 2. ล้างค่ารหัสผ่านก่อนส่งกลับไปเพื่อความปลอดภัย
	for i := range users {
		users[i].Password = ""
	}

	// 3. ส่งข้อมูลกลับไป
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Users fetched successfully",
		"data":    users,
	})
}

// GetUserByIDHandler ดึงข้อมูลผู้ใช้คนเดียวตาม ID (สำหรับ Admin)
func GetUserByIDHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var user model.User
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	user.Password = "" // ไม่ส่งรหัสผ่านกลับไป
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": user})
}

// GetUserOrdersHandler ดึงประวัติการซื้อเกมทั้งหมดของผู้ใช้ (สำหรับ Admin)
func GetUserOrdersHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var orders []model.Order

	// --- 👇 [แก้ไข Preload ตรงนี้] ---
	// ใช้ Preload เพื่อดึงข้อมูล OrderDetails และ Game ที่ซ้อนอยู่ข้างในมาด้วย
	err := db.Preload("OrderDetails").
		// 💡 จุดสำคัญ: Preload Game โดยใช้ Unscoped() เช่นกัน
		Preload("OrderDetails.Game", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		Where("user_id = ?", id).
		Order("order_date desc").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user orders"})
		return
	}

	if len(orders) == 0 { // ใช้ len() จะปลอดภัยกว่าการเช็ค nil โดยตรง
		orders = []model.Order{}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": orders})
}

// GetUserWalletHistoryHandler ดึงประวัติการเติมเงินทั้งหมดของผู้ใช้ (สำหรับ Admin)
func GetUserWalletHistoryHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var history []model.WalletHistory
	if err := db.Where("user_id = ?", id).Order("transaction_date desc").Find(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch wallet history"})
		return
	}

	if history == nil {
		history = []model.WalletHistory{}
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": history})
}
