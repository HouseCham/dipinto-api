package handlers

import (
	"github.com/HouseCham/dipinto-api/internal/application/dto"
	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
)

type AdminHandler struct {
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
}

// GetAdminDashboard is a handler function that retrieves the admin dashboard data
func (h *AdminHandler) GetAdminDashboard(c fiber.Ctx) error {
	cards, err := h.RepositoryService.GetAdminCardsData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.AdminDashboardDTO{})
	}

	salesChart, err := h.RepositoryService.GetMonthlySalesData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.AdminDashboardDTO{})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Admin dashboard data retrieved successfully",
		Data: dto.AdminDashboardDTO{
			Cards:      cards,
			SalesChart: salesChart,
		},
	})
}