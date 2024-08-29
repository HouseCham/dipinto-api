package services

import (
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
	dbResponse := r.repo.DB.Find(&users)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve users from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return users, nil
}

// GetUser retrieves a user from the database
func (r *RepositoryService) GetUserById(userID uint64) (*model.User, error) {
	var user model.User
	dbResponse := r.repo.DB.First(&user, userID)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve user from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &user, nil
}

// GetUserByEmail retrieves a user from the database by email
func (r *RepositoryService) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	dbResponse := r.repo.DB.Where("email = ?", email).First(&user)
	if dbResponse.Error != nil {
		log.Warnf("Failed to retrieve user from the database: %v", dbResponse.Error)
		return nil, dbResponse.Error
	}
	return &user, nil
}

// UpdateUser updates a user in the database
func (r *RepositoryService) UpdateUser(updatedUser *model.User) error {
	dbResponse := r.repo.DB.Save(updatedUser)
	if dbResponse.Error != nil {
		log.Warnf("Failed to update user in the database: %v", dbResponse.Error)
		return dbResponse.Error
	}
	return nil
}

// ValidateEmailUsed checks if an email is already used within the database and returns true if it is
func (r *RepositoryService) ValidateEmailUsed(email string) bool {
	var user model.User
	dbResponse := r.repo.DB.Where("email = ?", email).First(&user)
	return dbResponse.Error == nil
}
