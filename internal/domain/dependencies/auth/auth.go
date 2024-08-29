package auth

import (
	"time"

	"github.com/HouseCham/dipinto-api/internal/infrastructure/config"
	"github.com/dgrijalva/jwt-go"
)

// Claims structure
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

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
func (auth *AuthService) CreateToken(username string, userID int64) (string, error) {
	// Set the expiration time of the token
	expirationTime := time.Now().Add(24 * time.Hour) // Token expires in 24 hours

	// Create the claims
	claims := &Claims{
		Username: username,
		StandardClaims: jwt.StandardClaims{
			Id:        string(rune(userID)),
			ExpiresAt: expirationTime.Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "user-service",
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