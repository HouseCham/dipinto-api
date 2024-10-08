package services

import (
	"encoding/json"
	"time"

	"github.com/HouseCham/dipinto-api/internal/application/dto"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/db"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/HouseCham/dipinto-api/utils"
	"github.com/gofiber/fiber/v3/log"
)

// RepositoryService contains the business logic related to the repository.
type RepositoryService struct {
	repo *db.Database
}

// NewRepositoryService creates a new RepositoryService
func NewRepositoryService(db *db.Database) *RepositoryService {
	return &RepositoryService{
		repo: db,
	}
}

// #region USER

// InsertUser inserts a new user into the database
func (r *RepositoryService) InsertUser(newUser *model.User) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newUser)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert user into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newUser.ID, nil
}

// GetAllUsers retrieves all users from the database
func (r *RepositoryService) GetAllUsers() ([]model.User, error) {
	var users []model.User
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Omit("CreatedAt", "UpdatedAt", "DeletedAt").Find(&users)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve users from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return users, nil
}

// GetAllCustomers retrieves all customers from the database
func (r *RepositoryService) GetAllCustomers(offset int, limit int, searchValue string) ([]model.User, error) {
	var customers []model.User
	query := r.repo.DB.Table("users u").
		Select("u.id, u.name, u.email, u.phone, u.created_at").
		Where("deleted_at IS NULL").
		Where("role='customer'")

	// Add search filter if searchValue is provided
	if searchValue != "" {
		searchPattern := "%" + searchValue + "%"
		query = query.Where("u.name ILIKE ?", searchPattern)
	}
	// Adding offset and limit for pagination
	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	dbResponse := query.Scan(&customers)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve customers from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return customers, nil
}

// GetUser retrieves a user from the database
func (r *RepositoryService) GetUserById(userID uint64) (*model.User, error) {
	var user model.User
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Omit("DeletedAt").First(&user, userID)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve user from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &user, nil
}

// GetUserByEmail retrieves a user from the database by email
func (r *RepositoryService) GetUserByEmail(email string, isAdmin bool) (*model.User, error) {
	var user model.User
	query := r.repo.DB.Table("users u").
		Select("u.id, u.name, u.email, u.phone, u.role, u.password").
		Where("u.deleted_at IS NULL").
		Where("u.email = ?", email)

	if isAdmin {
		query = query.Where("u.role = 'admin'")
	} else {
		query = query.Where("u.role = 'customer'")
	}

	dbResponse := query.Scan(&user)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve user from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &user, nil
}

// UpdateUser updates a user in the database
func (r *RepositoryService) UpdateUser(updatedUser *model.User) error {
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Save(updatedUser)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update user in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

//#region PRODUCT

// InsertProduct inserts a new product into the database
func (r *RepositoryService) InsertProduct(newProduct *model.Product, sizes *[]model.ProductSize) (uint64, error) {
	// Marshal images to JSON
	imagesJSON, err := json.Marshal(newProduct.Images)
	if err != nil {
		return 0, err
	}

	// Set the images field to the marshaled JSON
	newProduct.Images = imagesJSON

	// start a new transaction
	tx := r.repo.DB.Begin()
	if tx.Error != nil {
		log.Warnf("Failed to begin transaction: %v", tx.Error)
		return 0, tx.Error
	}

	// Insert the product
	dbResponse := tx.Omit("ID", "Category", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert product into the database: %v", dbResponse.Error)
		tx.Rollback()
		return 0, dbResponse.Error
	}

	// Insert the product sizes
	for _, size := range *sizes {
		size.ProductID = newProduct.ID
		dbResponse := tx.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&size)
		if dbResponse.Error != nil {
			log.Warnf("Failed to insert product size into the database: %v", dbResponse.Error)
			tx.Rollback()
			return 0, dbResponse.Error
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Warnf("Failed to commit transaction: %v", err)
		tx.Rollback()
		return 0, err
	}

	return newProduct.ID, nil
}

// UpdateProduct updates a product in the database
func (r *RepositoryService) UpdateProduct(updatedProduct *model.Product, sizes *[]model.ProductSize) error {
	// Marshal images to JSON
	imagesJSON, err := json.Marshal(updatedProduct.Images)
	if err != nil {
		return err
	}

	// Set the images field to the marshaled JSON
	updatedProduct.Images = imagesJSON
	updatedProduct.UpdatedAt = time.Now()

	// start a new transaction
	tx := r.repo.DB.Begin()
	if tx.Error != nil {
		log.Warnf("Failed to begin transaction: %v", tx.Error)
		return tx.Error
	}

	// Update the product
	dbResponse := tx.Omit("CreatedAt", "DeletedAt").Save(updatedProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update product in the database: %v", dbResponse.Error)
		tx.Rollback()
		return dbResponse.Error
	}

	// Update the product sizes
	for _, size := range *sizes {
		size.ProductID = updatedProduct.ID
		size.UpdatedAt = time.Now()
		dbResponse := tx.Omit("CreatedAt", "DeletedAt").Save(&size)
		if dbResponse.Error != nil {
			log.Warnf("Failed to update product size into the database: %v", dbResponse.Error)
			tx.Rollback()
			return dbResponse.Error
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Warnf("Failed to commit transaction: %v", err)
		tx.Rollback()
		return err
	}

	return nil
}

// GetAllProducts retrieves all products from the database
func (r *RepositoryService) GetAllProductsCatalog(params model.ClientSearchParams) (*[]model.CatalogProduct, error) {
	query := r.repo.DB.Table("products p").
		Select("p.id, p.slug, p.name, c.name as category, p.images, s.price, s.discount").
		Joins("INNER JOIN categories c ON p.category_id = c.id").
		Joins("INNER JOIN product_sizes s ON p.id = s.product_id").
		Joins("INNER JOIN (SELECT product_id, MIN(price) AS min_price FROM product_sizes WHERE is_available = true GROUP BY product_id) min_sizes ON s.product_id = min_sizes.product_id AND s.price = min_sizes.min_price").
		Where("s.is_available = true")

	// Add category filter if categoryId is provided
	if params.CategoryID != 0 {
		query = query.Where("p.category_id = ?", params.CategoryID)
	}

	// Add search filter if searchValue is provided
	if params.SearchValue != "" {
		searchPattern := "%" + params.SearchValue + "%"
		query = query.Where("p.name ILIKE ? OR p.description ILIKE ?", searchPattern, searchPattern)
	}

	// Add price filter if MinPrice and MaxPrice are provided
	if params.MinPrice > 0 {
		query = query.Where("s.price >= ?", params.MinPrice)
	}
	if params.MaxPrice > 0 {
		query = query.Where("s.price <= ?", params.MaxPrice)
	}

	// Adding offset and limit for pagination
	if params.Offset >= 0 {
		query = query.Offset(params.Offset)
	}
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}

	var products []model.CatalogProduct
	dbResponse := query.Scan(&products)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve products from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &products, nil
}

// GetAllProductsAdmin retrieves all products from the database for admin purposes
func (r *RepositoryService) GetAllProductsAdmin(categoryId uint64, searchValue string, offset int, limit int) (*[]dto.AdminProductDTO, error) {
	// get products from the database
	var products []dto.AdminProductDTO
	// Define the base query
	query := r.repo.DB.Table("products p").
		Select("p.id, p.slug, p.name, c.name as category, p.description, p.images").
		Joins("INNER JOIN categories c ON p.category_id = c.id").
		Where("p.deleted_at IS NULL")

	// Add category filter if categoryId is provided
	if categoryId != 0 {
		query = query.Where("p.category_id = ?", categoryId)
	}

	// Add search filter if searchValue is provided
	if searchValue != "" {
		searchPattern := "%" + searchValue + "%"
		query = query.Where("p.name ILIKE ? OR p.description ILIKE ?", searchPattern, searchPattern)
	}

	// Adding offset and limit for pagination
	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	dbResponse := query.Scan(&products)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve products from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	// get product sizes from the database
	var sizes []dto.ProductSizeDTO
	for i := range products {
		var productSizes []dto.ProductSizeDTO
		query = r.repo.DB.Table("product_sizes s").
			Select("s.id, s.product_id, s.is_available, s.size_slug, s.size, s.price, s.discount").
			Where("s.deleted_at IS NULL").
			Where("s.product_id = ?", products[i].ID)

		// perform query search
		dbResponse = query.Scan(&productSizes)
		if dbResponse.Error != nil {
			log.Warnf("Failed to retrieve product sizes from the database: %v", dbResponse.Error)
			return nil, dbResponse.Error
		}
		sizes = append(sizes, productSizes...)
	}

	// Assign sizes to products
	utils.AssignSizesToProducts(products, sizes)
	return &products, nil
}

// GetProductBySlug retrieves a product from the database by its slug
func (r *RepositoryService) GetProductBySlug(slug string) (*model.Product, *[]model.ProductSize, error) {
	// Retrieve the product from the database
	var product model.Product
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Where("slug = ?", slug).First(&product)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve product from the database: %v", dbResponse.Error)
		return nil, nil, dbResponse.Error
	}
	// retrieve the product sizes from the database
	var sizes []model.ProductSize
	dbResponse = r.repo.DB.Where("product_id = ?", product.ID).Omit("ProductID", "CreatedAt", "UpdatedAt", "DeletedAt").Find(&sizes)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve product sizes from the database: %v", dbResponse.Error)
		return nil, nil, dbResponse.Error
	}
	return &product, &sizes, nil
}

//#region CATEGORY

// GetAllCategories retrieves all categories from the database
func (r *RepositoryService) GetAllCategories(offset int, limit int, searchValue string) (*[]model.Category, error) {
	var categories []model.Category
	query := r.repo.DB.Table("categories c").
		Select("c.id, c.name, c.description")

	// Add search filter if searchValue is provided
	if searchValue != "" {
		searchPattern := "%" + searchValue + "%"
		query = query.Where("c.name ILIKE ?", searchPattern)
	}
	// Adding offset and limit for pagination
	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	// Perform query search
	dbResponse := query.Scan(&categories)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve categories from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &categories, nil
}

// UpdateCategory updates a category in the database by its ID
func (r *RepositoryService) UpdateCategory(updatedCategory *model.Category) error {
	dbResponse := r.repo.DB.Save(updatedCategory)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update category in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// InsertCategory inserts a new category into the database
func (r *RepositoryService) InsertCategory(newCategory *model.Category) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt").Create(&newCategory)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert category into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newCategory.ID, nil
}

// CheckCategoryExists checks if a category already exists in the database
func (r *RepositoryService) CheckCategoryExists(category *model.Category) (bool, error) {
	var count int64
	dbResponse := r.repo.DB.Model(&model.Category{}).Where("name = ?", category.Name).Count(&count)
	if dbResponse.Error != nil {
		log.Warnf("Failed to check if category exists in the database: %v", dbResponse.Error)
		return false, dbResponse.Error
	}
	return count > 0, nil
}

// GetTopFiveCategoriesSold retrieves the top five categories sold from the database
func (r *RepositoryService) GetTopFiveCategoriesSold() (*[]dto.TopCategoryDTO, error) {
	query := `
	SELECT c.name, SUM(oi.quantity) as total_sales
	FROM product_sizes s
	INNER JOIN products p ON p.id = s.product_id
	INNER JOIN categories c ON p.category_id = c.id
	INNER JOIN order_items oi ON oi.product_id = s.id
	INNER JOIN orders o ON o.id = oi.order_id
	WHERE o.status != 'cancelled'
	GROUP BY c.id
	ORDER BY SUM(oi.quantity) DESC
	LIMIT 5;`
	var categories []dto.TopCategoryDTO
	dbResponse := r.repo.DB.Raw(query).Scan(&categories)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve top five categories sold from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &categories, nil
}

// #region ORDER
// GetOrderList retrieves a list of orders from the database
func (r *RepositoryService) GetAdminOrderList(offset int, limit int, searchValue string, orderStatus string, paymentType string) (*[]model.OrderListItem, error) {
	var orders []model.OrderListItem
	query := r.repo.DB.Table("orders o").
		Select("o.id, o.order_date, u.name, o.payment_method, o.total_amount, o.status, o.delivery_date, o.tracking_id, o.shipping_company").
		Joins("INNER JOIN users u ON o.user_id = u.id")

	// Add search filter if searchValue is provided
	if searchValue != "" {
		searchPattern := "%" + searchValue + "%"
		query = query.Where("u.name ILIKE ?", searchPattern)
	}
	// Add status filter if orderStatus is provided
	if orderStatus != "" {
		query = query.Where("o.status = ?", orderStatus)
	}
	// Add payment filter if paymentType is provided
	if paymentType != "" {
		query = query.Where("o.payment_method = ?", paymentType)
	}
	// Adding offset and limit for pagination
	if offset >= 0 {
		query = query.Offset(offset)
	}
	if limit > 0 {
		query = query.Limit(limit)
	}

	dbResponse := query.Scan(&orders)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve orders from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &orders, nil
}

// GetOrderDetails retrieves the details of an order from the database
func (r *RepositoryService) GetOrderDetails(orderID uint64) (*dto.OrderDetailsDTO, error) {
	var order dto.OrderDetailsDTO
	// Fetch Order Details along with User and Address
	dbResponse := r.repo.DB.Raw(`
		SELECT 
			o.id, o.order_date, o.status, o.total_amount, o.delivery_date, o.payment_method, 
			o.tracking_id, o.shipping_company, o.delivery_cost, u.name, u.email, 
			u.phone, a.address_line1 as address_line1, 
			a.address_line2 as address_line2, a.city as city, 
			a.state as state, a.postal_code as postal_code
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id
		INNER JOIN addresses a ON o.address_id = a.id
		WHERE o.id = ?`, orderID).Scan(&order)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve order details from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}

	var items []dto.OrderItemDTO
	dbResponse = r.repo.DB.Raw(`
		SELECT 
			oi.id, p.images, p.name, s.size, oi.quantity, oi.price, oi.discount
		FROM order_items oi 
		INNER JOIN product_sizes s ON oi.product_id = s.id
		INNER JOIN products p ON p.id = s.product_id
		WHERE oi.order_id = ?`, orderID).Scan(&items)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve order items from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}

	order.Items = items
	return &order, nil
}

// PatchOrderStatus updates the status of an order in the database
func (r *RepositoryService) PatchOrderStatus(orderID uint64, status string) error {
	dbResponse := r.repo.DB.Model(&model.Order{}).Where("id = ?", orderID).Update("status", status)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update order status in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// #region ADMIN DASHBOARD

// GetAdminCardsData retrieves the data for the admin dashboard cards (total pending orders, total sales, total expenses, total customers)
func (r *RepositoryService) GetAdminCardsData() (*dto.AdminCardsDTO, error) {
	var cards dto.AdminCardsDTO
	// Fetch total pending orders
	dbResponse := r.repo.DB.Raw(`SELECT * FROM total_pending_orders`).Scan(&cards.TotalPendingOrders)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve total pending orders from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	// Fetch total sales
	dbResponse = r.repo.DB.Raw(`SELECT * FROM total_sales`).Scan(&cards.TotalSales)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve total sales from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	// Fetch total expenses
	dbResponse = r.repo.DB.Raw(`SELECT * FROM total_expenses`).Scan(&cards.TotalExpenses)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve total expenses from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	// Fetch total customers
	dbResponse = r.repo.DB.Raw(`SELECT * FROM total_customers`).Scan(&cards.TotalCustomers)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve total customers from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &cards, nil
}

// GetMonthlySalesData retrieves the monthly sales data for the admin dashboard
func (r *RepositoryService) GetMonthlySalesData() (*[]dto.MonthlySalesDTO, error) {
	var sales []dto.MonthlySalesDTO
	// query for retrieving monthly sales data
	query := r.repo.DB.Table("orders o").
		Select("DATE_TRUNC('month', o.order_date) as month, SUM(o.total_amount) as total_sales").
		Where("o.status = 'delivered'").
		Group("month").
		Order("month DESC").
		Limit(12)

	dbResponse := query.Scan(&sales)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve monthly sales data from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &sales, nil
}

// #region CARTS
// InsertCart inserts a new cart into the database
func (r *RepositoryService) InsertCart(newCart *model.Cart) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newCart)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert cart into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newCart.ID, nil
}

// GetCartByUserId retrieves a cart from the database by user ID
func (r *RepositoryService) GetCartByUserId(userID uint64) (*model.Cart, error) {
	var cart model.Cart
	dbResponse := r.repo.DB.Where("user_id = ?", userID).First(&cart)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve cart from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &cart, nil
}

// InsertCartProduct inserts a new cart product into the database
func (r *RepositoryService) InsertCartProduct(newCartProduct *model.CartProduct) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newCartProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert cart product into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newCartProduct.ID, nil
}

// UpdateCartProduct updates a cart product in the database
func (r *RepositoryService) UpdateCartProduct(updatedCartProduct *model.CartProduct) error {
	dbResponse := r.repo.DB.Save(updatedCartProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update cart product in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// DeleteCartProduct deletes a cart product from the database
func (r *RepositoryService) DeleteCartProduct(cartProductID uint64) error {
	dbResponse := r.repo.DB.Delete(&model.CartProduct{}, cartProductID)
	if dbResponse.Error != nil {
		log.Warnf("Failed to delete cart product from the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// GetCartProductByCartProductId retrieves a cart product from the database by cart and product ID
func (r *RepositoryService) GetCartProductByCartProductId(cartID uint64, productID uint64) (*model.CartProduct, error) {
	var cartProduct model.CartProduct
	dbResponse := r.repo.DB.Where("cart_id = ? AND product_id = ?", cartID, productID).First(&cartProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve cart product from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &cartProduct, nil
}

// UpdateProductCartQuantity updates the quantity of a product in the cart
func (r *RepositoryService) UpdateProductCartQuantity(cartID uint64, productID uint64, quantity int) error {
	dbResponse := r.repo.DB.Model(&model.CartProduct{}).
		Where("cart_id = ? AND product_id = ?", cartID, productID).
		Update("quantity", quantity)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update product cart quantity in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// #region WISHLIST
// InsertWishlist inserts a new wishlist into the database
func (r *RepositoryService) InsertWishlist(newWishlist *model.Wishlist) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newWishlist)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert wishlist into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newWishlist.ID, nil
}

// GetWishlistByUserId retrieves a wishlist from the database by user ID
func (r *RepositoryService) GetWishlistByUserId(userID uint64) (*model.Wishlist, error) {
	var wishlist model.Wishlist
	dbResponse := r.repo.DB.Where("user_id = ?", userID).First(&wishlist)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve wishlist from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &wishlist, nil
}

// InsertWishlistProduct inserts a new wishlist product into the database
func (r *RepositoryService) InsertWishlistProduct(newWishlistProduct *model.WishlistProduct) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newWishlistProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert wishlist product into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newWishlistProduct.ID, nil
}

// GetWishlistProductsByWishlistId retrieves all wishlist products from the database by wishlist ID
func (r *RepositoryService) GetWishlistProductsById(wishlistID uint64, userID uint64) (*[]model.WishlistProduct, error) {
	var wishlistProducts []model.WishlistProduct
	dbResponse := r.repo.DB.Where("wishlist_id = ?", wishlistID).Where("user_id = ?", userID).Find(&wishlistProducts)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve wishlist products from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &wishlistProducts, nil
}

// UpdateWishlistProduct updates a wishlist product in the database
func (r *RepositoryService) UpdateWishlistProduct(updatedWishlistProduct *model.WishlistProduct) error {
	dbResponse := r.repo.DB.Save(updatedWishlistProduct)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update wishlist product in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// DeleteWishlistProduct deletes a wishlist product from the database
func (r *RepositoryService) DeleteWishlistProduct(wishlistProductID uint64) error {
	dbResponse := r.repo.DB.Delete(&model.WishlistProduct{}, wishlistProductID)
	if dbResponse.Error != nil {
		log.Warnf("Failed to delete wishlist product from the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// #region ADDRESS
// InsertAddress inserts a new address into the database
func (r *RepositoryService) InsertAddress(newAddress *model.Address) (uint64, error) {
	dbResponse := r.repo.DB.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newAddress)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert address into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newAddress.ID, nil
}

// GetAddressesListByUserId retrieves a list of addresses from the database by user ID
func (r *RepositoryService) GetAddressesListByUserId(userID uint64) (*[]model.Address, error) {
	var addresses []model.Address
	dbResponse := r.repo.DB.Where("user_id = ?", userID).Where("deleted_at IS NULL").Omit("CreatedAt", "UpdatedAt", "DeletedAt", "Addressee", "Phone", "AddresseeEmail").Find(&addresses)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve addresses from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &addresses, nil
}
