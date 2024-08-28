package routes

import (
	"github.com/HouseCham/dipinto-api/internal/infrastructure/http/handlers"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler) {
	// user routes
	app.Post("/api/v1/users", userHandler.InsertUser)

	userRoutes := app.Group("/api/v1/users")
	userRoutes.Use(userHandler.MiddlewareService.VerifyJWT())
	userRoutes.Get("/", userHandler.GetUser)
}