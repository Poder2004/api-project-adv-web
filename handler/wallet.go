package handlers

import (
	"api-game/model"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ✅ ฟังก์ชันเติมเงินเข้ากระเป๋า (รองรับ transaction_date)
func AddWalletHandler(c *gin.Context, db *gorm.DB) {
	var req struct {
		UserID          uint    `json:"user_id"`
		Amount          float64 `json:"amount"`
		TransactionDate *string `json:"transaction_date"` // optional: 'YYYY-MM-DD' หรือ RFC3339
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}
	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "จำนวนเงินต้องมากกว่า 0"})
		return
	}

	// หา user
	var user model.User
	if err := db.First(&user, req.UserID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "ไม่พบผู้ใช้"})
		return
	}

	// อัปเดตยอด wallet
	user.Wallet += req.Amount
	if err := db.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "ไม่สามารถอัปเดตยอดเงินได้"})
		return
	}

	// ✅ แปลงวันที่จาก client (ถ้ามี) แล้ว fallback เป็นวันนี้เมื่อไม่มี/parse ไม่ได้
	loc, _ := time.LoadLocation("Asia/Bangkok")
	var txTime time.Time
	if req.TransactionDate != nil && *req.TransactionDate != "" {
		// 1) RFC3339 เช่น "2025-10-10T00:00:00+07:00"
		if t, err := time.Parse(time.RFC3339, *req.TransactionDate); err == nil {
			txTime = t.In(loc)
		} else if t2, err2 := time.ParseInLocation("2006-01-02", *req.TransactionDate, loc); err2 == nil {
			// 2) 'YYYY-MM-DD'
			txTime = t2
		} else {
			// 3) parse ไม่ได้ → วันนี้ (กันพัง)
			txTime = time.Now().In(loc)
		}
	} else {
		// ไม่ส่งมา → วันนี้
		txTime = time.Now().In(loc)
	}

	// บันทึกประวัติด้วยวันที่ที่ได้
	history := model.WalletHistory{
		UserID:          req.UserID,
		Amount:          req.Amount,
		TransactionDate: txTime,
	}
	if err := db.Create(&history).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "บันทึกประวัติเติมเงินไม่สำเร็จ"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "เติมเงินสำเร็จ",
		"wallet":  user.Wallet,
	})
}

// ✅ ประวัติการเติมเงิน (คงเดิม)
func GetWalletHistoryHandler(c *gin.Context, db *gorm.DB) {
	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"status": "error", "message": "missing user_id"})
		return
	}

	var histories []model.WalletHistory
	if err := db.Where("user_id = ?", userIDStr).
		Order("transaction_date DESC, history_id DESC").
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
