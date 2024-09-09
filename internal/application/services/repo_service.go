package services

import (
	"encoding/json"

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
	dbResponse := r.repo.DB.Omit("ID","CreatedAt","UpdatedAt","DeletedAt").Create(&newUser)
	if dbResponse.Error != nil {
		log.Warnf("Failed to insert user into the database: %v", dbResponse.Error)
		return 0, dbResponse.Error
	}
	return newUser.ID, nil
}

// GetAllUsers retrieves all users from the database
func (r *RepositoryService) GetAllUsers() ([]model.User, error) {
	var users []model.User
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Omit("CreatedAt","UpdatedAt","DeletedAt").Find(&users)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve users from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return users, nil
}

// GetAllCustomers retrieves all customers from the database
func (r *RepositoryService) GetAllCustomers() ([]model.User, error) {
	var customers []model.User
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Where("role='customer'").Omit("UpdatedAt","DeletedAt").Find(&customers)
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
func (r *RepositoryService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	dbResponse := r.repo.DB.Where("deleted_at IS NULL").Where("email = ?", email).Omit("CreatedAt","UpdatedAt","DeletedAt").First(&user)
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
// GetAllProducts retrieves all products from the database
func (r *RepositoryService) GetAllProductsCatalogue() (*[]model.CatalogueProduct, error) {
	var products []model.CatalogueProduct
	query := `
		SELECT p.id, p.slug, p.name, p.images, s.price, s.discount
		FROM
			products p
			INNER JOIN product_sizes s ON p.id = s.product_id
			INNER JOIN (
				SELECT product_id, MIN(price) AS min_price
				FROM product_sizes
				WHERE
					is_available = true
				GROUP BY
					product_id
			) min_sizes ON s.product_id = min_sizes.product_id
			AND s.price = min_sizes.min_price
		WHERE
			s.is_available = true;
	`
	dbResponse := r.repo.DB.Raw(query).Scan(&products)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve products from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &products, nil
}
// GetAllProductsAdmin retrieves all products from the database for admin purposes
func (r *RepositoryService) GetAllProductsAdmin() (*[]model.AdminProduct, error) {
	// get products from the database
	var products []model.Product
	query := `
		SELECT p.id, p.slug, p.name, c.name as category, p.description, p.images
		FROM products p
			INNER JOIN categories c ON p.category_id = c.id
		WHERE
			p.deleted_at IS NULL;
	`
	dbResponse := r.repo.DB.Raw(query).Scan(&products)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve products from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	// get product sizes from the database
	var sizes []model.ProductSize
	for i := range products {
		var productSizes []model.ProductSize
		dbResponse = r.repo.DB.Where("product_id = ?", products[i].ID).Where("deleted_at IS NULL").Omit("CreatedAt", "UpdatedAt", "DeletedAt").Find(&productSizes)
		if dbResponse.Error != nil {
			log.Warnf("Failed to retrieve product sizes from the database: %v", dbResponse.Error)
			return nil, dbResponse.Error
		}
		sizes = append(sizes, productSizes...)
	}

	// create a slice of AdminProduct
	var adminProducts []model.AdminProduct
	for _, product := range products {
		// get the product sizes for the current product filtering by product ID on sizes slice
		var productSizes []model.ProductSize
		for _, size := range sizes {
			if size.ProductID == product.ID {
				productSizes = append(productSizes, size)
			}
		}
		adminProducts = append(adminProducts, *utils.ParseProductToAdminProduct(&product, &productSizes))
	}

	return &adminProducts, nil
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
func (r *RepositoryService) GetAllCategories() (*[]model.Category, error) {
	var categories []model.Category
	dbResponse := r.repo.DB.Omit("CreatedAt", "UpdatedAt").Find(&categories)
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

//#region ORDER
// GetOrderList retrieves a list of orders from the database
func (r *RepositoryService) GetAdminOrderList() (*[]model.OrderListItem, error) {
	var orders []model.OrderListItem
	query := `
		SELECT o.id, o.order_date, u.name, o.payment_method, o.total_amount, o.status, o.delivery_date, o.tracking_id, o.shipping_company
		FROM orders o
		INNER JOIN users u ON o.user_id = u.id;
	`
	dbResponse := r.repo.DB.Raw(query).Scan(&orders)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve orders from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &orders, nil
}