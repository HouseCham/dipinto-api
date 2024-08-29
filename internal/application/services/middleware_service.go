package services

import (
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
)

// MiddlewareService contains the business logic related to the middleware.
type MiddlewareService struct {
	middleware *middleware.MiddlewareService
}

// NewMiddlewareService creates a new MiddlewareService
func NewMiddlewareService(m *middleware.MiddlewareService) *MiddlewareService {
	return &MiddlewareService{
		middleware: m,
	}
}

// VerifyJWT validates the JWT token and checks the claim ID is greater than 0
func (m *MiddlewareService) VerifyJWT() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the JWT token from the Authorization header
		tokenStr := c.Get("Authorization")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
				StatusCode: fiber.StatusUnauthorized,
			})
		}

		// Extract the token from the "Bearer " prefix
		if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
			tokenStr = tokenStr[7:]
		} else {
			return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
				StatusCode: fiber.StatusUnauthorized,
			})
		}

		// Parse and validate the token
		claims, err := m.middleware.ValidateToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
				StatusCode: fiber.StatusUnauthorized,
			})
		}

		// Store the claims in the request context
		c.Locals("claims", claims)

		// Proceed to the next middleware/handler
		return c.Next()
	}
}

// VerifyOrigin is a Fiber middleware function that checks if the request's Origin header is from "ramsesramva.com"
func (m *MiddlewareService) VerifyOrigin() fiber.Handler {
	return func(c fiber.Ctx) error {
		// Get the Origin header
		origin := c.Get("Origin")
		// Validate the Origin header
		if !(m.middleware.ValidateOrigin(origin)) {
			return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
				StatusCode: fiber.StatusUnauthorized,
			})
		}
		// Proceed to the next middleware/handler
		return c.Next()
	}
}