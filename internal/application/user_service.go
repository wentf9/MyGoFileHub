package application

import (
	"context"
	"errors"

	"github.com/wentf9/MyGoFileHub/internal/domain/model"
	"github.com/wentf9/MyGoFileHub/internal/domain/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{userRepo: repo}
}

func (s *UserService) ListUsers(ctx context.Context) ([]*model.User, error) {
	return s.userRepo.FindAll(ctx)
}

func (s *UserService) GetUser(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.FindByID(ctx, id)
}

func (s *UserService) CreateUser(ctx context.Context, username, password, role string) error {
	// 检查用户是否已存在
	if _, err := s.userRepo.FindByUsername(ctx, username); err == nil {
		return errors.New("user already exists")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		IsActive:     true,
	}

	return s.userRepo.Save(ctx, user)
}

func (s *UserService) UpdateUser(ctx context.Context, id uint, password, role string, isActive *bool) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}

	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hash)
	}

	if role != "" {
		user.Role = role
	}

	if isActive != nil {
		user.IsActive = *isActive
	}

	return s.userRepo.Save(ctx, user)
}

func (s *UserService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}
