// internal/service/user_service.go
package service

import (
	"context"
	"errors"

	"github.com/akhmadzaqiriyadi/stmadb-portal-go/prisma/db"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	db *db.PrismaClient
}

func NewUserService(db *db.PrismaClient) *UserService {
	return &UserService{db: db}
}

// GetUsersParams adalah struct untuk parameter GetUsers.
type GetUsersParams struct {
	Page     int
	Limit    int
	Search   string
	Role     string
	IsActive *bool // Pointer agar bisa handle true/false/nil
}

// GetUsers mengambil daftar pengguna dengan paginasi dan filter.
func (s *UserService) GetUsers(params GetUsersParams) ([]db.UserModel, int, error) {
	var where []db.UserWhereParam

	// Terapkan filter jika ada
	if params.Search != "" {
		where = append(where, db.User.Username.Contains(params.Search))
	}
	if params.Role != "" {
		where = append(where, db.User.Role.Equals(db.UserRole(params.Role)))
	}
	if params.IsActive != nil {
		where = append(where, db.User.IsActive.Equals(*params.IsActive))
	}

	// Hitung total dengan FindMany karena Aggregate tidak tersedia.
	allUsers, err := s.db.User.FindMany(where...).Exec(context.Background())
	if err != nil {
		return nil, 0, errors.New("failed to count users")
	}
	total := len(allUsers)

	// Ambil data dengan paginasi
	users, err := s.db.User.FindMany(where...).
		OrderBy(db.User.CreatedAt.Order(db.SortOrderDesc)).
		Skip((params.Page - 1) * params.Limit).
		Take(params.Limit).
		Exec(context.Background())
		
	if err != nil {
		return nil, 0, errors.New("failed to retrieve users")
	}
	return users, total, nil
}

// GetUserByID mengambil satu pengguna berdasarkan ID.
func (s *UserService) GetUserByID(id int) (*db.UserModel, error) {
	user, err := s.db.User.FindUnique(db.User.ID.Equals(db.BigInt(id))).Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return user, nil
}

// CreateUser membuat pengguna baru.
func (s *UserService) CreateUser(username, password, role string) (*db.UserModel, error) {
	_, err := s.db.User.FindFirst(db.User.Username.Equals(username)).Exec(context.Background())
	if !errors.Is(err, db.ErrNotFound) {
		return nil, errors.New("username already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	return s.db.User.CreateOne(
		db.User.Username.Set(username),
		db.User.Password.Set(string(hashedPassword)),
		db.User.Role.Set(db.UserRole(role)),
	).Exec(context.Background())
}

// UpdateUser memperbarui data pengguna.
func (s *UserService) UpdateUser(id int, role *string, isActive *bool) (*db.UserModel, error) {
	var params []db.UserSetParam
	if role != nil && *role != "" {
		params = append(params, db.User.Role.Set(db.UserRole(*role)))
	}
	if isActive != nil {
		params = append(params, db.User.IsActive.Set(*isActive))
	}

	updatedUser, err := s.db.User.FindUnique(db.User.ID.Equals(db.BigInt(id))).Update(params...).Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	return updatedUser, nil
}

// DeleteUser menghapus pengguna.
func (s *UserService) DeleteUser(id int) error {
	_, err := s.db.User.FindUnique(db.User.ID.Equals(db.BigInt(id))).Delete().Exec(context.Background())
	if err != nil {
		if errors.Is(err, db.ErrNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	return nil
}