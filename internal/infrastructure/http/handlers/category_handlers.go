package handlers

import (
	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
)

type CategoryHandler struct {
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
	ModelService      *services.ModelService
}

// GetAllCategories is a handler function that retrieves all categories from the database
func (h *CategoryHandler) GetAllCategories(c fiber.Ctx) error {
	categories, err := h.RepositoryService.GetAllCategories()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve categories from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Categories retrieved successfully",
		Data:       categories,
	})
}

// InsertCategory is a handler function that inserts a new category into the database
func (h *CategoryHandler) InsertCategory(c fiber.Ctx) error {
	var request model.Category
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing category request body",
		})
	}

	// Validate the request body
	invalidResponse := h.ModelService.ValidateRequestBody(&request)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       invalidResponse,
		})
	}

	// Check if the category already exists
	categoryExists, err := h.RepositoryService.CheckCategoryExists(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to check if category exists in the database",
		})
	} else if categoryExists {
		return c.Status(fiber.StatusConflict).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusConflict,
			Message:    "Category already exists",
		})
	}

	// Insert the category into the database
	categoryID, err := h.RepositoryService.InsertCategory(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create category in the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Category created successfully",
		Data:       categoryID,
	})
}

// UpdateCategory is a handler function that updates a category in the database
func (h *CategoryHandler) UpdateCategory(c fiber.Ctx) error {
	var request model.Category
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing category request body",
		})
	}

	// Validate the request body
	invalidResponse := h.ModelService.ValidateRequestBody(&request)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       invalidResponse,
		})
	}

	// Update the category in the database
	err := h.RepositoryService.UpdateCategory(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update category in the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Category updated successfully",
	})
}