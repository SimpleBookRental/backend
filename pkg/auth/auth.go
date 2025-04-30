package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/SimpleBookRental/backend/internal/domain"
	"github.com/SimpleBookRental/backend/pkg/config"
	"github.com/golang-jwt/jwt/v5"
)

// TokenType defines the type of token
type TokenType string

const (
	// AccessToken represents an access token
	AccessToken TokenType = "access"
	// RefreshToken represents a refresh token
	RefreshToken TokenType = "refresh"
)

// Claims represents the JWT claims
type Claims struct {
	UserID   int64     `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// JWTService provides JWT token generation and validation
type JWTService struct {
	config *config.JWTConfig
}

// NewJWTService creates a new JWTService
func NewJWTService(config *config.JWTConfig) *JWTService {
	return &JWTService{
		config: config,
	}
}

// GenerateAccessToken generates a new access token
func (s *JWTService) GenerateAccessToken(user *domain.User) (string, error) {
	return s.generateToken(user, AccessToken, s.config.ExpirationHours)
}

// GenerateRefreshToken generates a new refresh token
func (s *JWTService) GenerateRefreshToken(user *domain.User) (string, error) {
	return s.generateToken(user, RefreshToken, s.config.RefreshExpiration)
}

// ValidateToken validates a token and returns the claims
func (s *JWTService) ValidateToken(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// generateToken generates a new token
func (s *JWTService) generateToken(user *domain.User, tokenType TokenType, expiration time.Duration) (string, error) {
	claims := &Claims{
		UserID:   user.ID,
		Username: user.Username,
		Role:     string(user.Role),
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			Issuer:    "book-rental-api",
			Subject:   fmt.Sprintf("%d", user.ID),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.Secret))
}

// GetUserIDFromToken extracts the user ID from a token
func (s *JWTService) GetUserIDFromToken(tokenString string) (int64, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return 0, err
	}
	return claims.UserID, nil
}

// GetUserRoleFromToken extracts the user role from a token
func (s *JWTService) GetUserRoleFromToken(tokenString string) (domain.UserRole, error) {
	claims, err := s.ValidateToken(tokenString)
	if err != nil {
		return "", err
	}
	return domain.UserRole(claims.Role), nil
}

// IsAdmin checks if the user has admin role
func IsAdmin(role domain.UserRole) bool {
	return role == domain.RoleAdmin
}

// IsLibrarian checks if the user has librarian role
func IsLibrarian(role domain.UserRole) bool {
	return role == domain.RoleLibrarian || role == domain.RoleAdmin
}

// IsMember checks if the user has member role
func IsMember(role domain.UserRole) bool {
	return role == domain.RoleMember || role == domain.RoleLibrarian || role == domain.RoleAdmin
}
