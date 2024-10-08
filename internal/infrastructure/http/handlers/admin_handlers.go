package handlers

import (
	"strconv"
	"time"

	"github.com/HouseCham/dipinto-api/internal/application/dto"
	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/HouseCham/dipinto-api/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type AdminHandler struct {
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
	ModelService      *services.ModelService
	AuthService       *services.AuthService
}

const CookieSecure bool = false

//#region Admin
// LoginUser is a handler function that logs a user into the application
func (h *AdminHandler) LoginAdmin(c fiber.Ctx) error {
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
	// Retrieve the admin from the database
	dbAdmin, err := h.RepositoryService.GetUserByEmail(request.Email, true)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve user information",
		})
	}

	// Compare the user's password with the hashed password in the database
	if err := h.AuthService.ValidatePassword(request.Password, dbAdmin.Password); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusUnauthorized,
			Message:    "Incorrect email or password",
		})
	}

	// Generate a JWT token
	token, err := h.AuthService.GenerateToken(dbAdmin, request.Remember)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to generate token",
		})
	}

	// create http only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "dipinto-token",
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
		Data:       "",
	})
}
// GetAllCustomers is a handler function that retrieves all customers from the database
func (h *AdminHandler) GetAllCustomers(c fiber.Ctx) error {
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
// GetUserById is a handler function that retrieves a user from the database by ID
func (h *AdminHandler) GetUserById(c fiber.Ctx) error {
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
// LogoutUser is a handler function that removes the http only cookie from the client
func (h *AdminHandler) LogoutUser(c fiber.Ctx) error {
	// create http only cookie
	c.Cookie(&fiber.Cookie{
		Name:     "dipinto-token",
		Value:    "",
		HTTPOnly: true,
		Secure:   CookieSecure,
	})

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "User logged out successfully",
	})
}
//#endregion Admin

//#region Categories
// InsertCategory is a handler function that inserts a new category into the database
func (h *AdminHandler) InsertCategory(c fiber.Ctx) error {
	var request model.Category
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing category request body",
		})
	}

	// Validate the request body
	invalidResponse := h.ModelService.ValidateRequestBody(&request)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       invalidResponse,
		})
	}

	// Check if the category already exists
	categoryExists, err := h.RepositoryService.CheckCategoryExists(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to check if category exists in the database",
		})
	} else if categoryExists {
		return c.Status(fiber.StatusConflict).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusConflict,
			Message:    "Category already exists",
		})
	}

	// Insert the category into the database
	categoryID, err := h.RepositoryService.InsertCategory(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create category in the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Category created successfully",
		Data:       categoryID,
	})
}

// UpdateCategory is a handler function that updates a category in the database
func (h *AdminHandler) UpdateCategory(c fiber.Ctx) error {
	var request model.Category
	if err := c.Bind().JSON(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing category request body",
		})
	}

	// Validate the request body
	invalidResponse := h.ModelService.ValidateRequestBody(&request)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       invalidResponse,
		})
	}

	// Update the category in the database
	err := h.RepositoryService.UpdateCategory(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update category in the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Category updated successfully",
	})
}
//#endregion Categories

//#region Products
// InsertProduct is a handler function that inserts a new product into the database
func (h *AdminHandler) InsertProduct(c fiber.Ctx) error {
	var requestDto dto.ProductDTO

	err := c.Bind().JSON(&requestDto); 
	if err != nil {
		log.Warnf("Failed to parse user request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing user request body",
		})
	}

	// Parse the request DTO to a model
	product, sizes := utils.ParseProductDTOToModel(requestDto)

	// Validate the request body
	invalidResponse := validateProductStruct(product, sizes, requestDto.Images, h)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(invalidResponse)
	}

	// Insert the product into the database
	productID, err := h.RepositoryService.InsertProduct(product, sizes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create product in the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Product created successfully",
		Data:       productID,
	})
}

// UpdateProduct is a handler function that updates a product in the database
func (h *AdminHandler) UpdateProduct(c fiber.Ctx) error {
	var requestDto dto.ProductDTO

	err := c.Bind().JSON(&requestDto); 
	if err != nil {
		log.Warnf("Failed to parse user request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing user request body",
		})
	}

	// Parse the request DTO to a model
	product, sizes := utils.ParseProductDTOToModel(requestDto)

	// Validate the request body
	invalidResponse := validateProductStruct(product, sizes, requestDto.Images, h)
	if invalidResponse != nil {
		return c.Status(fiber.StatusBadRequest).JSON(invalidResponse)
	}

	// Update the product in the database
	err = h.RepositoryService.UpdateProduct(product, sizes)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to update product in the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Product updated successfully",
	})
}
// GetAllProductsAdmin is a handler function that retrieves all products from the database
func (h *AdminHandler) GetAllProductsAdmin(c fiber.Ctx) error {
	searchValue := c.Query("searchValue", "")
	categoryID := c.Query("categoryID", "0")
	offset := c.Query("offset", "0")
	limit := c.Query("limit", "10")
	// convert the orderID to a uint64
	uintID, err := strconv.ParseUint(categoryID, 10, 64)
	offsetInt, err1 := strconv.Atoi(offset)
	limitInt, err2 := strconv.Atoi(limit)
	if (err != nil || err1 != nil || err2 != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid order ID",
		})
	}
	// retrieve the order details from the database
	products, err := h.RepositoryService.GetAllProductsAdmin(uintID, searchValue, offsetInt, limitInt)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve products from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		Message:    "Products retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       products,
	})
}

// validateProductStruct isolates the validation logic for the product struct
func validateProductStruct(product *model.Product, sizes *[]model.ProductSize, imagesDto []dto.ImageDTO, h *AdminHandler) *model.HTTPResponse {
	if errors := h.ModelService.ValidateRequestBody(product); errors != nil {
		return utils.ReturnBadRequestResponse(errors)
	}
	for _, size := range *sizes {
		if errors := h.ModelService.ValidateRequestBody(size); errors != nil {
			return utils.ReturnBadRequestResponse(errors)
		}
	}
	for _, image := range imagesDto {
		if errors := h.ModelService.ValidateRequestBody(image); errors != nil {
			return utils.ReturnBadRequestResponse(errors)
		}
	}
	return nil
}
//#endregion Products

// GetAdminDashboard is a handler function that retrieves the admin dashboard data
func (h *AdminHandler) GetAdminDashboard(c fiber.Ctx) error {
	cards, err := h.RepositoryService.GetAdminCardsData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.AdminDashboardDTO{})
	}

	salesChart, err := h.RepositoryService.GetMonthlySalesData()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(dto.AdminDashboardDTO{})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Admin dashboard data retrieved successfully",
		Data: dto.AdminDashboardDTO{
			Cards:      cards,
			SalesChart: salesChart,
		},
	})
}

//#region Orders
// GetAdminOrderList is a handler function that retrieves all orders from the database
func (h *AdminHandler) GetAdminOrderList(c fiber.Ctx) error {
	offset := c.Query("offset", "0")
	limit := c.Query("limit", "10")
	searchValue := c.Query("searchValue", "")
	status := c.Query("status", "")
	payment := c.Query("payment", "")

	offsetInt, err := strconv.Atoi(offset)
	limitInt, err1 := strconv.Atoi(limit)
	if (err != nil || err1 != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid pagination parameters",
		})
	}
	orders, err := h.RepositoryService.GetAdminOrderList(offsetInt, limitInt, searchValue, status, payment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve orders from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Orders retrieved successfully",
		Data:       orders,
	})
}

// GetAdminOrderDetails is a handler function that retrieves a specific order from the database
func (h *AdminHandler) GetAdminOrderDetailsById(c fiber.Ctx) error {
	orderID := c.Params("id")
	// convert the orderID to a uint64
	uintID, err := strconv.ParseUint(orderID, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid order ID",
		})
	}
	// retrieve the order details from the database
	order, err := h.RepositoryService.GetOrderDetails(uintID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve order details from the database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Order details retrieved successfully",
		Data:       order,
	})
}
//#endregion Orders