package services

import (
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/auth"
	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/security"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
)

// AuthService contains the business logic related to authentication.
type AuthService struct {
	auth *auth.AuthService
	security *security.SecurityService
}
// NewAuthService creates a new AuthService
func NewAuthService(auth *auth.AuthService) *AuthService {
	return &AuthService{
		auth: auth,
		security: &security.SecurityService{},
	}
}

// HashPassword hashes a password using bcrypt
func (a *AuthService) HashPassword(password string) (string, error) {
	return a.security.HashPassword(password)
}
// ValidatePassword validates a password against a hashed password. Returns nil if the password is correct.
func (a *AuthService) ValidatePassword(password, hashedPassword string) error {
	return a.security.CheckPassword(password, hashedPassword)
}

// GenerateToken generates a JWT token for the given username and user ID
func (a *AuthService) GenerateToken(user *model.User) (string, error) {
	return a.auth.CreateToken(user)
}