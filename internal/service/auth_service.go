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

// Login memvalidasi kredensial user dan mengembalikan token jika valid
func (s *AuthService) Login(username, password string) (string, error) {
	// 1. Cari user berdasarkan username
	user, err := s.db.User.FindFirst(db.User.Username.Equals(username)).Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return "", errors.New("invalid credentials")
		}
		return "", err
	}

	// 2. Bandingkan password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid credentials") // Password tidak cocok
	}

	// 3. Buat JWT Token
	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"role":     user.Role,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), // Token berlaku 24 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	jwtSecret := viper.GetString("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

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