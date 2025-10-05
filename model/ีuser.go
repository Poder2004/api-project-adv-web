package models

type User struct {
	UserID       uint    `json:"user_id" gorm:"column:user_id;primaryKey;autoIncrement"`
	Username     string  `json:"username" gorm:"column:username;type:varchar(50);not null;unique"`
	Email        string  `json:"email" gorm:"column:email;type:varchar(255);not null;unique"`
	Password     string  `json:"password" gorm:"column:password;type:varchar(255);not null"`
	Role         string  `json:"role" gorm:"column:role;type:enum('member','admin');default:'member';not null"`
	ImageProfile string  `json:"image_profile" gorm:"column:image_profile;type:varchar(255)"`
	Wallet       float64 `json:"wallet" gorm:"column:wallet;type:decimal(10,2);default:0"`
}


