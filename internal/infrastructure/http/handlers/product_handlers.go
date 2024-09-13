package handlers

import (
	"strconv"

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

// UpdateProduct is a handler function that updates a product in the database
func (h *ProductHandler) UpdateProduct(c fiber.Ctx) error {
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

	// Update the product in the database
	err = h.RepositoryService.UpdateProduct(product, sizes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update product in the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Product updated successfully",
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
// GetAllProductsAdmin is a handler function that retrieves all products from the database
func (h *ProductHandler) GetAllProductsAdmin(c fiber.Ctx) error {
	searchValue := c.Query("searchValue", "")
	categoryID := c.Query("categoryID", "0")
	offset := c.Query("offset", "0")
	limit := c.Query("limit", "10")
	// convert the orderID to a uint64
	uintID, err := strconv.ParseUint(categoryID, 10, 64)
	offsetInt, err1 := strconv.Atoi(offset)
	limitInt, err2 := strconv.Atoi(limit)
	if (err != nil || err1 != nil || err2 != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid order ID",
		})
	}
	// retrieve the order details from the database
	products, err := h.RepositoryService.GetAllProductsAdmin(uintID, searchValue, offsetInt, limitInt)
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

// GetProductBySlug is a handler function that retrieves a product by its slug
func (h *ProductHandler) GetProductBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid slug",
		})
	}
	product, sizes, err := h.RepositoryService.GetProductBySlug(slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    "Product not found",
		})
	}
	// Parse the product model to a DTO
	productDto := utils.ParseProductModelToDTO(product, sizes)
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		Message:    "Product retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       productDto,
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