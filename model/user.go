package model

type User struct {
	UserID       uint    `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
	Username     string  `json:"username" gorm:"column:username;type:varchar(50);not null;unique"`
	Email        string  `json:"email" gorm:"column:email;type:varchar(255);not null;unique"`
	Password     string  `json:"password" gorm:"column:password;type:varchar(255);not null"`
	Role         string  `json:"role" gorm:"column:role;type:enum('member','admin');default:'member';not null"`
	ImageProfile string  `json:"image_profile" gorm:"column:image_profile;type:varchar(255)"`
	Wallet       float64 `json:"wallet" gorm:"column:wallet;type:decimal(10,2);default:0"`

	// Relationships
	// Orders   []Order `gorm:"foreignKey:UserID" json:"orders,omitempty"`
	// Library  []Game  `gorm:"many2many:user_library;" json:"library,omitempty"` // ความสัมพันธ์ many-to-many
}

type UserLibrary struct {
	UserID uint `gorm:"primaryKey"`
	GameID uint `gorm:"primaryKey"`
}

// ✅ เพิ่มฟังก์ชันนี้เข้าไปใต้ struct UserLibrary
// เพื่อบอก GORM ให้ใช้ชื่อตาราง "user_library" ที่ถูกต้อง
func (UserLibrary) TableName() string {
	return "user_library"
}

type UserCoupon struct {
	UserCouponID uint `gorm:"primaryKey;autoIncrement" json:"user_coupon_id"`
	UserID       uint `gorm:"not null" json:"user_id"`
	DID          uint `gorm:"column:did;not null" json:"did"`
	IsUsed       bool `gorm:"not null;default:false" json:"is_used"`
	// แบบใหม่ที่ถูกต้อง

}

func (UserCoupon) TableName() string {
	return "user_coupons" // บอก GORM ให้ใช้ชื่อตารางนี้
}