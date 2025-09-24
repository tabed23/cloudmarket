package jwt

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

var jwtSecret = []byte("cooki")

// JwtClaims represents the JWT claims for authentication
type JwtClaims struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Role  string `json:"role"` // Added role for authorization
	jwt.StandardClaims
}

// GenreateJwt generates a JWT token with claims including user ID, email, and role
func GenreateJwt(ctx context.Context, id, email, role string) (string, error) {
	claims := JwtClaims{
		ID:    id,
		Email: email,
		Role:  role, // Adding role to claims
		StandardClaims: jwt.StandardClaims{
			Issuer:    "cloudmarket",
			ExpiresAt: time.Now().Add(24 * time.Hour).Unix(), // Token expires in 24 hours
		},
	}

	// Generate the token with the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

// ValidateJwt validates the JWT token and returns the claims
func ValidateJwt(ctx context.Context, tokenString string) (*JwtClaims, error) {
	// Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &JwtClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	// Extract the claims from the token
	claims, ok := token.Claims.(*JwtClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check if the token has expired
	if claims.ExpiresAt < time.Now().Unix() {
		return nil, fmt.Errorf("token has expired")
	}

	return claims, nil
}
