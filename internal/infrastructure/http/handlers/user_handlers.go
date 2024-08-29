package handlers

import (
	"github.com/HouseCham/dipinto-api/internal/domain/services"
	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	AuthService *services.AuthService
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
}

func (h *UserHandler) InsertUser(c fiber.Ctx) error {
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Hello World",
	})
}

func (h *UserHandler) GetUser(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Hello World for GET user",
	})
}