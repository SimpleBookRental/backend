package controllers

import (
	"github.com/SimpleBookRental/backend/internal/models"
	"github.com/SimpleBookRental/backend/internal/services"
	"github.com/SimpleBookRental/backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

// TokenController handles HTTP requests for tokens
type TokenController struct {
	tokenService services.TokenServiceInterface
}

// NewTokenController creates a new token controller
func NewTokenController(tokenService services.TokenServiceInterface) *TokenController {
	return &TokenController{tokenService: tokenService}
}

// RefreshToken godoc
// @Summary      Refresh token
// @Description  Refresh JWT token
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      models.RefreshTokenRequest  true  "Refresh token payload"
// @Success      200     {object}  models.RefreshTokenResponse
// @Failure      400     {object}  models.ErrorResponse
// @Router       /api/v1/refresh-token [post]
func (c *TokenController) RefreshToken(ctx *gin.Context) {
	var request models.RefreshTokenRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	response, err := c.tokenService.RefreshToken(&request)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, response)
}

// Logout godoc
// @Summary      Logout
// @Description  User logout
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body      models.LogoutRequest  true  "Logout payload"
// @Success      200     {object}  map[string]bool
// @Failure      400     {object}  models.ErrorResponse
// @Router       /api/v1/logout [post]
func (c *TokenController) Logout(ctx *gin.Context) {
	var request models.LogoutRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	err := c.tokenService.Logout(&request)
	if err != nil {
		utils.BadRequest(ctx, err.Error())
		return
	}

	utils.OK(ctx, gin.H{"logged_out": true})
}
