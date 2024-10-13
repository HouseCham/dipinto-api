package dto

import "encoding/json"

type CartDTO struct {
	ID           uint64           `json:"id"`
	CartProducts []CartProductDTO `gorm:"-" json:"cart_products"`
}

type CartProductDTO struct {
	ID       uint64          `json:"id"`
	Name     string          `json:"name"`
	Slug     string          `json:"slug"`
	Size     string          `json:"size"`
	Price    float64         `json:"price"`
	Discount float64         `json:"discount"`
	Quantity int             `json:"quantity"`
	Images   json.RawMessage `json:"images"`
}
