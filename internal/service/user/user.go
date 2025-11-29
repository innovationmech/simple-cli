package user

import (
	"context"
	"errors"

	"github.com/innovationmech/simple-cli/internal/interfaces"
	"github.com/innovationmech/simple-cli/internal/model"
	"github.com/innovationmech/simple-cli/internal/repository"
)

type UserSrv = interfaces.UserService

type UserServiceConfig struct {
	UserRepository repository.UserRepository
}

type UserServiceOption func(*UserServiceConfig)

type userService struct {
	config *UserServiceConfig
}

func WithUserRepository(repo repository.UserRepository) UserServiceOption {
	return func(config *UserServiceConfig) {
		config.UserRepository = repo
	}
}

func NewUserService(opts ...UserServiceOption) (UserSrv, error) {
	config := &UserServiceConfig{}
	for _, opt := range opts {
		opt(config)
	}
	if config.UserRepository == nil {
		return nil, errors.New("user repository is required")
	}
	return &userService{config: config}, nil
}

func (s *userService) CreateUser(ctx context.Context, user *model.User) error {
	return s.config.UserRepository.CreateUser(ctx, user)
}

func (s *userService) GetUser(ctx context.Context, id string) (*model.User, error) {
	return s.config.UserRepository.GetUser(ctx, id)
}

func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.config.UserRepository.UpdateUser(ctx, user)
}

func (s *userService) DeleteUser(ctx context.Context, id string) error {
	return s.config.UserRepository.DeleteUser(ctx, id)
}
