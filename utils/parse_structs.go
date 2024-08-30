package utils

import (
	"encoding/json"

	"github.com/HouseCham/dipinto-api/internal/application/dto"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/validator"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
)

// ParseProductDTOToModel parses a ProductDTO to a Product model
func ParseProductDTOToModel(dto dto.ProductDTO) (*model.Product, *[]model.ProductSize) {
	// Marshal images to JSON
    imagesJSON, err := json.Marshal(dto.Images)
    if err != nil {
        return nil, nil
    }
	// Create a slice of ProductSize
	var sizes []model.ProductSize
	for _, size := range dto.Sizes {
		sizes = append(sizes, model.ProductSize{
			ID:          size.ID,
			ProductID:   dto.ID,
			IsAvailable: size.IsAvailable,
			SizeSlug:    size.SizeSlug,
			Size:        size.Size,
			Price:       size.Price,
			Discount:    size.Discount,
		})
	}
	// Return the Product and the slice of ProductSize
	return &model.Product{
		ID:          dto.ID,
		CategoryID:  dto.CategoryID,
		Slug:        dto.Slug,
		Name:        dto.Name,
		Description: dto.Description,
		Images:      imagesJSON,
	}, &sizes
}

// returnBadRequestResponse isolates the response logic for a bad request
func ReturnBadRequestResponse(errors *[]validator.ValidatorError) *model.HTTPResponse {
	return &model.HTTPResponse{
		StatusCode: fiber.StatusBadRequest,
		Message:    "Invalid request body",
		Data:       errors,
	}
}
