package services

import (
	"encoding/json"

	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/db"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
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
	dbResponse := tx.Omit("ID", "CreatedAt", "UpdatedAt", "DeletedAt").Create(&newProduct)
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