// internal/handler/auth_handler.go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/service"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db"
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

// Login handles user login requests.
// @Summary      User login
// @Description  Authenticate a user and return a JWT token pair.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Login Credentials"
// @Success      200  {object}  GenericResponse{data=TokenResponse} "Login Successful"
// @Failure      400  {object}  GenericResponse "Invalid Request"
// @Failure      401  {object}  GenericResponse "Unauthorized"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid request body"})
		return
	}

	tokens, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "Login successful",
		Data:    TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken},
	})
}

// RefreshToken handles requests to refresh JWT tokens.
// @Summary      Refresh token
// @Description  Get a new pair of access and refresh tokens.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        token body RefreshTokenRequest true "Refresh Token"
// @Success      200  {object}  GenericResponse{data=TokenResponse} "Tokens refreshed"
// @Failure      401  {object}  GenericResponse "Unauthorized"
// @Router       /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid request body"})
		return
	}

	tokens, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, GenericResponse{Success: false, Message: err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "Tokens refreshed successfully",
		Data:    TokenResponse{AccessToken: tokens.AccessToken, RefreshToken: tokens.RefreshToken},
	})
}


// GetProfile handles requests to get the current user's profile.
// @Summary      Get user profile
// @Description  Get the profile of the currently authenticated user.
// @Tags         Authentication
// @Security     BearerAuth
// @Produce      json
// @Success      200  {object}  ProfileResponse "Profile Retrieved"
// @Failure      401  {object}  GenericResponse "Unauthorized"
// @Router       /auth/profile [get]
func (h *AuthHandler) GetProfile(c *gin.Context) {
	userCtx, exists := c.Get("user")
	if !exists {
		c.JSON(http.StatusUnauthorized, GenericResponse{Success: false, Message: "User not found in context"})
		return
	}

	user := userCtx.(*db.UserModel)
	
	profileData := ProfileData{
		ID:        int64(user.ID),
		Username:  user.Username,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.String(),
	}

	c.JSON(http.StatusOK, ProfileResponse{
		Success: true,
		Message: "Profile retrieved successfully",
		Data:    profileData,
	})
}

// ChangePassword handles requests to change the user's password.
// @Summary      Change user password
// @Description  Allows an authenticated user to change their password.
// @Tags         Authentication
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        passwords body ChangePasswordRequest true "Current and new passwords"
// @Success      200  {object}  GenericResponse "Password changed successfully"
// @Failure      400  {object}  GenericResponse "Invalid request body"
// @Failure      401  {object}  GenericResponse "Unauthorized or incorrect password"
// @Router       /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	userCtx, _ := c.Get("user")
	user := userCtx.(*db.UserModel)

	err := h.service.ChangePassword(int(user.ID), req.CurrentPassword, req.NewPassword)
	if err != nil {
		c.JSON(http.StatusUnauthorized, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "Password changed successfully",
	})
}