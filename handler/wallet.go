package handlers

import (
	"api-game/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ✅ ฟังก์ชันเติมเงินเข้ากระเป๋า
func AddWalletHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		UserID uint    `json:"user_id"`
		Amount float64 `json:"amount"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "จำนวนเงินต้องมากกว่า 0"})
		return
	}

	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
		return
	}

	// บวกเงินเพิ่มใน wallet
	user.Wallet += req.Amount
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตยอดเงินได้"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "เติมเงินสำเร็จ",
		"wallet":  user.Wallet,
	})
}
