package model

import (
	"time"

	"gorm.io/gorm"
)

// User แทนข้อมูลในตาราง users
type User struct {
	UserID       uint    `gorm:"primaryKey;column:user_id;autoIncrement" json:"user_id"`
	Username     string  `gorm:"column:username;type:varchar(50);not null;unique" json:"username"`
	Email        string  `gorm:"column:email;type:varchar(255);not null;unique" json:"email"`
	Password     string  `gorm:"column:password;type:varchar(255);not null" json:"-"` // ซ่อน Password จาก JSON
	Role         string  `gorm:"column:role;type:enum('member','admin');default:'member';not null" json:"role"`
	ImageProfile string  `gorm:"column:image_profile;type:varchar(255)" json:"image_profile"`
	Wallet       float64 `gorm:"column:wallet;type:decimal(10,2);default:0" json:"wallet"`
}

// Category แทนข้อมูลในตาราง categories
type Category struct {
	CategoryID   uint   `gorm:"primaryKey;column:category_id" json:"category_id"`
	CategoryName string `gorm:"column:category_name;type:varchar(255);not null" json:"category_name"`

	// Relationship: One-to-Many
	Games []Game `gorm:"foreignKey:CategoryID" json:"games,omitempty"`
}

// Game แทนข้อมูลในตาราง games
type Game struct {
	GameID      uint           `gorm:"primaryKey;column:game_id" json:"game_id"`
	Title       string         `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description string         `gorm:"column:description;type:text" json:"description"`
	Price       float64        `gorm:"column:price;type:decimal(10,2);not null" json:"price"`
	ImageGame   string         `gorm:"column:image_game;type:varchar(255)" json:"image_game"`
	ReleaseDate time.Time      `gorm:"column:release_date;type:date" json:"release_date"`
	CategoryID  uint           `gorm:"column:category_id" json:"category_id"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Relationships
	Category Category `gorm:"foreignKey:CategoryID;references:CategoryID"`
}

// DiscountCode แทนข้อมูลในตาราง discount_code
type DiscountCode struct {
	DID           uint    `gorm:"primaryKey;column:did" json:"did"`
	NameCode      string  `gorm:"column:name_code;type:varchar(50);unique;not null" json:"name_code"`
	Description   string  `gorm:"column:description;type:text" json:"description"`
	DiscountValue float64 `gorm:"column:discount_value;type:decimal(10,2);not null" json:"discount_value"`
	DiscountType  string  `gorm:"column:discount_type;type:varchar(10);not null;default:'fixed'" json:"discount_type"`
	MinValue      float64 `gorm:"column:min_value;type:decimal(10,2);not null;default:0.00" json:"min_value"`
	LimitUsage    int     `gorm:"column:limit_usage" json:"limit_usage"`
	UsedCount     int     `gorm:"column:used_count;default:0" json:"used_count"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"` 
}

func (DiscountCode) TableName() string {
	return "discount_code"
}

// Order แทนข้อมูลในตาราง orders (ฉบับแก้ไขสมบูรณ์)
type Order struct {
	OrdersID   uint      `gorm:"primaryKey;column:orders_id" json:"orders_id"`
	UserID     uint      `gorm:"column:user_id" json:"user_id"`
	DID        *uint     `gorm:"column:did" json:"did"`
	Discount   float64   `gorm:"column:discount;type:decimal(10,2)" json:"discount"`
	SumTotal   float64   `gorm:"column:sum_total;type:decimal(10,2);not null" json:"sum_total"`
	FinalTotal float64   `gorm:"column:final_total;type:decimal(10,2);not null" json:"final_total"`
	OrderDate  time.Time `gorm:"column:order_date;type:datetime" json:"order_date"`

	// --- 👇 [แก้ไขความสัมพันธ์ตรงนี้] ---
	User         User          `gorm:"foreignKey:UserID;references:UserID"` // บอกว่า: ให้ใช้ UserID ของตารางนี้ ไปอ้างอิงกับ UserID ของตาราง User
	DiscountCode DiscountCode  `gorm:"foreignKey:DID;references:DID"`
	OrderDetails []OrderDetail `gorm:"foreignKey:OrdersID;references:OrdersID"`
}

// OrderDetail แทนข้อมูลในตาราง orders_detail (ฉบับแก้ไขสมบูรณ์)
type OrderDetail struct {
	OdID     uint `gorm:"primaryKey;column:od_id" json:"od_id"`
	OrdersID uint `gorm:"column:orders_id" json:"orders_id"`
	GameID   uint `gorm:"column:game_id" json:"game_id"`

	// --- 👇 [แก้ไขความสัมพันธ์ตรงนี้] ---
	Game Game `gorm:"foreignKey:GameID;references:GameID"` // บอกว่า: ให้ใช้ GameID ของตารางนี้ ไปอ้างอิงกับ GameID ของตาราง Game
}

func (OrderDetail) TableName() string {
	return "orders_detail"
}

// WalletHistory แทนข้อมูลในตาราง wallet_history
type WalletHistory struct {
	HistoryID       uint      `gorm:"primaryKey;column:history_id" json:"history_id"`
	UserID          uint      `gorm:"column:user_id;not null" json:"user_id"`
	Amount          float64   `gorm:"column:amount;type:decimal(10,2);not null" json:"amount"`
	TransactionDate time.Time `gorm:"column:transaction_date;type:datetime;not null" json:"transaction_date"`
}

func (WalletHistory) TableName() string {
	return "wallet_history"
}

// UserLibrary แทนข้อมูลในตาราง user_library (Join Table)
type UserLibrary struct {
	UserID uint `gorm:"primaryKey;column:user_id"`
	GameID uint `gorm:"primaryKey;column:game_id"`
}

func (UserLibrary) TableName() string {
	return "user_library"
}

// UserCoupon แทนข้อมูลในตาราง user_coupons
type UserCoupon struct {
	UserCouponID uint `gorm:"primaryKey;column:user_coupon_id;autoIncrement" json:"user_coupon_id"`
	UserID       uint `gorm:"column:user_id;not null" json:"user_id"`
	DID          uint `gorm:"column:did;not null" json:"did"`
	IsUsed       bool `gorm:"column:is_used;not null;default:false" json:"is_used"`
}

func (UserCoupon) TableName() string {
	return "user_coupons"
}
