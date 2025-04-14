package services

import (
	"fmt"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/repositories"
	"github.com/SimpleBookRental/backend/pkg/utils"
)

// TokenService handles business logic for tokens
type TokenService struct {
	tokenRepo repositories.TokenRepositoryInterface
	userRepo  repositories.UserRepositoryInterface
}

// Ensure TokenService implements TokenServiceInterface
var _ TokenServiceInterface = (*TokenService)(nil)

// NewTokenService creates a new token service
func NewTokenService(tokenRepo repositories.TokenRepositoryInterface, userRepo repositories.UserRepositoryInterface) *TokenService {
	return &TokenService{tokenRepo: tokenRepo, userRepo: userRepo}
}

// RefreshToken refreshes an access token using a refresh token
func (s *TokenService) RefreshToken(request *models.RefreshTokenRequest) (*models.RefreshTokenResponse, error) {
	// Find token in database
	issuedToken, err := s.tokenRepo.FindTokenByValue(request.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("error finding token: %w", err)
	}
	if issuedToken == nil {
		return nil, fmt.Errorf("token not found")
	}

	// Check if token is revoked or expired
	if issuedToken.IsRevoked {
		return nil, fmt.Errorf("token has been revoked")
	}
	if issuedToken.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token has expired")
	}

	// Check if token is a refresh token
	if issuedToken.TokenType != string(models.RefreshToken) {
		return nil, fmt.Errorf("invalid token type")
	}

	// Get user from token
	user, err := s.userRepo.FindByID(issuedToken.UserID)
	if err != nil {
		return nil, fmt.Errorf("error finding user: %w", err)
	}
	if user == nil {
		return nil, fmt.Errorf("user not found")
	}

	// Generate new token pair
	tokenPair, err := utils.GenerateTokenPair(user.ID, user.Email)
	if err != nil {
		return nil, fmt.Errorf("error generating tokens: %w", err)
	}

	// Save new tokens in database
	accessExpiration := time.Now().Add(time.Hour * 24)       // 24 hours
	refreshExpiration := time.Now().Add(time.Hour * 24 * 30) // 30 days

	// Save access token
	accessToken := &models.IssuedToken{
		UserID:    user.ID,
		Token:     tokenPair.AccessToken,
		TokenType: string(models.AccessToken),
		ExpiresAt: accessExpiration,
	}
	err = s.tokenRepo.CreateToken(accessToken)
	if err != nil {
		return nil, fmt.Errorf("error saving access token: %w", err)
	}

	// Save refresh token
	refreshToken := &models.IssuedToken{
		UserID:    user.ID,
		Token:     tokenPair.RefreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: refreshExpiration,
	}
	err = s.tokenRepo.CreateToken(refreshToken)
	if err != nil {
		return nil, fmt.Errorf("error saving refresh token: %w", err)
	}

	// Revoke old refresh token
	err = s.tokenRepo.RevokeToken(issuedToken)
	if err != nil {
		return nil, fmt.Errorf("error revoking old token: %w", err)
	}

	// Create response
	response := &models.RefreshTokenResponse{
		AccessToken:  tokenPair.AccessToken,
		RefreshToken: tokenPair.RefreshToken,
		ExpiresAt:    accessExpiration.Unix(),
	}

	return response, nil
}

// Logout revokes a refresh token
func (s *TokenService) Logout(request *models.LogoutRequest) error {
	// Find token in database
	issuedToken, err := s.tokenRepo.FindTokenByValue(request.RefreshToken)
	if err != nil {
		return fmt.Errorf("error finding token: %w", err)
	}

	// If token not found, nothing to do
	if issuedToken == nil {
		return nil
	}

	// Check if token is a refresh token
	if issuedToken.TokenType != string(models.RefreshToken) {
		return fmt.Errorf("invalid token type")
	}

	// Revoke all tokens for the user
	err = s.tokenRepo.RevokeAllUserTokens(issuedToken.UserID)
	if err != nil {
		return fmt.Errorf("error revoking user tokens: %w", err)
	}

	// Cleanup expired tokens
	_ = s.tokenRepo.CleanupExpiredTokens()

	return nil
}
