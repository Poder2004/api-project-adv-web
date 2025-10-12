package handlers

import (
	"api-game/model"
	"errors"
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

// handler/admin/game_handler.go (ต่อท้ายไฟล์เดิม)
func GetGameByIDHandler(c *gin.Context, db *gorm.DB) {
	idStr := c.Param("id")
	var game model.Game

	// preload category มาด้วย
	if err := db.Preload("Category").First(&game, "game_id = ?", idStr).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"status": "error", "message": "game not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "db error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "game fetched",
		"data":    game,
	})
}

// GetMyOrdersHandler ดึงประวัติการซื้อเกมของผู้ใช้ที่กำลังล็อกอินอยู่
func GetMyOrdersHandler(c *gin.Context, db *gorm.DB) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	var orders []model.Order

	// --- 👇 [แก้ไข Preload ทั้งหมดตรงนี้] ---
	err := db.Preload("User").
		Preload("DiscountCode").
		Preload("OrderDetails").
		// 💡 จุดสำคัญ: Preload Game โดยใช้ Unscoped() เพื่อให้ดึงเกมที่ถูก soft delete มาด้วย
		Preload("OrderDetails.Game", func(db *gorm.DB) *gorm.DB {
			return db.Unscoped()
		}).
		// Preload Category ของ Game ที่ดึงมาแล้ว
		Preload("OrderDetails.Game.Category").
		Where("user_id = ?", userID).
		Order("order_date desc").
		Find(&orders).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user orders"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "success", "data": orders})
}
