// internal/handler/user_handler.go
package handler

import (
	"math"
	"net/http"
	"strconv"
	"strings"

	"github.com/akhmadzaqiriyadi/stmadb-portal-go/internal/service"
	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db" // <-- Tambahkan import ini
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// PERBAIKAN: Pindahkan fungsi ToProfileDTO ke sini
func ToProfileDTO(user db.UserModel) ProfileData {
	return ProfileData{
		ID:        int64(user.ID),
		Username:  user.Username,
		Role:      string(user.Role),
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt.String(),
	}
}

// GetUsers godoc
// @Summary      Get all users with filters and pagination
// @Description  Retrieves a list of all users. Only accessible by admins.
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        page query int false "Page number"
// @Param        limit query int false "Items per page"
// @Param        search query string false "Search by username"
// @Param        role query string false "Filter by role (admin, teacher, student, staff)"
// @Param        is_active query boolean false "Filter by active status (true/false)"
// @Success      200 {object}  GenericResponse "List of users"
// @Failure      500 {object}  GenericResponse "Internal Server Error"
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	var filters UserQueryFilters
	if err := c.ShouldBindQuery(&filters); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid query parameters"})
		return
	}

	if filters.Page <= 0 {
		filters.Page = 1
	}
	if filters.Limit <= 0 {
		filters.Limit = 10
	}

	var isActive *bool
	if strings.ToLower(filters.IsActive) == "true" {
		val := true
		isActive = &val
	} else if strings.ToLower(filters.IsActive) == "false" {
		val := false
		isActive = &val
	}

	params := service.GetUsersParams{
		Page:     filters.Page,
		Limit:    filters.Limit,
		Search:   filters.Search,
		Role:     filters.Role,
		IsActive: isActive,
	}

	users, total, err := h.service.GetUsers(params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	var userProfiles []ProfileData
	for _, user := range users {
		userProfiles = append(userProfiles, ToProfileDTO(user))
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Users retrieved successfully",
		"data":    userProfiles,
		"meta": gin.H{
			"page":       filters.Page,
			"limit":      filters.Limit,
			"total":      total,
			"totalPages": int(math.Ceil(float64(total) / float64(filters.Limit))),
		},
	})
}

// GetUserByID godoc
// @Summary      Get a single user by ID
// @Description  Retrieves details of a single user.
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} GenericResponse{data=ProfileData} "User details"
// @Failure      404 {object} GenericResponse "User not found"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid user ID"})
		return
	}

	user, err := h.service.GetUserByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "User retrieved successfully",
		Data:    ToProfileDTO(*user),
	})
}

// CreateUser godoc
// @Summary      Create a new user
// @Description  Creates a new user account. Only accessible by admins.
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        user body CreateUserRequest true "New User Data"
// @Success      201 {object} GenericResponse{data=ProfileData} "User created successfully"
// @Failure      400 {object} GenericResponse "Invalid request body or username exists"
// @Router       /users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	user, err := h.service.CreateUser(req.Username, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, GenericResponse{
		Success: true,
		Message: "User created successfully",
		Data:    ToProfileDTO(*user),
	})
}

// UpdateUser godoc
// @Summary      Update a user
// @Description  Updates a user's role or active status.
// @Tags         Users
// @Security     BearerAuth
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Param        user body UpdateUserRequest true "User Update Data"
// @Success      200 {object} GenericResponse{data=ProfileData} "User updated successfully"
// @Failure      404 {object} GenericResponse "User not found"
// @Router       /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid user ID"})
		return
	}
	
	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	// Buat pointer untuk role jika tidak kosong
	var rolePtr *string
	if req.Role != "" {
		rolePtr = &req.Role
	}

	user, err := h.service.UpdateUser(id, rolePtr, req.IsActive)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "User updated successfully",
		Data:    ToProfileDTO(*user),
	})
}

// DeleteUser godoc
// @Summary      Delete a user
// @Description  Deletes a user account.
// @Tags         Users
// @Security     BearerAuth
// @Produce      json
// @Param        id   path      int  true  "User ID"
// @Success      200 {object} GenericResponse "User deleted successfully"
// @Failure      404 {object} GenericResponse "User not found"
// @Router       /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, GenericResponse{Success: false, Message: "Invalid user ID"})
		return
	}

	err = h.service.DeleteUser(id)
	if err != nil {
		c.JSON(http.StatusNotFound, GenericResponse{Success: false, Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, GenericResponse{
		Success: true,
		Message: "User deleted successfully",
	})
}