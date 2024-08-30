package handlers

import (
	"github.com/HouseCham/dipinto-api/internal/application/dto"
	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/HouseCham/dipinto-api/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type ProductHandler struct {
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
	ModelService      *services.ModelService	
}

// InsertProduct is a handler function that inserts a new product into the database
func (h *ProductHandler) InsertProduct(c fiber.Ctx) error {
	var requestDto dto.ProductDTO

	err := c.Bind().JSON(&requestDto); 
	if err != nil {
		log.Warnf("Failed to parse user request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing user request body",
		})
	}

	// Parse the request DTO to a model
	product, sizes := utils.ParseProductDTOToModel(requestDto)

	// Validate the request body
	invalidResponse := validateProductStruct(product, sizes, requestDto.Images, h)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(invalidResponse)
	}

	// Insert the product into the database
	productID, err := h.RepositoryService.InsertProduct(product, sizes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create product in the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Product created successfully",
		Data:       productID,
	})
}

// GetAllProducts is a handler function that retrieves all products from the database
func (h *ProductHandler) GetAllProductsCatalogue(c fiber.Ctx) error {
	products, err := h.RepositoryService.GetAllProductsCatalogue()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve products from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		Message:    "Products retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       products,
	})
}

// validateProductStruct isolates the validation logic for the product struct
func validateProductStruct(product *model.Product, sizes *[]model.ProductSize, imagesDto []dto.ImageDTO, h *ProductHandler) *model.HTTPResponse {
	if errors := h.ModelService.ValidateRequestBody(product); errors != nil {
		return utils.ReturnBadRequestResponse(errors)
	}
	for _, size := range *sizes {
		if errors := h.ModelService.ValidateRequestBody(size); errors != nil {
			return utils.ReturnBadRequestResponse(errors)
		}
	}
	for _, image := range imagesDto {
		if errors := h.ModelService.ValidateRequestBody(image); errors != nil {
			return utils.ReturnBadRequestResponse(errors)
		}
	}
	return nil
}