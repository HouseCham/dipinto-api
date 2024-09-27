package model

import (
	"time"
)

type Cart struct {
	ID           uint64         `gorm:"primaryKey;autoIncrement:true"`
	UserID       uint64         `gorm:"index;not null"`
	CreatedAt    time.Time     `gorm:"autoCreateTime"`
	UpdatedAt    time.Time     `gorm:"autoUpdateTime"`
	CartProducts []CartProduct `gorm:"foreignKey:CartID"`
}

type CartProduct struct {
	ID        uint64     `gorm:"primaryKey;autoIncrement:true"`
	CartID    uint64     `gorm:"index;not null"`
	ProductID uint64     `gorm:"index;not null"`
	Quantity  int       `gorm:"not null"`
	AddedAt   time.Time `gorm:"autoCreateTime"`

	Cart    Cart        `gorm:"foreignKey:CartID"`
	Product ProductSize `gorm:"foreignKey:ProductID"`
}
