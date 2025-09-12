// internal/handler/auth_handler.go
package handler

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/service"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db" // Prisma client
)

type AuthHandler struct {
	service *service.AuthService
}

func NewAuthHandler(service *service.AuthService) *AuthHandler {
	return &AuthHandler{service: service}
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// Login handles user login requests.
// @Summary      User login
// @Description  Authenticate a user and return a JWT token.
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        credentials body LoginRequest true "Login Credentials"
// @Success      200  {object}  GenericResponse{data=LoginResponse} "Login Successful"
// @Failure      400  {object}  GenericResponse "Invalid Request"
// @Failure      401  {object}  GenericResponse "Unauthorized"
// @Router       /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
    // ... (logika fungsi tidak berubah)
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid request body"})
		return
	}

	token, err := h.service.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "Login successful",
		Data:    LoginResponse{AccessToken: token},
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
	
	// PENTING: Jangan kirim seluruh objek user dari database yang berisi password.
	// Buat objek DTO yang aman.
	profileData := ProfileData{
		ID:        int64(user.ID), // Konversi BigInt ke int64
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