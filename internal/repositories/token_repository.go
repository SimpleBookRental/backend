package repositories

import (
	"time"
	"fmt"

	"github.com/SimpleBookRental/backend/internal/models"
	"gorm.io/gorm"
)

// TokenRepository handles database operations for tokens
type TokenRepository struct {
	db *gorm.DB
}

// GetDB returns the database connection
func (r *TokenRepository) GetDB() interface{} {
	return r.db
}


// WithTx returns a new TokenRepository with the given transaction
func (r *TokenRepository) WithTx(tx interface{}) (TokenRepositoryInterface, error) {
	db, ok := tx.(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}
	return &TokenRepository{db: db}, nil
}

// Ensure TokenRepository implements TokenRepositoryInterface
var _ TokenRepositoryInterface = (*TokenRepository)(nil)

// NewTokenRepository creates a new token repository
func NewTokenRepository(db *gorm.DB) *TokenRepository {
	return &TokenRepository{db: db}
}

// CreateToken creates a new token
func (r *TokenRepository) CreateToken(token *models.IssuedToken) error {
	return r.db.Create(token).Error
}

// FindTokenByValue finds a token by its value
func (r *TokenRepository) FindTokenByValue(tokenString string) (*models.IssuedToken, error) {
	var token models.IssuedToken
	err := r.db.Where("token = ?", tokenString).First(&token).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &token, nil
}

// FindActiveTokensByUserID finds all active tokens for a user
func (r *TokenRepository) FindActiveTokensByUserID(userID string) ([]models.IssuedToken, error) {
	var tokens []models.IssuedToken
	err := r.db.Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now()).Find(&tokens).Error
	return tokens, err
}

// RevokeToken revokes a token
func (r *TokenRepository) RevokeToken(token *models.IssuedToken) error {
	now := time.Now()
	token.IsRevoked = true
	token.RevokedAt = &now
	return r.db.Save(token).Error
}

// RevokeAllUserTokens revokes all tokens for a user
func (r *TokenRepository) RevokeAllUserTokens(userID string) error {
	now := time.Now()
	return r.db.Model(&models.IssuedToken{}).
		Where("user_id = ? AND is_revoked = ? AND expires_at > ?", userID, false, time.Now()).
		Updates(map[string]interface{}{
			"is_revoked": true,
			"revoked_at": now,
			"updated_at": now,
		}).Error
}

// CleanupExpiredTokens removes expired tokens
func (r *TokenRepository) CleanupExpiredTokens() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.IssuedToken{}).Error
}
