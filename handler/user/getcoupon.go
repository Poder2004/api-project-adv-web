// handlers/coupon_handler.go

package handlers

import (
	"api-game/model"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// ClaimCouponHandler handles a user's request to claim a discount coupon
func ClaimCouponHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง user_id จาก Middleware และแปลงประเภทข้อมูลอย่างปลอดภัย
	userIDAny, exists := c.Get("user_id") // 👈 เปลี่ยนชื่อตัวแปรเพื่อความชัดเจน
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// แปลงประเภทจาก interface{} เป็น float64 (ซึ่งเป็นค่า default จาก JWT)
	userIDFloat, ok := userIDAny.(float64)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Invalid user ID type in context"})
		return
	}
	// แปลงจาก float64 เป็น uint เพื่อใช้งานจริง
	finalUserID := uint(userIDFloat) // 👈 ได้ user_id ที่เป็น uint แล้ว

	// 2. ดึง did (coupon id) จาก URL parameter
	didStr := c.Param("did")
	did, err := strconv.ParseUint(didStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid coupon ID format"})
		return
	}

	// 3. ใช้ Transaction เพื่อให้แน่ใจว่าทุกอย่างสำเร็จหรือล้มเหลวพร้อมกัน
	err = db.Transaction(func(tx *gorm.DB) error {
		var coupon model.DiscountCode

		// 4. ตรวจสอบว่าคูปองมีอยู่จริง และยังเหลือให้กดรับหรือไม่
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&coupon, did).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("coupon not found")
			}
			return err
		}

		if coupon.UsedCount >= coupon.LimitUsage {
			return errors.New("coupon has reached its usage limit")
		}

		// 5. สร้างข้อมูลการรับคูปองในตาราง user_coupons
		userCoupon := model.UserCoupon{
			UserID: finalUserID, // 👈 ใช้ user_id ที่แปลงค่าแล้ว
			DID:    uint(did),
			IsUsed: false,
		}

		if err := tx.Create(&userCoupon).Error; err != nil {
			return errors.New("you have already claimed this coupon")
		}

		// 6. อัปเดต used_count ในตาราง discount_code (+1)
		coupon.UsedCount++
		if err := tx.Save(&coupon).Error; err != nil {
			return err
		}

		return nil
	})

	// 7. จัดการผลลัพธ์จาก Transaction
	if err != nil {
		switch err.Error() {
		case "coupon not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "coupon has reached its usage limit", "you have already claimed this coupon":
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to claim coupon", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Coupon claimed successfully!",
	})

	
}

// GetMyCouponsHandler handles fetching all claimed coupon IDs for the current user
func GetMyCouponsHandler(c *gin.Context, db *gorm.DB) {
	// 1. ดึง user_id จาก Middleware
	userIDAny, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	userID := uint(userIDAny.(float64)) // แปลงประเภทให้ถูกต้อง

	var claimedCouponIDs []uint

	// 2. คิวรีจากตาราง user_coupons โดยเลือกเฉพาะคอลัมน์ 'did'
	// สำหรับ user_id ที่ตรงกับคนที่ล็อกอินอยู่
	err := db.Model(&model.UserCoupon{}).
		Where("user_id = ?", userID).
		Pluck("did", &claimedCouponIDs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch user's coupons"})
		return
	}
	
	// ถ้าไม่มีคูปองเลย ให้ส่งกลับเป็น array ว่างๆ แทนที่จะเป็น null
	if claimedCouponIDs == nil {
		claimedCouponIDs = []uint{}
	}

	// 3. ส่งข้อมูลกลับไปเป็น JSON
	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User's claimed coupons fetched successfully",
		"data":    claimedCouponIDs,
	})

}

// GetMyAvailableCouponsHandler fetches full coupon details that the user has claimed but not yet used.
func GetMyAvailableCouponsHandler(c *gin.Context, db *gorm.DB) {
	userIDAny, _ := c.Get("user_id")
	userID := uint(userIDAny.(float64))

	var availableCoupons []model.DiscountCode

	// ใช้ JOIN เพื่อดึงข้อมูลจากตาราง discount_code
	// โดยอ้างอิงจากตาราง user_coupons ที่ is_used = false
	err := db.Joins("JOIN user_coupons ON user_coupons.did = discount_code.did").
		Where("user_coupons.user_id = ? AND user_coupons.is_used = ?", userID, false).
		Find(&availableCoupons).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch available coupons"})
		return
	}

	if availableCoupons == nil {
		availableCoupons = []model.DiscountCode{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "User's available coupons fetched successfully",
		"data":    availableCoupons,
	})
}

	

