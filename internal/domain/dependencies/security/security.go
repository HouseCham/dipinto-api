package security

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)
// SecurityService is a service that provides security-related functionalities
type SecurityService struct {}

// HashPassword hashes the provided password
func (s *SecurityService) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %s", err)
	}
	return string(hashedPassword), nil
}
// CheckPassword checks if the provided password is correct or not
func (s *SecurityService) CheckPassword(password string, hashedPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}