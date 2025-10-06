package handlersadmin // หรือ package handler ตามโครงสร้างของคุณ

import (
	"api-game/model"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateGame handles the creation of a new game.
func CreateGame(c *gin.Context, db *gorm.DB) {
	// --- 1. รับข้อมูลจาก Form (Multipart/form-data) ---
	title := c.PostForm("title")
	description := c.PostForm("description")
	priceStr := c.PostForm("price")
	categoryIDStr := c.PostForm("category_id")

	// --- 2. ตรวจสอบข้อมูลเบื้องต้น ---
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Game title is required"})
		return
	}

	// --- 3. แปลงข้อมูล String เป็น Number ---
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price format"})
		return
	}

	categoryID, err := strconv.Atoi(categoryIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID format"})
		return
	}

	// --- 4. จัดการไฟล์ที่อัปโหลด ---
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image upload failed: " + err.Error()})
		return
	}

	// สร้างชื่อไฟล์ใหม่ที่ไม่ซ้ำกัน โดยใช้ timestamp
	extension := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	filePath := "uploads/" + newFileName

	// บันทึกไฟล์ไปยัง Foleder 'uploads'
	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
		return
	}

	// --- 5. สร้าง Object Game เพื่อบันทึกลง Database ---
	game := model.Game{
		Title:       title,
		Description: description,
		Price:       price,
		CategoryID:  uint(categoryID),
		ImageGame:   filePath, // เก็บแค่ path ของไฟล์
		ReleaseDate: time.Now(),
	}


	result := db.Create(&game)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create game: " + result.Error.Error()})
		return
	}

	// --- 7. ส่งผลลัพธ์กลับไป ---
	c.JSON(http.StatusCreated, gin.H{
		"message": "Game created successfully",
		"data":    game,
	})
}