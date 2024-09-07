package routes

import (
	"github.com/HouseCham/dipinto-api/internal/infrastructure/http/handlers"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(app *fiber.App, userHandler *handlers.UserHandler, productHandler *handlers.ProductHandler) {
	/* ========== GLOBAL  ========== */
	// === MIDDLEWARE ===
	app.Use(userHandler.MiddlewareService.VerifyOrigin())
	// === HANDLERS ===
	app.Post("/api/v1/users", userHandler.InsertUser)
	app.Get("/api/v1/products", productHandler.GetAllProductsCatalogue)
	app.Post("/api/v1/users/login", userHandler.LoginUser)
	app.Post("/api/v1/users/sign-up", userHandler.InsertUser)

	/* ========== User routes  ========== */
	userRoutes := app.Group("/api/v1/users")
	// === MIDDLEWARE ===
	userRoutes.Use(userHandler.MiddlewareService.VerifyJWT())
	// === HANDLERS ===
	userRoutes.Get("/", userHandler.GetUserById)

	/* ========== CUSTOMER Product routes  ========== */
	// === GROUP ===
	productRoutes := app.Group("/api/v1/products")
	// === HANDLERS ===
	productRoutes.Get("/", productHandler.GetAllProductsCatalogue)
	productRoutes.Get("/:slug", productHandler.GetProductBySlug)

	/* ========== ADMIN Product routes  ========== */
	// === MIDDLEWARE ===
	productRoutes.Use(productHandler.MiddlewareService.VerifyJWT()).Use(productHandler.MiddlewareService.VerifyAdmin())
	// === HANDLERS ===
	productRoutes.Get("/get-products/admin", productHandler.GetAllProductsAdmin)
	productRoutes.Post("/insert", productHandler.InsertProduct)
}