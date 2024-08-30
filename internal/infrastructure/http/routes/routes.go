package routes

import (
	"github.com/HouseCham/dipinto-api/internal/infrastructure/http/handlers"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, productHandler *handlers.ProductHandler) {
	// user routes
	app.Use(userHandler.MiddlewareService.VerifyOrigin())
	app.Post("/api/v1/users", userHandler.InsertUser)

	userRoutes := app.Group("/api/v1/users")
	userRoutes.Post("/login", userHandler.LoginUser)
	userRoutes.Use(userHandler.MiddlewareService.VerifyJWT())
	// Get by ID
	userRoutes.Get("/", userHandler.GetUserById)

	/* ========== Product routes  ========== */
	// === ADMIN ===
	productAdminRoutes := app.Group("/api/v1/products")
	productAdminRoutes.Use(productHandler.MiddlewareService.VerifyJWT()).Use(productHandler.MiddlewareService.VerifyAdmin())
	productAdminRoutes.Post("/", productHandler.InsertProduct)
	// === CUSTOMER ===
}