package service

import (
	"context"
	"fmt"
	"todo/internal/model"
	"todo/internal/repository"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
	Compare(hashedPassword, password string) bool
}

type ServiceUserI interface {
	CreateUser(ctx context.Context, name, password string) error
	FindByName(ctx context.Context, name string) (*model.User, error)
	Login(ctx context.Context, name, password string) (string, error)
}

var _ ServiceUserI = (*ServiceUser)(nil)

type ServiceUser struct {
	repo   repository.RepositoryUser
	hasher PasswordHasher
}

func NewUser(repo repository.RepositoryUser, hasher PasswordHasher) *ServiceUser {
	return &ServiceUser{
		repo:   repo,
		hasher: hasher,
	}
}

func (s *ServiceUser) CreateUser(ctx context.Context, name, password string) error {
	if err := model.UserValidatePassword(password); err != nil {
		return fmt.Errorf("create user (validate): %w", err)
	}

	hashedPassword, err := s.hasher.Hash(password)
	if err != nil {
		return fmt.Errorf("create user (hash pw): %w", err)
	}

	user, err := model.NewUser(name, hashedPassword)
	if err != nil {
		return fmt.Errorf("create user (new user): %w", err)
	}

	if err = s.repo.CreateUser(ctx, user); err != nil {
		return fmt.Errorf("create user (repo): %w", err)
	}

	return nil
}

func (s *ServiceUser) FindByName(ctx context.Context, name string) (*model.User, error) {
	if err := model.UserValidateName(name); err != nil {
		return nil, fmt.Errorf("find by name (validate): %w", err)
	}

	user, err := s.repo.FindByName(ctx, name)
	if err != nil {
		return nil, fmt.Errorf("find by name (repo): %w", err)
	}
	return user, nil
}

func (s *ServiceUser) Login(ctx context.Context, name, password string) (string, error) {
	user, err := s.FindByName(ctx, name)
	if err != nil {
		return "", fmt.Errorf("login (find): %w", err)
	}

	if !s.hasher.Compare(user.Password, password) {
		return "", model.ErrUserInvalidPW
	}

	return user.ID, nil
}
