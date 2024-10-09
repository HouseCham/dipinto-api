package routes

import (
	"github.com/HouseCham/dipinto-api/internal/infrastructure/http/handlers"
	"github.com/gofiber/fiber/v3"
)

// SetupRoutes sets up the routes for the application
func SetupRoutes(app *fiber.App, adminHandler *handlers.AdminHandler, clientHandler *handlers.ClientHandler) {
	/* ========== GLOBAL  ========== */
	app.Post("/api/v1/users/logout", adminHandler.LogoutUser)
	
	/* ========== CLIENT  ========== */
	clientRoutes := app.Group("/api/v1")
	// === CUSTOMER ENDPOINTS ===
	clientRoutes.Post("/customers/login", clientHandler.LoginCustomer)
	clientRoutes.Post("/customers/sign-up", clientHandler.InsertCustomer)
	// === PRODUCTS ENDPOINTS ===
	clientRoutes.Get("/products", clientHandler.GetAllProductsCatalog)
	clientRoutes.Get("/categories", clientHandler.GetAllCategories)
	clientRoutes.Get("/products/:slug", clientHandler.GetProductBySlug)
	
	/* ========== CUSTOMER JWT ENDPOINTS  ========== */
	clientRoutes.Use(adminHandler.MiddlewareService.VerifyJWT())
	// === RENEW TOKEN ===
	clientRoutes.Post("/users/refresh-token", clientHandler.RefreshToken)
	// === ADDRESS ENDPOINTS ===
	clientRoutes.Post("/user/address/insert", clientHandler.InsertCustomerAddress)
	clientRoutes.Get("/user/address", clientHandler.GetCustomerAddresses)
	// === WISHLIST ENDPOINTS ===
	clientRoutes.Get("/wishlist", clientHandler.GetCustomerWishlist)
	clientRoutes.Post("/wishlist/add-product", clientHandler.AddProductToWishlist)
	clientRoutes.Delete("/wishlist/remove-product/:id", clientHandler.RemoveProductFromWishlist)
	// === CART ENDPOINTS ===
	clientRoutes.Get("/cart", clientHandler.GetCustomerCart)
	clientRoutes.Post("/cart/add-product", clientHandler.AddProductToCart)
	clientRoutes.Delete("/cart/remove-product/:id", clientHandler.RemoveProductFromCart)
	
	/* ========== ADMIN  ========== */
	adminRoutes := app.Group("/api/v1")
	adminRoutes.Post("/users/login", adminHandler.LoginAdmin)
	// === MIDDLEWARE ===
	adminRoutes.Use(adminHandler.MiddlewareService.VerifyJWT())

	adminRoutes.Use(adminHandler.MiddlewareService.VerifyAdmin())
	adminRoutes.Get("/admin/dashboard", adminHandler.GetAdminDashboard)
	// === PRODUCTS ENDPOINTS ===
	adminRoutes.Get("/products/get-products/admin", adminHandler.GetAllProductsAdmin)
	adminRoutes.Post("/products/insert", adminHandler.InsertProduct)
	adminRoutes.Put("/products/update", adminHandler.UpdateProduct)
	// === CATEGORIES ENDPOINTS ===
	adminRoutes.Post("/categories/insert", adminHandler.InsertCategory)
	adminRoutes.Put("/categories/update", adminHandler.UpdateCategory)
	// === ORDERS ENDPOINTS ===
	adminRoutes.Get("/orders/get-admin-list", adminHandler.GetAdminOrderList)
	adminRoutes.Get("/orders/details/:id", adminHandler.GetAdminOrderDetailsById)
	// === CUSTOMERS ENDPOINTS ===
	adminRoutes.Get("/get-customers", adminHandler.GetAllCustomers)
}