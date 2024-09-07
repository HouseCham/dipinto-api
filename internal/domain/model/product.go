package model

import (
	"encoding/json"
	"time"
)

type Product struct {
	ID          uint64          `gorm:"primaryKey;autoIncrement" json:"id"`
	CategoryID  uint64          `gorm:"not null" validate:"required,numeric" json:"category_id"`
	Slug        string          `gorm:"unique;not null" validate:"required,slug" json:"slug"`
	Name        string          `gorm:"not null" validate:"required" json:"name"`
	Description string          `gorm:"type:text" validate:"required" json:"description"`
	Images      json.RawMessage `gorm:"type:jsonb" json:"images"`
	CreatedAt   time.Time       `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time       `gorm:"default:now()" json:"updated_at"`
	DeletedAt   *time.Time
}

type ProductSize struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ProductID   uint64    `gorm:"not null" validate:"numeric" json:"product_id"`
	IsAvailable bool      `gorm:"not null" validate:"boolean" json:"is_available"`
	SizeSlug    string    `gorm:"not null" validate:"required,slug" json:"size_slug"`
	Size        string    `gorm:"not null" validate:"required" json:"size"`
	Price       float64   `gorm:"type:numeric(10,2);not null" validate:"required,numeric" json:"price"`
	Discount    float64   `gorm:"type:numeric" validate:"numeric" json:"discount"`
	CreatedAt   time.Time `gorm:"default:now()" json:"created_at"`
	UpdatedAt   time.Time `gorm:"default:now()" json:"updated_at"`
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

type AdminProduct struct {
	ID       uint64        `json:"id"`
	ImageUrl string        `json:"image_url"`
	Name     string        `json:"name"`
	Category string        `json:"category"`
	Slug     string        `json:"slug"`
	Sizes    []ProductSize `json:"sizes"`
}

func (ProductSize) TableName() string {
	return "product_sizes"
}
func (Product) TableName() string {
	return "products"
}
