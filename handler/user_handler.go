package handlers

import (
	// ตรวจสอบ path ของ model ให้ถูกต้อง
	models "api-game/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt" // 👈 1. Import bcrypt สำหรับเข้ารหัส
	"gorm.io/gorm"
)

// EditProfileHandler จัดการการอัปเดตข้อมูลโปรไฟล์
func EditProfileHandler(c *gin.Context, db *gorm.DB) {
	// 2. อัปเดต struct ที่รับข้อมูลเข้ามา
	var input struct {
		Username     string `json:"username"`
		Email        string `json:"email"`
		ImageProfile string `json:"imageProfile"`
		Password     string `json:"password"` // รับรหัสผ่านใหม่ (อาจจะเป็นค่าว่าง)
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. เตรียมข้อมูลที่จะอัปเดต
	updateData := models.User{
		Username:     input.Username,
		Email:        input.Email,
		ImageProfile: input.ImageProfile,
	}

	// 4. ตรวจสอบว่ามีการส่งรหัสผ่านใหม่มาหรือไม่
	if input.Password != "" {
		// ถ้ามี ให้ทำการ hash รหัสผ่านใหม่
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		// เพิ่มรหัสผ่านที่เข้ารหัสแล้วเข้าไปในข้อมูลที่จะอัปเดต
		updateData.Password = string(hashedPassword)
	}

	// 5. อัปเดตข้อมูลในฐานข้อมูล
	if err := db.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile", "details": err.Error()})
		return
	}

	// ส่งข้อมูลที่อัปเดตแล้วกลับไป
	user.Password = "" // ไม่ส่งรหัสผ่านกลับไป
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
		"user":    user,
	})
}
