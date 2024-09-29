package model

import (
	"time"
)

type Cart struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement:true" json:"id"`
	UserID       uint64         `gorm:"index;not null" json:"user_id"`
	CreatedAt    time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	CartProducts []CartProduct `gorm:"foreignKey:CartID" json:"cart_products"`
}

type CartProduct struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement:true" json:"id"`
	CartID    uint64     `gorm:"index;not null" json:"cart_id"`
	ProductID uint64     `gorm:"index;not null" json:"product_id"`
	Quantity  int       `gorm:"not null" json:"quantity"`
	AddedAt   time.Time `gorm:"autoCreateTime" json:"added_at"`

	Cart    Cart        `gorm:"foreignKey:CartID" json:"cart"`
	Product ProductSize `gorm:"foreignKey:ProductID" json:"product"`
}
