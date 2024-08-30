package dto

// ProductDTO is a data transfer object for the product model
type ProductDTO struct {
	ID          uint64        `json:"id"`
	CategoryID  uint64        `json:"category_id"`
	Slug        string        `json:"slug"`
	Name        string        `json:"name"`
	Description string        `json:"description"`
	Sizes       []ProductSizeDTO     `json:"sizes"`
	Images      []ImageDTO `json:"images"`
}

// ProductSizeDTO is a data transfer object for the product size model
type ProductSizeDTO struct {
	ID          uint64  `json:"id"`
	IsAvailable bool    `json:"is_available"`
	SizeSlug    string  `json:"size_slug"`
	Size        string  `json:"size"`
	Price       float64 `json:"price"`
	Discount    float64 `json:"discount"`
}

type ImageDTO struct {
	URL       string `json:"url" validate:"required"`
	Alt       string `json:"alt" validate:"required"`
	IsPrimary bool   `json:"is_primary" validate:"boolean"`
}