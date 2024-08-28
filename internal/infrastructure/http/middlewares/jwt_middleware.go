package middlewares

import (
    "errors"
    "github.com/gofiber/fiber/v3"
    "github.com/dgrijalva/jwt-go"
    "time"
)

// Claims represents the JWT claims
type Claims struct {
    Username string `json:"username"`
    jwt.StandardClaims
}

// JWTMiddleware is a Fiber middleware function that validates JWT tokens
func JWTMiddleware(jwtKey []byte) fiber.Handler {
    return func(c fiber.Ctx) error {
        // Get the JWT token from the Authorization header
        tokenStr := c.Get("Authorization")
        if tokenStr == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing token"})
        }

        // Extract the token from the "Bearer " prefix
        if len(tokenStr) > 7 && tokenStr[:7] == "Bearer " {
            tokenStr = tokenStr[7:]
        } else {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token format"})
        }

        // Parse and validate the token
        claims, err := validateToken(tokenStr, jwtKey)
        if err != nil {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
        }

        // Store the claims in the request context
        c.Locals("claims", claims)

        // Proceed to the next middleware/handler
        return c.Next()
    }
}

// validateToken parses the JWT token and checks that the claim ID is greater than 0
func validateToken(tokenStr string, jwtKey []byte) (*Claims, error) {
    // Parse the token
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return jwtKey, nil
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