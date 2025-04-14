package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// TokenController handles HTTP requests for tokens
type TokenController struct {
	tokenService *services.TokenService
}

// NewTokenController creates a new token controller
func NewTokenController(tokenService *services.TokenService) *TokenController {
	return &TokenController{tokenService: tokenService}
}

// RefreshToken handles refreshing tokens
func (c *TokenController) RefreshToken(ctx *gin.Context) {
	var request models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, "Invalid request body", err.Error())
		return
	}

	response, err := c.tokenService.RefreshToken(&request)
	if err != nil {
		utils.BadRequest(ctx, "Failed to refresh token", err.Error())
		return
	}

	utils.OK(ctx, "Token refreshed successfully", response)
}
