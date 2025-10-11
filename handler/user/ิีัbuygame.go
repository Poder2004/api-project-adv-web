package handlers // หรือ package ที่คุณใช้

import (
	"api-game/model"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CheckoutHandler handles the entire game purchase process.
func CheckoutHandler(c *gin.Context, db *gorm.DB) {
	// --- 1. รับข้อมูลจาก Frontend ---
	var request struct {
		GameIDs  []uint `json:"game_ids" binding:"required"`
		CouponID *uint  `json:"coupon_id"` // ใช้ pointer เพื่อให้รับค่า null ได้
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request data", "details": err.Error()})
		return
	}

	if len(request.GameIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Cart is empty"})
		return
	}

	// --- 2. ดึงข้อมูล User ---
	userIDAny, _ := c.Get("user_id")
	userID := uint(userIDAny.(float64))

	var currentUser model.User
	if err := db.First(&currentUser, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}

	// --- 3. เริ่ม Transaction ---
	err := db.Transaction(func(tx *gorm.DB) error {
		// --- 3.1 ดึงข้อมูลเกมและคำนวณราคารวม (Subtotal) ---
		var games []model.Game
		if err := tx.Where("game_id IN ?", request.GameIDs).Find(&games).Error; err != nil {
			return err
		}
		if len(games) != len(request.GameIDs) {
			return errors.New("one or more games not found")
		}

		var sumTotal float64
		for _, game := range games {
			sumTotal += game.Price
		}

		var discount float64
		var finalTotal float64

		// --- 3.2 ตรวจสอบและคำนวณส่วนลดจากคูปอง (ถ้ามี) ---
		if request.CouponID != nil {
			var userCoupon model.UserCoupon
			// ตรวจสอบว่าผู้ใช้มีคูปองใบนี้จริงและยังไม่ได้ใช้
			err := tx.Where("user_id = ? AND did = ? AND is_used = ?", userID, *request.CouponID, false).First(&userCoupon).Error
			if err != nil {
				return errors.New("invalid or already used coupon")
			}
			
			var couponDetails model.DiscountCode
			if err := tx.First(&couponDetails, *request.CouponID).Error; err != nil {
				return errors.New("coupon details not found")
			}
			
			// ตรวจสอบยอดซื้อขั้นต่ำ
			if sumTotal < couponDetails.MinValue {
				return errors.New("subtotal does not meet coupon's minimum value")
			}

			// คำนวณส่วนลด
			if couponDetails.DiscountType == "fixed" {
				discount = couponDetails.DiscountValue
			} else if couponDetails.DiscountType == "percent" {
				discount = (sumTotal * couponDetails.DiscountValue) / 100
			}

			// อัปเดตสถานะคูปองเป็น "ใช้แล้ว"
			if err := tx.Model(&userCoupon).Update("is_used", true).Error; err != nil {
				return err
			}
		}

		finalTotal = sumTotal - discount

		// --- 3.3 ตรวจสอบ Wallet และหักเงิน ---
		if currentUser.Wallet < finalTotal {
			return errors.New("insufficient wallet balance")
		}
		
		newWalletBalance := currentUser.Wallet - finalTotal
		if err := tx.Model(&currentUser).Update("wallet", newWalletBalance).Error; err != nil {
			return err
		}

		// --- 3.4 สร้าง Order Record ---
		newOrder := model.Order{
			UserID:     userID,
			DID:        request.CouponID,
			Discount:   discount,
			SumTotal:   sumTotal,
			FinalTotal: finalTotal,
			OrderDate:  time.Now(),
		}
		if err := tx.Create(&newOrder).Error; err != nil {
			return err
		}

		// --- 3.5 สร้าง Order Detail และเพิ่มเกมเข้า Library ---
		for _, game := range games {
			// สร้าง Order Detail
			orderDetail := model.OrderDetail{
				OrdersID: newOrder.OrdersID,
				GameID:   game.GameID,
			}
			if err := tx.Create(&orderDetail).Error; err != nil {
				return err
			}

			// เพิ่มเกมเข้า User Library
			userLibraryEntry := model.UserLibrary{
				UserID: userID,
				GameID: game.GameID,
			}
			// ใช้ Clauses(clause.OnConflict{DoNothing: true}) เพื่อไม่ให้เกิด error ถ้ามีเกมนั้นอยู่แล้ว
			if err := tx.Create(&userLibraryEntry).Error; err != nil {
                // ถ้าไม่ต้องการให้ error ถ้ามีเกมซ้ำ ให้ comment บรรทัดล่างแล้วไปใช้ OnConflict แทน
				return errors.New("game already in library")
			}
		}

		return nil // Commit Transaction
	}) // --- สิ้นสุด Transaction ---

	// --- 4. ส่งผลลัพธ์กลับ ---
	if err != nil {
		// แยกประเภท Error เพื่อให้ Frontend แสดงผลได้ถูกต้อง
		switch err.Error() {
		case "insufficient wallet balance":
			c.JSON(http.StatusPaymentRequired, gin.H{"error": err.Error()})
		case "one or more games not found", "invalid or already used coupon", "subtotal does not meet coupon's minimum value", "game already in library":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "An unexpected error occurred", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Checkout successful!",
	})
}