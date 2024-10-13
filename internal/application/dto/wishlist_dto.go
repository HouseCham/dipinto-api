package dto

import (
	"encoding/json"

	"github.com/HouseCham/dipinto-api/internal/domain/model"
)

type WishListDTO struct {
	ID               uint64                 `json:"id"`
	UserID           uint64                 `json:"user_id"`
	WishlistProducts []model.CatalogProduct `json:"wishlist_products"`
}

type WishlistProductDTO struct {
	ID          uint64          `json:"id"`
	Name        string          `json:"name" `
	Slug        string          `json:"slug"`
	Images      json.RawMessage `json:"images"`
	Description string          `json:"description"`
	Category    string          `json:"category"`
}