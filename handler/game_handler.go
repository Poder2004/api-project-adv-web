package handlers

import (
	"api-game/model"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// SearchHandler จัดการการค้นหาเกมตามชื่อ, ประเภท, หรือ ID
func SearchHandler(c *gin.Context, db *gorm.DB) {
	query := c.Query("q")

	if query == "" {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "Search query is empty",
			"data":    []model.Game{},
		})
		return
	}

	var games []model.Game
	searchQuery := "%" + query + "%"

	gameID, err := strconv.Atoi(query)
	if err != nil {
		gameID = 0 // ไม่ใช่ตัวเลข
	}

	// vvvvvvvvvv --- แก้ไขครั้งสุดท้าย --- vvvvvvvvvv
	result := db.Joins("Category").
		Preload("Category").
		// เปลี่ยนจาก "Category" (Double Quote) เป็น `Category` (Backtick)
		// เพื่อให้เข้ากับ Syntax ของ MySQL
		Where("games.title LIKE ? OR `Category`.category_name LIKE ? OR games.game_id = ?", searchQuery, searchQuery, gameID).
		Find(&games)
	// ^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^^

	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Database error: " + result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Games found",
		"data":    games,
	})
}
