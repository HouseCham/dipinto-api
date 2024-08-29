package services

import "github.com/HouseCham/dipinto-api/internal/domain/dependencies/db"

// RepositoryService contains the business logic related to the repository.
type RepositoryService struct {
	db *db.Database
}
// NewRepositoryService creates a new RepositoryService
func NewRepositoryService(db *db.Database) *RepositoryService {
	return &RepositoryService{
		db: db,
	}
}