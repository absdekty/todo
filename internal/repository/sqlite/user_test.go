package sqlite

import (
	"context"
	"errors"
	"testing"
	"todo/internal/model"
)

func TestCreateUser(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		user      *model.User
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:      "Валидный юзер",
			user:      &model.User{ID: "userid", Name: "name"},
			setup:     func(repo *dbRepo) {},
			wantError: nil,
		},
		{
			name: "name прошлого",
			user: &model.User{ID: "userid", Name: "name"},
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid", Name: "name"})
			},
			wantError: model.ErrUserAlreadyExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			err := repo.user.CreateUser(ctx, tt.user)
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
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		userName  string
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:     "Существует",
			userName: "name",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid", Name: "name"})
			},
			wantError: nil,
		},
		{
			name:      "Не существует",
			userName:  "name",
			setup:     func(repo *dbRepo) {},
			wantError: model.ErrUserNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			_, err := repo.user.FindByName(ctx, tt.userName)
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
