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
	ID         uint64    `gorm:"primaryKey;autoIncrement:true" json:"id"`
	WishlistID uint64    `gorm:"index;not null" json:"wishlist_id"`
	ProductID  uint64    `gorm:"index;not null" json:"product_id"`
	AddedAt    time.Time `gorm:"autoCreateTime" json:"added_at"`

	Wishlist Wishlist `gorm:"foreignKey:WishlistID" json:"wishlist"`
	Product  Product  `gorm:"foreignKey:ProductID" json:"product"`
}