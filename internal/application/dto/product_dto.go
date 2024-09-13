package dto

import "encoding/json"

// ProductDTO is a data transfer object for the product model
type ProductDTO struct {
	ID          uint64           `json:"id"`
	CategoryID  uint64           `json:"category_id"`
	Slug        string           `json:"slug"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	Sizes       []ProductSizeDTO `json:"sizes"`
	Images      []ImageDTO       `json:"images"`
}

// AdminProduct is a struct that represents a product in the admin panel catalog
type AdminProductDTO struct {
	ID       uint64           `gorm:"primaryKey" json:"id"`
	Images   json.RawMessage  `gorm:"type:jsonb" json:"images"`
	Name     string           `json:"name"`
	Category string           `json:"category"`
	Slug     string           `json:"slug"`
	Sizes    []ProductSizeDTO `gorm:"-" json:"sizes"`
}

// ProductSizeDTO is a data transfer object for the product size model
type ProductSizeDTO struct {
	ID          uint64  `json:"id"`
	ProductID   uint64  `json:"product_id"`
	IsAvailable bool    `json:"is_available"`
	SizeSlug    string  `json:"size_slug"`
	Size        string  `json:"size"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
	UpdatedAt   string  `json:"updated_at"`
	DeletedAt   string  `json:"deleted_at"`
}

type ImageDTO struct {
	URL       string `json:"url" validate:"required"`
	Alt       string `json:"alt" validate:"required"`
	IsPrimary bool   `json:"is_primary" validate:"boolean"`
}
