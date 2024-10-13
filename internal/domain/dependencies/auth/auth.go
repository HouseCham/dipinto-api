package auth

import "fmt"

import (
	"time"

	"github.com/HouseCham/dipinto-api/internal/domain/dependencies/middleware"
	"github.com/HouseCham/dipinto-api/internal/domain/model"
	"github.com/HouseCham/dipinto-api/internal/infrastructure/config"
	"github.com/dgrijalva/jwt-go"
)

type AuthService struct {
	jwtKey []byte
}

// AuthService contains the business logic related to authentication.
func SetUpAuthService(config *config.Config) *AuthService {
	authService := AuthService{
		jwtKey: []byte(config.JWT.SecretKey),
	}
	return &authService
}

// CreateToken creates a JWT token with the given username
func (auth *AuthService) CreateToken(user *model.User, remember bool) (string, error) {
	// Set the expiration time of the token
	expirationTime := time.Now().Add(8 * time.Hour) // Token expires in 8 hours

	// Create the claims
	claims := &middleware.Claims{
		ID:       fmt.Sprint(user.ID),
		Username: user.Name,
		Role:     user.Role,
		Remember: remember,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "dipinto-api",
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with the secret key
	signedToken, err := token.SignedString(auth.jwtKey)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}
