package handlersadmin

import (
	"api-game/model"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateCouponHandler handles the creation of a new discount code (ฉบับอัปเดต)
func CreateCouponHandler(c *gin.Context, db *gorm.DB) {
	// 1. อัปเดต struct ที่รับข้อมูลจาก JSON body ให้ตรงกับ UI และ Model ใหม่
	var input struct {
		NameCode      string  `json:"name_code" binding:"required"`
		Description   string  `json:"description"`
		DiscountValue float64 `json:"discount_value" binding:"required"`
		DiscountType  string  `json:"discount_type" binding:"required"` // 'fixed' (บาท) or 'percent' (%)
		MinValue      float64 `json:"min_value"`                      // ยอดซื้อขั้นต่ำ (optional)
		LimitUsage    int     `json:"limit_usage" binding:"required"`
	}

	// Bind ข้อมูล JSON และตรวจสอบความถูกต้องเบื้องต้น
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	// 2. ตรวจสอบค่าของ DiscountType ที่ส่งมา
	if input.DiscountType != "fixed" && input.DiscountType != "percent" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid discount_type. Must be 'fixed' or 'percent'"})
		return
	}

	// 3. สร้าง instance ของ DiscountCode model ด้วยข้อมูลใหม่ทั้งหมด
	coupon := model.DiscountCode{
		NameCode:      input.NameCode,
		Description:   input.Description,
		DiscountValue: input.DiscountValue,
		DiscountType:  input.DiscountType,
		MinValue:      input.MinValue,
		LimitUsage:    input.LimitUsage,
		UsedCount:     0, // ค่าเริ่มต้นตอนสร้าง
	}

	// 4. บันทึกข้อมูลลงฐานข้อมูล
	if err := db.Create(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create coupon", "details": err.Error()})
		return
	}

	// 5. ส่งผลลัพธ์กลับไป
	c.JSON(http.StatusCreated, gin.H{
		"status":  "success",
		"message": "Discount coupon created successfully",
		"data":    coupon,
	})
}

// GetAllCouponsHandler จัดการการดึงข้อมูลคูปองทั้งหมด
func GetAllCouponsHandler(c *gin.Context, db *gorm.DB) {
	var coupons []model.DiscountCode

	// ดึงข้อมูลทั้งหมด โดยเรียงจากใหม่ไปเก่า
	if err := db.Order("did desc").Find(&coupons).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not fetch coupons"})
		return
	}

	if coupons == nil {
		coupons = []model.DiscountCode{}
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Coupons fetched successfully",
		"data":    coupons,
	})
}


// --- [เพิ่ม] Struct สำหรับรับข้อมูลตอน Update ---
// เราไม่รับ NameCode เพราะปกติแล้วเราไม่ควรให้แก้ไขรหัสคูปอง
type UpdateCouponInput struct {
	NameCode      string  `json:"name_code" binding:"required"`
	Description   string  `json:"description" binding:"required"`
	DiscountValue float64 `json:"discount_value" binding:"required"`
	DiscountType  string  `json:"discount_type" binding:"required"`
	MinValue      float64 `json:"min_value"`
	LimitUsage    int     `json:"limit_usage" binding:"required"`
}

// --- [เพิ่ม] ฟังก์ชัน UpdateCouponHandler ---
func UpdateCouponHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id")
	var input UpdateCouponInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	var coupon model.DiscountCode
	if err := db.First(&coupon, id).Error; err != nil {
        // ... (error handling เหมือนเดิม) ...
		return
	}

	// --- 👇 อัปเดตค่าต่างๆ รวมถึง NameCode ---
    coupon.NameCode = input.NameCode // 👈 เพิ่มบรรทัดนี้
	coupon.Description = input.Description
	coupon.DiscountValue = input.DiscountValue
	coupon.DiscountType = input.DiscountType
	coupon.MinValue = input.MinValue
	coupon.LimitUsage = input.LimitUsage

	// บันทึกการเปลี่ยนแปลงลงฐานข้อมูล
	if err := db.Save(&coupon).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update coupon"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success", "data": coupon})
}


// --- [เพิ่ม] ฟังก์ชัน DeleteCouponHandler (Soft Delete) ---
func DeleteCouponHandler(c *gin.Context, db *gorm.DB) {
	id := c.Param("id") // ดึง ID ของคูปองที่จะลบจาก URL

	// GORM จะทำ Soft Delete อัตโนมัติ เพราะใน model มี gorm.DeletedAt
	// มันจะรันคำสั่ง: UPDATE "discount_code" SET "deleted_at"='...' WHERE "did" = ?
	if err := db.Delete(&model.DiscountCode{}, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Coupon not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete coupon"})
		return
	}

	// ส่ง Status 204 No Content กลับไปเมื่อสำเร็จ (เป็นมาตรฐานของ REST API)
	c.Status(http.StatusNoContent)
}