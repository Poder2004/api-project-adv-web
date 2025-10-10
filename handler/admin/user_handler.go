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

// GetUserByIDHandler ดึงข้อมูลผู้ใช้คนเดียวตาม ID
func GetUserByIDHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง "id" จาก URL parameter
	id := c.Param("id")

	var user model.User
	// 2. ค้นหา user ในฐานข้อมูลด้วย ID ที่ได้รับมา
	if err := db.First(&user, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. ไม่ส่งรหัสผ่านกลับไปเพื่อความปลอดภัย
	user.Password = ""

	// 4. ส่งข้อมูลผู้ใช้กลับไป
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   user,
	})
}

// GetUserOrdersHandler ดึงประวัติการซื้อเกมทั้งหมดของผู้ใช้
func GetUserOrdersHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง "id" จาก URL parameter
	id := c.Param("id")
	var orders []model.Order

	// 2. ค้นหา Orders ทั้งหมดที่ user_id ตรงกัน
	//    ใช้ Preload("OrderDetails.Game") เพื่อดึงข้อมูลรายละเอียดการสั่งซื้อและข้อมูลเกมที่ซ้อนอยู่ข้างในมาด้วย
	if err := db.Preload("OrderDetails.Game").Where("user_id = ?", id).Order("order_date desc").Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user orders"})
		return
	}

	// ถ้าไม่พบข้อมูล GORM จะคืนค่าเป็น slice ว่างๆ (ไม่ใช่ error)
	if orders == nil {
		orders = []model.Order{}
	}

	// 3. ส่งข้อมูลประวัติการซื้อทั้งหมดกลับไป
	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   orders,
	})
}