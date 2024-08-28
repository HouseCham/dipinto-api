package services

import (
	"github.com/HouseCham/dipinto-api/internal/domain/auth"
	"github.com/HouseCham/dipinto-api/internal/domain/security"
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