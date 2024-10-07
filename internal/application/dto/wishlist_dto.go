package dto

import "github.com/HouseCham/dipinto-api/internal/domain/model"

type WishListDTO struct {
	ID               uint64                 `json:"id"`
	UserID           uint64                 `json:"user_id"`
	WishlistProducts []model.CatalogProduct `json:"wishlist_products"`
}