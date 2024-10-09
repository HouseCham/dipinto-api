package handlers

import (
	"strconv"
	"time"

	"github.com/HouseCham/dipinto-api/internal/application/services"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/HouseCham/dipinto-api/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

type ClientHandler struct {
	MiddlewareService *services.MiddlewareService
	RepositoryService *services.RepositoryService
	ModelService      *services.ModelService
	AuthService       *services.AuthService
}

// RenewToken is a handler function that renews the user's token
func (h *ClientHandler) RefreshToken(c fiber.Ctx) error {
	// get claims from context
	claims := c.Locals("claims").(*middleware.Claims)

	user := &model.User{
		ID:   utils.StringToUint64(claims.ID), // Ensure this function exists in the utils package
		Name: claims.Username,
		Role: claims.Role,
	}
	// generate new token
	token, err := h.AuthService.GenerateToken(user, true)
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

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Token renewed successfully",
	})
}

// #region Categories
// GetAllCategories is a handler function that retrieves all categories from the database
func (h *ClientHandler) GetAllCategories(c fiber.Ctx) error {
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

	// Retrieve all categories from the database
	categories, err := h.RepositoryService.GetAllCategories(offsetInt, limitInt, searchValue)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve categories from the database",
		})
	}

	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Categories retrieved successfully",
		Data:       categories,
	})
}

//#endregion Categories

// #region Products
// GetAllProducts is a handler function that retrieves all products from the database
func (h *ClientHandler) GetAllProductsCatalog(c fiber.Ctx) error {
	searchParams := model.ClientSearchParams{
		SearchValue: c.Query("search", ""),
		CategoryID:  utils.StringToUint64(c.Query("category_id", "0")),
		MinPrice:    utils.StringToFloat64(c.Query("price_min", "0")),
		MaxPrice:    utils.StringToFloat64(c.Query("price_limit", "0")),
		OrderBy:     c.Query("sort_by", "created_at"),
		Offset:      utils.StringToInt(c.Query("offset", "0")),
		Limit:       utils.StringToInt(c.Query("limit", "10")),
	}

	products, err := h.RepositoryService.GetAllProductsCatalog(searchParams)
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

// GetProductBySlug is a handler function that retrieves a product by its slug
func (h *ClientHandler) GetProductBySlug(c fiber.Ctx) error {
	slug := c.Params("slug")
	if slug == "" {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid slug",
		})
	}
	product, sizes, err := h.RepositoryService.GetProductBySlug(slug)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusNotFound,
			Message:    "Product not found",
		})
	}
	// Parse the product model to a DTO
	productDto := utils.ParseProductModelToDTO(product, sizes)
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		Message:    "Product retrieved successfully",
		StatusCode: fiber.StatusOK,
		Data:       productDto,
	})
}

//#endregion Products

// #region Customers
// InsertCustomer is a handler function that inserts a new customer into the database
func (h *ClientHandler) InsertCustomer(c fiber.Ctx) error {
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
	request.Role = "customer"

	// Insert the user into the database
	userID, err := h.RepositoryService.InsertUser(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create user in the database",
		})
	}

	// Create a new cart for the user
	_, err = h.RepositoryService.InsertCart(&model.Cart{
		UserID: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create cart in the database",
		})
	}

	// Create a new wishlist for the user
	_, err = h.RepositoryService.InsertWishlist(&model.Wishlist{
		UserID: userID,
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create wishlist in the database",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "User created successfully",
		Data:       userID,
	})
}

// LoginCustomer is a handler function that logs a customer into the application
func (h *ClientHandler) LoginCustomer(c fiber.Ctx) error {
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
	dbUser, err := h.RepositoryService.GetUserByEmail(request.Email, false)
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
	token, err := h.AuthService.GenerateToken(dbUser, request.Remember)
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

	userLogged := model.User{
		Name:  dbUser.Name,
		Email: dbUser.Email,
		Phone: dbUser.Phone,
	}

	// Return the response with cookie
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "User logged in successfully",
		Data:       userLogged,
	})
}

//#endregion Customers

// #region Wishlist
// GetCustomerWishlist is a handler function that retrieves the customer's wishlist
func (h *ClientHandler) GetCustomerWishlist(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	userID := utils.StringToUint64(claims.ID)
	// Retrieve the user's wishlist from the database
	wishlist, err := h.RepositoryService.GetWishlistByUserId(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve wishlist information",
		})
	}
	// Retrieve the wishlist products from the database
	wishlistProducts, err := h.RepositoryService.GetWishlistProductsById(wishlist.ID, userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve wishlist products information",
		})
	}

	wishlist.WishlistProducts = *wishlistProducts
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Wishlist retrieved successfully",
		Data:       utils.ParseWishlistToDTO(wishlist),
	})
}

// AddProductToWishlist is a handler function that adds a product to the wishlist
func (h *ClientHandler) AddProductToWishlist(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	var request model.WishlistProduct
	if err := c.Bind().JSON(&request); err != nil {
		log.Warnf("Failed to parse wishlist product request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing wishlist product request body",
		})
	}
	// Get Wishlist from the database using the user ID
	wishlist, err := h.RepositoryService.GetWishlistByUserId(utils.StringToUint64(claims.ID))
	// Check if the wishlist exists
	if err != nil {
		// If the wishlist does not exist, create a new wishlist
		wishlistID, err := h.RepositoryService.InsertWishlist(&model.Wishlist{
			UserID: utils.StringToUint64(claims.ID),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
				StatusCode: fiber.StatusInternalServerError,
				Message:    "Failed to create wishlist in the database",
			})
		}
		wishlist.ID = wishlistID
	}

	wishlistProduct := model.WishlistProduct{
		WishlistID: wishlist.ID,
		ProductID:  request.ProductID,
	}
	// Insert the product into the wishlist
	_, err = h.RepositoryService.InsertWishlistProduct(&wishlistProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to add product to wishlist",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Product added to wishlist successfully",
	})
}

// RemoveProductFromWishlist is a handler function that removes a product from the wishlist
func (h *ClientHandler) RemoveProductFromWishlist(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)

	productID := utils.StringToUint64(c.Params("id"))
	userID := utils.StringToUint64(claims.ID)

	// get the wishlist from the database
	wishlist, err := h.RepositoryService.GetWishlistByUserId(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve wishlist information",
		})
	}

	// Remove the product from the wishlist
	err = h.RepositoryService.DeleteWishlistProduct(productID, wishlist.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to remove product from wishlist",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Product removed from wishlist successfully",
	})
}

//#endregion Wishlist

// #region Addresses
// InsertAddress is a handler function that inserts a new address into the database
func (h *ClientHandler) InsertCustomerAddress(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	var request model.Address
	if err := c.Bind().JSON(&request); err != nil {
		log.Warnf("Failed to parse address request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing address request body",
		})
	}
	request.UserID = utils.StringToUint64(claims.ID)
	request.Country = "México"
	// Validate the request body
	if errors := h.ModelService.ValidateRequestBody(request); errors != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Invalid request body",
			Data:       errors,
		})
	}
	// Insert the address into the database
	addressID, err := h.RepositoryService.InsertAddress(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create address in the database",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Address created successfully",
		Data:       addressID,
	})
}

// GetCustomerAddresses is a handler function that retrieves all addresses from the database
func (h *ClientHandler) GetCustomerAddresses(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	// Retrieve all addresses from the database
	addresses, err := h.RepositoryService.GetAddressesListByUserId(utils.StringToUint64(claims.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve addresses from the database",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Addresses retrieved successfully",
		Data:       addresses,
	})
}

//#endregion Customers

// #region Carts
// GetCustomerCart is a handler function that retrieves the customer's cart
func (h *ClientHandler) GetCustomerCart(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	// Retrieve the user's cart from the database
	cart, err := h.RepositoryService.GetCartByUserId(utils.StringToUint64(claims.ID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve cart information",
		})
	}
	// Retrieve the cart products from the database
	cartProducts, err := h.RepositoryService.GetCartProductsByCartId(cart.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve cart products information",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Cart retrieved successfully",
		Data:       cartProducts,
	})
}

// InsertCart is a handler function that inserts a new cart into the database
func (h *ClientHandler) InsertCart(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	request := model.Cart{
		UserID: utils.StringToUint64(claims.ID),
	}
	// Insert the cart into the database
	cartID, err := h.RepositoryService.InsertCart(&request)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to create cart in the database",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Cart created successfully",
		Data:       cartID,
	})
}

// AddProductToCart is a handler function that adds a product to the cart
func (h *ClientHandler) AddProductToCart(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	var request model.CartProduct
	if err := c.Bind().JSON(&request); err != nil {
		log.Warnf("Failed to parse cart product request body: %v", err)
		return c.Status(fiber.StatusBadRequest).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusBadRequest,
			Message:    "Error parsing cart product request body",
		})
	}
	// Get Cart from the database using the user ID
	cart, err := h.RepositoryService.GetCartByUserId(utils.StringToUint64(claims.ID))
	// Check if the cart exists
	if err != nil {
		// If the cart does not exist, create a new cart
		cartID, err := h.RepositoryService.InsertCart(&model.Cart{
			UserID: utils.StringToUint64(claims.ID),
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
				StatusCode: fiber.StatusInternalServerError,
				Message:    "Failed to create cart in the database",
			})
		}
		cart.ID = cartID
	}

	cartProduct := model.CartProduct{
		CartID:    cart.ID,
		ProductID: request.ProductID,
		Quantity:  request.Quantity,
	}
	// Insert the product into the cart
	_, err = h.RepositoryService.InsertCartProduct(&cartProduct)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to add product to cart",
		})
	}
	return c.Status(fiber.StatusCreated).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusCreated,
		Message:    "Product added to cart successfully",
	})
}

// RemoveProductFromCart is a handler function that removes a product from the cart
func (h *ClientHandler) RemoveProductFromCart(c fiber.Ctx) error {
	claims := c.Locals("claims").(*middleware.Claims)
	productID := utils.StringToUint64(c.Params("id"))
	userID := utils.StringToUint64(claims.ID)

	// get the cart from the database
	cart, err := h.RepositoryService.GetCartByUserId(userID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to retrieve cart information",
		})
	}
	// Remove the product from the cart
	err = h.RepositoryService.DeleteCartProduct(productID, cart.ID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.HTTPResponse{
			StatusCode: fiber.StatusInternalServerError,
			Message:    "Failed to remove product from cart",
		})
	}
	return c.Status(fiber.StatusOK).JSON(model.HTTPResponse{
		StatusCode: fiber.StatusOK,
		Message:    "Product removed from cart successfully",
	})
}
