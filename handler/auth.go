package handlers

import (
	models "api-game/model"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"github.com/golang-jwt/jwt/v5"
)

// RegisterHandler รับคำขอสมัครสมาชิก
func RegisterHandler(c *gin.Context, db *gorm.DB) {
	var input models.User

	// อ่านข้อมูลจาก JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// เข้ารหัสรหัสผ่าน
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	// เตรียมข้อมูลสำหรับบันทึก
	user := models.User{
		Username:     input.Username,
		Email:        input.Email,
		Password:     string(hashedPassword),
		ImageProfile: input.ImageProfile,
		Role:         "member",
		Wallet:       0.00,
	}

	// บันทึกลงฐานข้อมูล
	if err := db.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "User successfully created",
		"user_id": user.UserID,
	})
}

var jwtSecret = []byte(os.Getenv("JWT_SECRET")) // ตั้งใน .env เช่น JWT_SECRET=mysecretkey

func LoginHandler(c *gin.Context, db *gorm.DB) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	// ค้นหาผู้ใช้จาก username
	var user models.User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		return
	}

	// ตรวจรหัสผ่าน
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // หมดอายุใน 24 ชม.
	})

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "Failed to create token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Login successful",
		"token":   tokenString,
		"user": gin.H{
			"user_id":      user.UserID,
			"username":     user.Username,
			"email":        user.Email,
			"role":         user.Role,
			"wallet":       user.Wallet,
			"ImageProfile": user.ImageProfile,
		},
	})
}


func Profile(c *gin.Context, db *gorm.DB) {
	userIDStr := c.Query("user_id") // ดึง user_id จาก query string
	userID, err := strconv.ParseUint(userIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user_id"})
		return
	}

	var user models.User

	sql := "SELECT username, email FROM users WHERE user_id = ?"

	if err := db.Raw(sql, userID).Scan(&user).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	// --- สิ้นสุดส่วนที่แก้ไข ---

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "profile found",
		"user": gin.H{
			"username": user.Username,
			"email":    user.Email,
		},
	})
}
