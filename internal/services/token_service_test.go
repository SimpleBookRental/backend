// Unit tests for TokenService using gomock and testify.
package services

import (
	"errors"
	"testing"
	"time"

	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestTokenService_RefreshToken_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	refreshToken := "refresh-token"
	userID := "user-1"
	issuedToken := &models.IssuedToken{
		UserID:    userID,
		Token:     refreshToken,
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	user := &models.User{ID: userID, Email: "test@example.com"}
	mockTokenRepo.EXPECT().FindTokenByValue(refreshToken).Return(issuedToken, nil)
	mockUserRepo.EXPECT().FindByID(userID).Return(user, nil)
	mockTokenRepo.EXPECT().CreateToken(gomock.Any()).Return(nil).AnyTimes()
	mockTokenRepo.EXPECT().RevokeToken(issuedToken).Return(nil)

	req := &models.RefreshTokenRequest{RefreshToken: refreshToken}
	resp, err := service.RefreshToken(req)
	assert.NoError(t, err)
	assert.NotEmpty(t, resp.AccessToken)
	assert.NotEmpty(t, resp.RefreshToken)
}

func TestTokenService_RefreshToken_TokenNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	mockTokenRepo.EXPECT().FindTokenByValue("notfound").Return(nil, nil)

	req := &models.RefreshTokenRequest{RefreshToken: "notfound"}
	resp, err := service.RefreshToken(req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "token not found")
}

func TestTokenService_RefreshToken_Revoked(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
		IsRevoked: true,
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)

	req := &models.RefreshTokenRequest{RefreshToken: "refresh-token"}
	resp, err := service.RefreshToken(req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "token has been revoked")
}

func TestTokenService_RefreshToken_Expired(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(-time.Hour),
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)

	req := &models.RefreshTokenRequest{RefreshToken: "refresh-token"}
	resp, err := service.RefreshToken(req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "token has expired")
}

func TestTokenService_RefreshToken_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.AccessToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)

	req := &models.RefreshTokenRequest{RefreshToken: "refresh-token"}
	resp, err := service.RefreshToken(req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "invalid token type")
}

func TestTokenService_RefreshToken_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)
	mockUserRepo.EXPECT().FindByID("user-1").Return(nil, nil)

	req := &models.RefreshTokenRequest{RefreshToken: "refresh-token"}
	resp, err := service.RefreshToken(req)
	assert.Nil(t, resp)
	assert.ErrorContains(t, err, "user not found")
}

func TestTokenService_Logout_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)
	mockTokenRepo.EXPECT().RevokeAllUserTokens("user-1").Return(nil)
	mockTokenRepo.EXPECT().CleanupExpiredTokens().Return(nil)

	req := &models.LogoutRequest{RefreshToken: "refresh-token"}
	err := service.Logout(req)
	assert.NoError(t, err)
}

func TestTokenService_Logout_TokenNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	mockTokenRepo.EXPECT().FindTokenByValue("notfound").Return(nil, nil)

	req := &models.LogoutRequest{RefreshToken: "notfound"}
	err := service.Logout(req)
	assert.NoError(t, err)
}

func TestTokenService_Logout_InvalidType(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.AccessToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)

	req := &models.LogoutRequest{RefreshToken: "refresh-token"}
	err := service.Logout(req)
	assert.ErrorContains(t, err, "invalid token type")
}

func TestTokenService_Logout_RepoError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	mockTokenRepo := mocks.NewMockTokenRepositoryInterface(ctrl)
	mockUserRepo := mocks.NewMockUserRepositoryInterface(ctrl)
	service := NewTokenService(mockTokenRepo, mockUserRepo)

	issuedToken := &models.IssuedToken{
		UserID:    "user-1",
		Token:     "refresh-token",
		TokenType: string(models.RefreshToken),
		ExpiresAt: time.Now().Add(time.Hour),
	}
	mockTokenRepo.EXPECT().FindTokenByValue("refresh-token").Return(issuedToken, nil)
	mockTokenRepo.EXPECT().RevokeAllUserTokens("user-1").Return(errors.New("revoke error"))

	req := &models.LogoutRequest{RefreshToken: "refresh-token"}
	err := service.Logout(req)
	assert.ErrorContains(t, err, "error revoking user tokens")
}
