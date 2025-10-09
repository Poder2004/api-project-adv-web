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
