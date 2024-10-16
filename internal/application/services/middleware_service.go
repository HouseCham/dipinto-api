package services

import (
	"github.com/HouseCham/dipinto-api/internal/domain/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
)

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