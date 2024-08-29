package main

import (
	"fmt"

	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/auth"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/db"
	"github.com/HouseCham/dipinto-api/internal/domain/services"
	"github.com/HouseCham/dipinto-api/internal/infrastructure/config"
	"github.com/HouseCham/dipinto-api/internal/infrastructure/http/handlers"
	"github.com/HouseCham/dipinto-api/internal/infrastructure/http/routes"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

func main() {
	log.Info("Starting Dipinto API")
	// Set up the configuration
	cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load configuration: %v", err)
    }
	log.Info("Configuration is set up")

	// Set up database connection
	database, err := db.NewDBConn(cfg.Database.DNS)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}
	defer database.Close()
	log.Info("Database connection is set up")
	
	// Set up the Fiber app
	app := fiber.New()

	// Setting up CORS middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{cfg.Client.Origin},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
	}))
	log.Info("Fiber app is set up")

	userHandler := *injectDependencies(cfg, database)
	log.Info("Handlers are set up")
	
	// Set up the routes and handlers for the app
	routes.SetupRoutes(app, &userHandler)
	log.Info("Routes are set up")

	log.Infof("Server is running on port %d", cfg.Server.Port)
	app.Listen(fmt.Sprintf(":%d", cfg.Server.Port))
}

// injectDependencies injects the dependencies into the handlers
func injectDependencies(cfg *config.Config, database *db.Database) *handlers.UserHandler {
	// Set up the services for dependency injection
	authService := services.NewAuthService(auth.SetUpAuthService(cfg))
	middlewareService := services.NewMiddlewareService(middleware.SetupMiddlewareService(cfg))
	repositoryService := services.NewRepositoryService(database)
	// Set up the http handlers
	return &handlers.UserHandler{
		AuthService: authService,
		MiddlewareService: middlewareService,
		RepositoryService: repositoryService,
	}
}