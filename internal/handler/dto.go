// internal/handler/dto.go
package handler

// GenericResponse is the base structure for all JSON responses.
type GenericResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Operation successful"`
	Data    interface{} `json:"data,omitempty"`
}

// LoginRequest is the structure for a login request.
type LoginRequest struct {
	Username string `json:"username" binding:"required" example:"admin"`
	Password string `json:"password" binding:"required" example:"admin123"`
}

// RefreshTokenRequest is the structure for a refresh token request.
type RefreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ChangePasswordRequest is the structure for a change password request.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"currentPassword" binding:"required" example:"old_password123"`
	NewPassword     string `json:"newPassword" binding:"required,min=6" example:"new_secure_password"`
}

// TokenResponse is the data structure for all responses that return tokens.
type TokenResponse struct {
	AccessToken  string `json:"accessToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	RefreshToken string `json:"refreshToken" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
}

// ProfileData is the sanitized user profile data sent to the client.
type ProfileData struct {
	ID        int64  `json:"id" example:"1"`
	Username  string `json:"username" example:"admin"`
	Role      string `json:"role" example:"admin"`
	IsActive  bool   `json:"is_active" example:"true"`
	CreatedAt string `json:"created_at" example:"2025-09-13T12:00:00Z"`
}

// ProfileResponse is the full structure for the get profile response.
type ProfileResponse struct {
	Success bool        `json:"success" example:"true"`
	Message string      `json:"message" example:"Profile retrieved successfully"`
	Data    ProfileData `json:"data"`
}