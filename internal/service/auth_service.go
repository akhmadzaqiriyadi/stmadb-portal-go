// internal/service/auth_service.go
package service

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"

	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db" // Prisma client
)

type AuthService struct {
	db *db.PrismaClient
}

func NewAuthService(db *db.PrismaClient) *AuthService {
	return &AuthService{db: db}
}

// Definisikan tipe data baru untuk menampung kedua token
type TokenPair struct {
	AccessToken  string
	RefreshToken string
}

// Login sekarang mengembalikan sepasang token
func (s *AuthService) Login(username, password string) (*TokenPair, error) {
	user, err := s.db.User.FindFirst(db.User.Username.Equals(username)).Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, errors.New("invalid credentials")
		}
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Buat Access Token (Masa berlaku pendek, misal 15 menit)
	accessClaims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(), // <-- Masa berlaku 15 menit
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Buat Refresh Token (Masa berlaku panjang, misal 7 hari)
	refreshClaims := jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(), // <-- Masa berlaku 7 hari
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	refreshTokenString, err := refreshToken.SignedString([]byte(viper.GetString("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: refreshTokenString,
	}, nil
}

// RefreshToken memvalidasi refresh token dan membuat sepasang token baru
func (s *AuthService) RefreshToken(refreshTokenString string) (*TokenPair, error) {
	token, err := jwt.Parse(refreshTokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(viper.GetString("JWT_REFRESH_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New("invalid or expired refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIdFloat := claims["userId"].(float64)
	userId := int(userIdFloat)

	user, err := s.db.User.FindUnique(db.User.ID.Equals(db.BigInt(userId))).Exec(context.Background())
	if err != nil {
		return nil, errors.New("user not found")
	}
	
	// Buat Access Token baru
	accessClaims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString([]byte(viper.GetString("JWT_SECRET")))
	if err != nil {
		return nil, err
	}

	// Buat Refresh Token baru
	refreshClaims := jwt.MapClaims{
		"userId": user.ID,
		"exp":    time.Now().Add(time.Hour * 24 * 7).Unix(),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	newRefreshTokenString, err := refreshToken.SignedString([]byte(viper.GetString("JWT_REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}

	return &TokenPair{
		AccessToken:  accessTokenString,
		RefreshToken: newRefreshTokenString,
	}, nil
}

// GetProfile mengambil profil pengguna berdasarkan userID.
func (s *AuthService) GetProfile(userID int) (*db.UserModel, error) {
	user, err := s.db.User.FindUnique(
		db.User.ID.Equals(db.BigInt(userID)),
	).With(
		db.User.Teacher.Fetch(),
		db.User.Student.Fetch().With(
			db.Student.CurrentClass.Fetch(),
		),
	).Exec(context.Background())

	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}

	return user, nil
}

// ChangePassword memvalidasi password lama dan menggantinya dengan yang baru.
func (s *AuthService) ChangePassword(userID int, currentPassword, newPassword string) error {
	// 1. Ambil data user dari database
	user, err := s.db.User.FindUnique(
		db.User.ID.Equals(db.BigInt(userID)),
	).Exec(context.Background())
	if err != nil {
		return errors.New("user not found")
	}

	// 2. Bandingkan password saat ini
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(currentPassword))
	if err != nil {
		return errors.New("current password is incorrect") // Password lama tidak cocok
	}

	// 3. Hash password baru
	hashedNewPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash new password")
	}

	// 4. Update password di database
	_, err = s.db.User.FindUnique(
		db.User.ID.Equals(db.BigInt(userID)),
	).Update(
		db.User.Password.Set(string(hashedNewPassword)),
	).Exec(context.Background())
	
	if err != nil {
		return errors.New("failed to update password")
	}

	return nil
}