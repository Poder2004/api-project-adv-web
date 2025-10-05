package handlers

import (
	models "api-game/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ... GetProfileHandler (ถ้ามี) ...

// UpdateProfileHandler จัดการการอัปเดตข้อมูลโปรไฟล์
func EditProfileHandler(c *gin.Context, db *gorm.DB) {
	// 1. สร้าง struct เพื่อรับข้อมูลจาก JSON body
	var input struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		// เพิ่ม field อื่นๆ ที่ต้องการให้แก้ไขได้ที่นี่
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	// 2. ดึง user_id จาก token ที่ middleware แปะมาให้
	// **หมายเหตุ:** เราจะถือว่ามี middleware ที่ทำการ decode token และเก็บ user_id ไว้ใน context แล้ว
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 3. ค้นหา user เดิมในฐานข้อมูล
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 4. อัปเดตข้อมูลเฉพาะ field ที่มีการส่งค่ามา
	// ใช้ .Model(&user) เพื่อระบุว่าจะอัปเดต record ไหน
	// ใช้ .Updates() เพื่ออัปเดตเฉพาะ field ที่ไม่เป็น zero-value
	if err := db.Model(&user).Updates(models.User{Username: input.Username, Email: input.Email}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile", "details": err.Error()})
		return
	}

	// 5. ส่งข้อมูลที่อัปเดตแล้วกลับไป (ยกเว้นรหัสผ่าน)
	user.Password = "" // ไม่ส่งรหัสผ่านกลับไป
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
		"user":    user,
	})
}
