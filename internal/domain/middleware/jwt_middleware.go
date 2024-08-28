package middleware

import (
	"errors"
	"time"

	"github.com/HouseCham/dipinto-api/internal/infrastructure/config"
	"github.com/dgrijalva/jwt-go"
)

// MiddlewareService contains the business logic related to middleware.
type MiddlewareService struct {
	jwtKey []byte
}
// Claims represents the JWT claims
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}
// SetupMiddlewareService creates a new MiddlewareService
func SetupMiddlewareService(config *config.Config) *MiddlewareService {
	middlewareService := MiddlewareService{
		jwtKey: []byte(config.JWT.SecretKey),
	}
	return &middlewareService
}

// ValidateToken validates the JWT token and checks the claim ID is greater than 0
func (m *MiddlewareService) ValidateToken(tokenStr string) (*Claims, error) {
    // Parse the token
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return m.jwtKey, nil
    })

    if err != nil {
        if err == jwt.ErrSignatureInvalid {
            return nil, errors.New("invalid token signature")
        }
        return nil, errors.New("invalid token")
    }

    // Check if the token is valid
    if !token.Valid {
        return nil, errors.New("invalid token")
    }

    // Extract the claims from the token
    claims, ok := token.Claims.(*Claims)
    if !ok {
        return nil, errors.New("could not parse claims")
    }

    // Check if the token has expired
    if claims.ExpiresAt < time.Now().Unix() {
        return nil, errors.New("token has expired")
    }

    // Check if the claim ID is greater than 0
    if claims.Id == "" || claims.Id == "0" {
        return nil, errors.New("invalid claim ID")
    }

    return claims, nil
}
