package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// GetAccessSecret gets the JWT secret key from environment variables
func GetAccessSecret() []byte {
	secret := os.Getenv("JWT_ACCESS_SECRET")
	if secret == "" {
		secret = "access_token_secret_key" // Default value
	}
	return []byte(secret)
}

// getRefreshSecret gets the JWT refresh secret key from environment variables
func getRefreshSecret() []byte {
	secret := os.Getenv("JWT_REFRESH_SECRET")
	if secret == "" {
		secret = "refresh_token_secret_key" // Default value
	}
	return []byte(secret)
}

// getAccessExpiration gets the JWT access token expiration from environment variables
func getAccessExpiration() time.Duration {
	expStr := os.Getenv("JWT_ACCESS_EXPIRATION")
	if expStr == "" {
		return time.Hour * 24 // Default: 24 hours
	}

	exp, err := time.ParseDuration(expStr)
	if err != nil {
		return time.Hour * 24 // Default: 24 hours
	}

	return exp
}

// getRefreshExpiration gets the JWT refresh token expiration from environment variables
func getRefreshExpiration() time.Duration {
	expStr := os.Getenv("JWT_REFRESH_EXPIRATION")
	if expStr == "" {
		return time.Hour * 24 * 30 // Default: 30 days
	}

	exp, err := time.ParseDuration(expStr)
	if err != nil {
		return time.Hour * 24 * 30 // Default: 30 days
	}

	return exp
}

// Claims represents the JWT claims
type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// TokenPair represents an access and refresh token pair
type TokenPair struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// GenerateTokenPair generates a new access and refresh token pair
func GenerateTokenPair(userID, email string) (*TokenPair, error) {
	// Generate access token
	accessToken, err := generateToken(userID, email, GetAccessSecret(), getAccessExpiration())
	if err != nil {
		return nil, fmt.Errorf("failed to generate access token: %w", err)
	}

	// Generate refresh token
	refreshToken, err := generateToken(userID, email, getRefreshSecret(), getRefreshExpiration())
	if err != nil {
		return nil, fmt.Errorf("failed to generate refresh token: %w", err)
	}

	return &TokenPair{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}

// generateToken generates a new JWT token
func generateToken(userID, email string, secret []byte, expiration time.Duration) (string, error) {
	// Create claims
	claims := &Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign token
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a JWT token
func ValidateToken(tokenString string, secret []byte) (*Claims, error) {
	// Parse token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})

	if err != nil {
		return nil, err
	}

	// Validate token
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// RefreshTokens refreshes an access token using a refresh token
func RefreshTokens(refreshToken string) (*TokenPair, error) {
	// Validate refresh token
	claims, err := ValidateToken(refreshToken, getRefreshSecret())
	if err != nil {
		return nil, fmt.Errorf("invalid refresh token: %w", err)
	}

	// Generate new token pair
	return GenerateTokenPair(claims.UserID, claims.Email)
}
