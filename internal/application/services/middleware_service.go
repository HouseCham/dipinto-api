package services

import (
	"time"

	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/gofiber/fiber/v3"
	"github.com/dgrijalva/jwt-go"
)

// MiddlewareService contains the business logic related to the middleware.
type MiddlewareService struct {
	jwtKey []byte
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
		tokenStr := c.Cookies("dipinto-token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
				Message:   "No token provided",
				StatusCode: fiber.StatusUnauthorized,
			})
		}

		// Parse and validate the token
		claims, err := m.middleware.ValidateToken(tokenStr)
		if err != nil {
			// Check if the token is expired and the user has the remember flag set to true
			if (claims != nil && claims.Remember) {
				// Generate a new token with the same claims
				newToken := refreshToken(claims, m.jwtKey)
				// Set the new token in the response
				c.Cookie(&fiber.Cookie{
					Name:     "dipinto-token",
					Value:    newToken,
					Expires:  time.Now().Add(24 * time.Hour),
					HTTPOnly: true,
				})
			} else {
				return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
					Message:  err.Error(),
					StatusCode: fiber.StatusUnauthorized,
				})
			}
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

// VerifyRole is a Fiber middleware function that checks if the user has the required role
func (m *MiddlewareService) VerifyAdmin() fiber.Handler {
	return func(c fiber.Ctx) error {
		claims, ok := c.Locals("claims").(*middleware.Claims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
				Message:    "Invalid token claims",
				StatusCode: fiber.StatusUnauthorized,
			})
		}
		if (claims.Role == "" || claims.Role != "admin") {
			return c.Status(fiber.StatusForbidden).JSON(model.HTTPResponse{
				Message:    "Insufficient permissions",
				StatusCode: fiber.StatusForbidden,
			})
		}
		return c.Next()
	}
}
// refreshToken generates a new JWT token with the same claims as the old token but with a new expiration time
func refreshToken(claims *middleware.Claims, jwtKey []byte) string {
	// Set the expiration time of the token
	expirationTime := time.Now().Add(8 * time.Hour) // Token expires in 8 hours

	// Create the claims
	refreshClaims := &middleware.Claims{
		ID:       claims.ID,
		Username: claims.Username,
		Role:     claims.Role,
		Remember: claims.Remember,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "dipinto-api",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(jwtKey)
	if err != nil {
		return ""
	}

	return signedToken
}