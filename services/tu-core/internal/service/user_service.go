package service

import (
	"context"
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/google/uuid"
	"sekolah-platform/services/tu-core/internal/repository"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) ListUsers(ctx context.Context, schoolID string) ([]repository.User, error) {
	return s.userRepo.ListUsersBySchool(ctx, schoolID)
}

func (s *UserService) CreateUser(ctx context.Context, schoolID, username, password, role string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password: %w", err)
	}

	user := &repository.User{
		UserID:   uuid.NewString(),
		Name:     username,
		Password: string(hash),
		Role:     role,
		SchoolID: schoolID,
		IsActive: true,
	}

	return s.userRepo.CreateUser(ctx, user)
}

func (s *UserService) UpdateRole(ctx context.Context, userID, role string) error {
	if role == "" {
		return errors.New("role tidak boleh kosong")
	}
	return s.userRepo.UpdateUserRole(ctx, userID, role)
}

func (s *UserService) ToggleStatus(ctx context.Context, userID string, isActive bool) error {
	return s.userRepo.ToggleUserStatus(ctx, userID, isActive)
}

// ResetPassword (BARU)
func (s *UserService) ResetPassword(ctx context.Context, userID, newPassword string) error {
	if len(newPassword) < 6 {
		return errors.New("password minimal 6 karakter")
	}
	
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal mengenkripsi password: %w", err)
	}
	
	return s.userRepo.UpdatePassword(ctx, userID, string(hash))
}
