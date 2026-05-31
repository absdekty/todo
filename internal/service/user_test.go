package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"todo/internal/model"
	"todo/internal/repository/mock"
)

type mockHasher struct{}

func (m *mockHasher) Hash(password string) (string, error) {
	return "hashed_" + password, nil
}

func (m *mockHasher) Compare(hashed, password string) bool {
	return hashed == "hashed_"+password
}

func TestCreateUser(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userName  string
		password  string
		wantError error
	}{
		{
			name:      "Обычный юзер",
			userName:  "john",
			password:  "pass123456",
			wantError: nil,
		},
		{
			name:      "Короткое имя",
			userName:  "jo",
			password:  "pass123456",
			wantError: model.ErrUserNameTooShort,
		},
		{
			name:      "Длинное имя",
			userName:  strings.Repeat("a", 30),
			password:  "pass123456",
			wantError: model.ErrUserNameTooLong,
		},
		{
			name:      "Короткий пароль",
			userName:  "john",
			password:  "123",
			wantError: model.ErrUserPasswordTooShort,
		},
		{
			name:      "Длинный пароль",
			userName:  "john",
			password:  strings.Repeat("a", 30),
			wantError: model.ErrUserPasswordTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			service := NewUser(mockRepo, &mockHasher{})

			err := service.CreateUser(ctx, tt.userName, tt.password)
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestFindByName(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userName  string
		setup     func(*mock.MockRepository)
		wantError error
	}{
		{
			name:     "Существует",
			userName: "john",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "random-id", Name: "john"})
			},
			wantError: nil,
		},
		{
			name:      "Не существует",
			userName:  "unknown",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNotExist,
		},
		{
			name:      "Невалидное имя",
			userName:  "jo",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNameTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewUser(mockRepo, &mockHasher{})
			_, err := service.FindByName(ctx, tt.userName)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestLogin(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userName  string
		password  string
		setup     func(*mock.MockRepository)
		wantError error
	}{
		{
			name:     "Существует",
			userName: "john",
			password: "pass123",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{
					ID:       "random-id",
					Name:     "john",
					Password: "hashed_pass123",
				})
			},
			wantError: nil,
		},
		{
			name:     "Неверный пароль",
			userName: "john",
			password: "wrong",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{
					ID:       "random-id",
					Name:     "john",
					Password: "hashed_pass123",
				})
			},
			wantError: model.ErrUserInvalidPW,
		},
		{
			name:      "Не существует",
			userName:  "unknown",
			password:  "pass",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewUser(mockRepo, &mockHasher{})
			_, err := service.Login(ctx, tt.userName, tt.password)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
