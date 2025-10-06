package handlers

import (
	"api-game/config" // <-- 1. Import package config ที่สร้างใหม่
	models "api-game/model"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// RegisterHandler รับคำขอสมัครสมาชิก
func RegisterHandler(c *gin.Context, db *gorm.DB) {
	// ... โค้ดส่วนนี้เหมือนเดิมทุกประการ ...
	username := c.PostForm("username")
	email := c.PostForm("email")
	password := c.PostForm("password")

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to encrypt password"})
		return
	}

	user := models.User{
		Username: username,
		Email:    email,
		Password: string(hashedPassword),
		Role:     "member",
		Wallet:   0.00,
	}

	file, err := c.FormFile("imageProfile")
	if err == nil {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		filePath := "uploads/" + newFileName
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			return
		}
		user.ImageProfile = filePath
	}

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


func LoginHandler(c *gin.Context, db *gorm.DB) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": err.Error()})
		return
	}

	var user models.User
	if err := db.Where("username = ?", input.Username).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "error", "message": "Invalid username or password"})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	// ✅ 2. เปลี่ยนมาใช้ Secret Key จาก config ที่โหลดมาแล้ว
	tokenString, err := token.SignedString([]byte(config.JWTSecret))
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