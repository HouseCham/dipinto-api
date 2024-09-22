package handlers

import (
	"strconv"
	"time"

	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

const CookieSecure = false

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
	claims, ok := c.Locals("claims").(*middleware.Claims)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve user claims",
		})
	}
	log.Infof("Claims: %v", claims)

	// parse string id to uint64
	userID, err := strconv.ParseUint(claims.ID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid user ID",
		})
	}
	// Retrieve the user from the database
	user, err := h.RepositoryService.GetUserById(userID)
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

// LoginUser is a handler function that logs a user into the application
func (h *UserHandler) LoginUser(c fiber.Ctx) error {
	var request model.LoginUser
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
	// Retrieve the user from the database
	dbUser, err := h.RepositoryService.GetUserByEmail(request.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve user information",
		})
	}

	// Compare the user's password with the hashed password in the database
	if err := h.AuthService.ValidatePassword(request.Password, dbUser.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Incorrect email or password",
		})
	}

	// Generate a JWT token
	token, err := h.AuthService.GenerateToken(dbUser)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to generate token",
		})
	}

	// create http only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		HTTPOnly: true,
		Secure:   CookieSecure,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	// Return the response with cookie
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "User logged in successfully",
		Data:       token,
	})
}

// LogoutUser is a handler function that removes the http only cookie from the client
func (h *UserHandler) LogoutUser(c fiber.Ctx) error {
	// create http only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		HTTPOnly: true,
		Secure:   CookieSecure,
	})

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "User logged out successfully",
	})
}

// GetAllCustomers is a handler function that retrieves all customers from the database
func (h *UserHandler) GetAllCustomers(c fiber.Ctx) error {
	offset := c.Query("offset", "0")
	limit := c.Query("limit", "10")
	searchValue := c.Query("searchValue", "")

	// Convert the offset and limit query parameters to integers
	offsetInt, err := strconv.Atoi(offset)
	limitInt, err1 := strconv.Atoi(limit)
	if err != nil || err1 != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid pagination parameters",
		})
	}

	// Retrieve all customers from the database
	customers, err := h.RepositoryService.GetAllCustomers(offsetInt, limitInt, searchValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve customers from the database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Customers retrieved successfully",
		Data:       customers,
	})
}
