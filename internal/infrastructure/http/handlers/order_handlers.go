package handlers

import (
	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
)

type OrderHandler struct {
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
	ModelService      *services.ModelService
}

// GetAdminOrderList is a handler function that retrieves all orders from the database
func (h *OrderHandler) GetAdminOrderList(c fiber.Ctx) error {
	orders, err := h.RepositoryService.GetAdminOrderList()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve orders from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Orders retrieved successfully",
		Data:       orders,
	})
}