package model

import "time"

// Category แทนข้อมูลในตาราง categories
type Category struct {
	CategoryID   uint   `gorm:"primaryKey" json:"category_id"`
	CategoryName string `gorm:"type:varchar(255);not null" json:"category_name"`

	// Relationship: One-to-Many
	Games []Game `gorm:"foreignKey:CategoryID" json:"games,omitempty"`
}

// Game แทนข้อมูลในตาราง games
type Game struct {
	GameID      uint      `gorm:"primaryKey" json:"game_id"`
	Title       string    `gorm:"type:varchar(255);not null" json:"title"`
	Description string    `gorm:"type:text" json:"description"`
	Price       float64   `gorm:"type:decimal(10,2);not null" json:"price"`
	ImageGame   string    `gorm:"type:varchar(255)" json:"image_game"`
	ReleaseDate time.Time `gorm:"type:date" json:"release_date"`
	CategoryID  uint      `json:"category_id"`

	// Relationships
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

// DiscountCode แทนข้อมูลในตาราง discount_code
type DiscountCode struct {
	DID         uint   `gorm:"primaryKey" json:"did"`
	NameCode    string `gorm:"type:varchar(50);unique;not null" json:"name_code"`
	Description string `gorm:"type:text" json:"description"`
	LimitUsage  int    `json:"limit_usage"`
	UsedCount   int    `gorm:"default:0" json:"used_count"`
}

// Order แทนข้อมูลในตาราง orders
type Order struct {
	OrdersID   uint      `gorm:"primaryKey" json:"orders_id"`
	UserID     uint      `json:"user_id"`
	DID        *uint     `json:"did"` // ใช้ pointer (*uint) เพื่อให้รับค่า NULL ได้
	Discount   float64   `gorm:"type:decimal(10,2)" json:"discount"`
	SumTotal   float64   `gorm:"type:decimal(10,2);not null" json:"sum_total"`
	FinalTotal float64   `gorm:"type:decimal(10,2);not null" json:"final_total"`
	OrderDate  time.Time `gorm:"type:datetime" json:"order_date"`

	// Relationships
	User         User          `gorm:"foreignKey:UserID" json:"user,omitempty"`
	DiscountCode DiscountCode  `gorm:"foreignKey:DID" json:"discount_code,omitempty"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrdersID" json:"order_details,omitempty"`
}

// OrderDetail แทนข้อมูลในตาราง orders_detail
type OrderDetail struct {
	OdID     uint `gorm:"primaryKey" json:"od_id"`
	OrdersID uint `json:"orders_id"`
	GameID   uint `json:"game_id"`

	// Relationship
	Game Game `gorm:"foreignKey:GameID" json:"game,omitempty"`
}