package handlers

import (
	"api-game/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"time"
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
	// ---- เก็บประวัติการเติมเงิน (เพิ่มบล็อกนี้) ----
	history := model.WalletHistory{
		UserID:          req.UserID,
		Amount:          req.Amount,
		TransactionDate: time.Now(),
	}
	if err := db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกประวัติเติมเงินไม่สำเร็จ"})
		return
	}
	// -----------------------------------------------
	c.JSON(http.StatusOK, gin.H{
		"message": "เติมเงินสำเร็จ",
		"wallet":  user.Wallet,
	})
}

// .ประวัติการเติมเงิน

func GetWalletHistoryHandler(c *gin.Context, db *gorm.DB) {
    userIDStr := c.Query("user_id")
    if userIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing user_id"})
        return
    }

    var histories []model.WalletHistory
    if err := db.Where("user_id = ?", userIDStr).
        Order("transaction_date DESC").
        Find(&histories).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"status": "error", "message": "query history failed"})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "status":  "success",
        "message": "ok",
        "data":    histories,
    })
}
