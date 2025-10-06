package handlers// หรือ handlersadmin ตามโครงสร้างของคุณ

import (
	"api-game/model" // แก้ไข path ตามโปรเจกต์ของคุณ
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetAllCategories ดึงข้อมูลประเภทเกมทั้งหมด
func GetAllCategories(c *gin.Context, db *gorm.DB) {
	var categories []model.Category

	// ค้นหา categories ทั้งหมดในฐานข้อมูล
	if err := db.Find(&categories).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch categories"})
		return
	}

	c.JSON(http.StatusOK, categories) // ส่งข้อมูลกลับไปเป็น JSON array
}