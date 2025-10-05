package handlers

import (
	models "api-game/model"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)
// GetProfileHandler ดึงข้อมูลโปรไฟล์ของผู้ใช้ที่ล็อกอินอยู่
func GetProfileHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง user_id จาก token ที่ middleware แปะมาให้
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// 2. ค้นหา user ในฐานข้อมูลด้วย user_id
	var user models.User
	if err := db.First(&user, userID).Error; err != nil {
		// GORM จะ return 'record not found' error ถ้าไม่เจอ
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// 3. ไม่ส่งรหัสผ่านกลับไปเพื่อความปลอดภัย
	user.Password = ""

	// 4. ส่งข้อมูลโปรไฟล์ทั้งหมดกลับไปเป็น JSON
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile fetched successfully",
		"user":    user,
	})
}


func EditProfileHandler(c *gin.Context, db *gorm.DB) {
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

	updateData := make(map[string]interface{})

	// --- 👇 ส่วนที่แก้ไขตามแนวคิดของคุณ + แก้ไข Key ---

	// 1. ตรวจสอบ Username: ถ้ามีค่าใหม่ และไม่เหมือนค่าเดิม ถึงจะอัปเดต
	if username := c.PostForm("username"); username != "" && username != user.Username {
		updateData["username"] = username // key "username" ถูกต้อง
	}

	// 2. ตรวจสอบ Email: ถ้ามีค่าใหม่ และไม่เหมือนค่าเดิม ถึงจะอัปเดต
	if email := c.PostForm("email"); email != "" && email != user.Email {
		updateData["email"] = email // key "email" ถูกต้อง
	}

	// 3. จัดการไฟล์รูปภาพ
	file, err := c.FormFile("imageProfile")
	if err == nil {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		filePath := "uploads/" + newFileName
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save uploaded file"})
			return
		}
		// 💥 แก้ไข Key ให้ตรงกับชื่อคอลัมน์ใน DB
		updateData["image_profile"] = filePath
	}

	// 4. จัดการรหัสผ่านใหม่ (ไม่ต้องเทียบของเก่า)
	if password := c.PostForm("password"); password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash new password"})
			return
		}
		updateData["password"] = string(hashedPassword) // key "password" ถูกต้อง
	}

	// 5. ตรวจสอบว่ามีข้อมูลให้อัปเดตหรือไม่
	if len(updateData) == 0 {
		// ถ้าไม่มีอะไรเปลี่ยนแปลงเลย ก็ส่ง response กลับไปว่าสำเร็จ แต่ไม่มีอะไรอัปเดต
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No changes detected",
			"user":    user,
		})
		return
	}

	// 6. อัปเดตข้อมูลลงฐานข้อมูล
	if err := db.Model(&user).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update profile", "details": err.Error()})
		return
	}

	// ดึงข้อมูลล่าสุดหลังจากอัปเดตเพื่อส่งกลับ
	db.First(&user, userID)
	user.Password = ""
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Profile updated successfully",
		"user":    user,
	})
}
