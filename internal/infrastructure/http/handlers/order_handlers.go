package handlers

import (
	"strconv"

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
	offset := c.Query("offset", "0")
	limit := c.Query("limit", "10")
	searchValue := c.Query("searchValue", "")
	status := c.Query("status", "")
	payment := c.Query("payment", "")

	offsetInt, err := strconv.Atoi(offset)
	limitInt, err1 := strconv.Atoi(limit)
	if (err != nil || err1 != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid pagination parameters",
		})
	}
	orders, err := h.RepositoryService.GetAdminOrderList(offsetInt, limitInt, searchValue, status, payment)
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

// GetAdminOrderDetails is a handler function that retrieves a specific order from the database
func (h *OrderHandler) GetAdminOrderDetailsById(c fiber.Ctx) error {
	orderID := c.Params("id")
	// convert the orderID to a uint64
	uintID, err := strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid order ID",
		})
	}
	// retrieve the order details from the database
	order, err := h.RepositoryService.GetOrderDetails(uintID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve order details from the database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Order details retrieved successfully",
		Data:       order,
	})
}