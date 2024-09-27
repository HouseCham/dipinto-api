package model

import (
	"time"
)

type Wishlist struct {
	ID               uint64            `gorm:"primaryKey;autoIncrement:true"`
	UserID           uint64            `gorm:"index;not null"` // Assuming the users table is already defined elsewhere
	CreatedAt        time.Time         `gorm:"autoCreateTime"`
	UpdatedAt        time.Time         `gorm:"autoUpdateTime"`
	WishlistProducts []WishlistProduct `gorm:"foreignKey:WishlistID"` // Association with WishlistProducts
}

type WishlistProduct struct {
	ID         uint64    `gorm:"primaryKey;autoIncrement:true"`
	WishlistID uint64    `gorm:"index;not null"`
	ProductID  uint64    `gorm:"index;not null"`
	AddedAt    time.Time `gorm:"autoCreateTime"`

	Wishlist Wishlist `gorm:"foreignKey:WishlistID"`
	Product  Product  `gorm:"foreignKey:ProductID"`
}
