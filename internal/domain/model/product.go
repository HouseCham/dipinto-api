package model

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement"`
	CategoryID  uint64          `gorm:"not null" validate:"required,numeric"`
	Slug        string          `gorm:"unique;not null" validate:"required,slug"`
	Name        string          `gorm:"not null" validate:"required"`
	Description string          `gorm:"type:text" validate:"required"`
	Images      json.RawMessage `gorm:"type:jsonb"`
	CreatedAt   time.Time       `gorm:"default:now()"`
	UpdatedAt   time.Time       `gorm:"default:now()"`
	DeletedAt   *time.Time
}

type ProductSize struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	ProductID   uint64    `gorm:"not null" validate:"required,numeric"`
	IsAvailable bool      `gorm:"not null" validate:"boolean"`
	SizeSlug    string    `gorm:"not null" validate:"required,slug"`
	Size        string    `gorm:"not null" validate:"required"`
	Price       float64   `gorm:"type:numeric(10,2);not null" validate:"required,numeric"`
	Discount    float64   `gorm:"type:numeric" validate:"numeric"`
	CreatedAt   time.Time `gorm:"default:now()"`
	UpdatedAt   time.Time `gorm:"default:now()"`
	DeletedAt   *time.Time
}

type CatalogueProduct struct {
	ID       uint64 `json:"id"`
	Slug     string `json:"slug"`
	Name     string `json:"name"`
	ImageUrl string `json:"image_url"`
	Price    string `json:"price"`
	Discount string `json:"discount"`
}

func (ProductSize) TableName() string {
	return "product_sizes"
}
func (Product) TableName() string {
	return "products"
}
