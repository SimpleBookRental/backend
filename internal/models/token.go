package models

// RefreshTokenRequest is the request body for refreshing tokens
// @Description Refresh token request
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required" example:"refresh-token"` // Refresh token
}

// RefreshTokenResponse is the response body for refreshing tokens
// @Description Refresh token response
type RefreshTokenResponse struct {
	AccessToken  string `json:"access_token" example:"jwt-token"`         // New access token
	RefreshToken string `json:"refresh_token" example:"refresh-token"`    // New refresh token
	ExpiresAt    int64  `json:"expires_at" example:"1713345600"`          // Expiry timestamp
}
