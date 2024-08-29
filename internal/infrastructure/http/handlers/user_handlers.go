package handlers

import (
	"strconv"

	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type UserHandler struct {
	AuthService       *services.AuthService
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
	ModelService      *services.ModelService
}

// InsertUser is a handler function that inserts a new user into the database
func (h *UserHandler) InsertUser(c fiber.Ctx) error {
	var request model.User
	if err := c.Bind().JSON(&request); err != nil {
		log.Warnf("Failed to parse user request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing user request body",
		})
	}
	// Validate the request body
	if errors := h.ModelService.ValidateRequestBody(request); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       errors,
		})
	}

	// Hash the user's password
	hashedPassword, err := h.AuthService.HashPassword(request.Password)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to hash user password",
		})
	}
	request.Password = hashedPassword

	// Insert the user into the database
	userID, err := h.RepositoryService.InsertUser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create user in the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "User created successfully",
		Data:       userID,
	})
}
// GetUserById is a handler function that retrieves a user from the database by ID
func (h *UserHandler) GetUserById(c fiber.Ctx) error {
	idStr := c.Params("id")
	// convert to uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		log.Warnf("Failed to parse user ID: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid user ID",
		})
	}
	// Retrieve the user from the database
	user, err := h.RepositoryService.GetUserById(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve user from the database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "User retrieved successfully",
		Data:       user,
	})
}
