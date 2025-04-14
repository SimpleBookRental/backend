package services

import (
	"fmt"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/pkg/utils"
)

// TokenService handles business logic for tokens
type TokenService struct{}

// NewTokenService creates a new token service
func NewTokenService() *TokenService {
	return &TokenService{}
}

// RefreshToken refreshes an access token using a refresh token
func (s *TokenService) RefreshToken(request *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	// Refresh tokens
	tokenPair, err := utils.RefreshTokens(request.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error refreshing tokens: %w", err)
	}

	// Create response
	response := &models.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 24).Unix(), // 24 hours
	}

	return response, nil
}
