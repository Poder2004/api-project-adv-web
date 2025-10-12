package handlersadmin

import (
	"api-game/model"
	"fmt"
	"net/http"
	"os"
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

// GetAllGamesHandler handles fetching all games.
func GetAllGamesHandler(c *gin.Context, db *gorm.DB) {
	var games []model.Game

	// Use Preload("Category") to automatically fetch the associated category for each game.
	// This is known as "Eager Loading".
	if err := db.Preload("Category").Find(&games).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch games from database"})
		return
	}

	// If no games are found, GORM returns an empty slice, not an error.
	// It's good practice to ensure the slice is not nil for JSON marshalling.
	if games == nil {
		games = []model.Game{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Games fetched successfully",
		"data":    games,
	})
}

// UpdateGameHandler จัดการการอัปเดตข้อมูลเกม
func UpdateGameHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง ID ของเกมจาก URL
	id := c.Param("id")

	var game model.Game
	// 2. ค้นหาเกมเดิมในฐานข้อมูล
	if err := db.First(&game, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}

	// 3. สร้าง map เพื่อเก็บข้อมูลที่จะอัปเดตแบบไดนามิก
	updateData := make(map[string]interface{})

	// 4. ตรวจสอบแต่ละ field จาก form-data และเพิ่มลงใน map ถ้ามีค่า
	if title := c.PostForm("title"); title != "" {
		updateData["title"] = title
	}
	if description := c.PostForm("description"); description != "" {
		updateData["description"] = description
	}
	if priceStr := c.PostForm("price"); priceStr != "" {
		if price, err := strconv.ParseFloat(priceStr, 64); err == nil {
			updateData["price"] = price
		}
	}
	if categoryIDStr := c.PostForm("category_id"); categoryIDStr != "" {
		if categoryID, err := strconv.Atoi(categoryIDStr); err == nil {
			updateData["category_id"] = uint(categoryID)
		}
	}

	// 5. จัดการไฟล์รูปภาพ (ถ้ามีการส่งไฟล์ใหม่มา)
	if file, err := c.FormFile("image"); err == nil {
		extension := filepath.Ext(file.Filename)
		newFileName := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
		filePath := "uploads/" + newFileName
		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save uploaded file"})
			return
		}
		updateData["image_game"] = filePath
	}

	// 6. ตรวจสอบว่ามีข้อมูลให้อัปเดตหรือไม่
	if len(updateData) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No data provided for update"})
		return
	}

	// 7. อัปเดตข้อมูลลงฐานข้อมูล
	if err := db.Model(&game).Updates(updateData).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update game", "details": err.Error()})
		return
	}

	// ดึงข้อมูลล่าสุดหลังจากอัปเดตเพื่อส่งกลับ
	db.Preload("Category").First(&game, id)
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Game updated successfully",
		"data":    game,
	})
}

// DeleteGameHandler จัดการการลบเกม (Soft Delete)
func DeleteGameHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง ID จาก URL parameter
	id := c.Param("id")

	var game model.Game
	// 2. ค้นหาเกมที่จะลบก่อน เพื่อจะเอา path ของรูปภาพมาใช้ลบไฟล์ทิ้ง
	if err := db.First(&game, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Game not found"})
		return
	}
	
	imagePath := game.ImageGame

	// 3. สั่งลบเกมด้วย GORM
	//    เพราะเราเพิ่ม gorm.DeletedAt ใน model แล้ว GORM จะทำ Soft Delete ให้เอง
	//    (GORM จะรันคำสั่ง UPDATE games SET deleted_at = 'เวลาปัจจุบัน' WHERE id = ...)
	if err := db.Delete(&model.Game{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete game"})
		return
	}

	// 4. (แนะนำ) ลบไฟล์รูปภาพออกจากเซิร์ฟเวอร์ด้วย เพื่อไม่ให้เปลืองพื้นที่
	if imagePath != "" {
		os.Remove(imagePath)
	}

	// 5. ส่ง Status 204 No Content กลับไป ซึ่งเป็นมาตรฐานสำหรับ DELETE request ที่สำเร็จ
	c.Status(http.StatusNoContent)
}


// GetTopSellingGamesHandler handles fetching the top 5 best-selling games.
func GetTopSellingGamesHandler(c *gin.Context, db *gorm.DB) {

	// สร้าง struct ชั่วคราวเพื่อรับผลลัพธ์จากการจัดอันดับ
	type GameRank struct {
		GameID        uint
		PurchaseCount int
	}

	var ranks []GameRank

	// --- 1. Query เพื่อหา 5 game_id ที่ถูกซื้อมากที่สุด ---
	err := db.Model(&model.OrderDetail{}).
		Select("game_id, count(game_id) as purchase_count").
		Group("game_id").
		Order("purchase_count desc").
		Limit(5).
		Find(&ranks).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not query game ranks"})
		return
	}

	// กรณีที่ยังไม่มีการซื้อขายเลย
	if len(ranks) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"status":  "success",
			"message": "No sales data available yet",
			"data":    []model.Game{}, // ส่ง array ว่างกลับไป
		})
		return
	}

	// --- 2. ดึง Game ID ทั้งหมดออกมาใส่ใน slice ---
	var gameIDs []uint
	for _, rank := range ranks {
		gameIDs = append(gameIDs, rank.GameID)
	}

	// --- 3. ดึงข้อมูลเกมฉบับเต็มของ ID ทั้ง 5 อันดับ ---
	var topGames []model.Game
	if err := db.Preload("Category").Where("game_id IN ?", gameIDs).Find(&topGames).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch top game details"})
		return
	}

	// --- 4. จัดเรียงข้อมูลเกมให้ตรงกับอันดับที่ได้มา ---
	// (เพราะ .Where("... IN ...") ไม่ได้รับประกันลำดับ)
	gameMap := make(map[uint]model.Game)
	for _, game := range topGames {
		gameMap[game.GameID] = game
	}

	var rankedGames []model.Game
	for _, rank := range ranks {
		if game, ok := gameMap[rank.GameID]; ok {
			rankedGames = append(rankedGames, game)
		}
	}

	// --- 5. ส่งผลลัพธ์กลับไป ---
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Top 5 selling games fetched successfully",
		"data":    rankedGames,
	})
}