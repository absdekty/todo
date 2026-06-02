package sqlite

import (
	"context"
	"errors"
	"testing"
	"todo/internal/model"
)

func TestCreateToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		token     *model.Token
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:  "Валидный токен",
			token: &model.Token{UserID: "userid", Token: "token"},
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
			},
			wantError: nil,
		},
		{
			name:  "Идентичный token",
			token: &model.Token{Token: "token"},
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.token.CreateToken(context.Background(), &model.Token{UserID: "userid", Token: "token"})
			},
			wantError: model.ErrTokenAlreadyExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			err := repo.token.CreateToken(ctx, tt.token)
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

func TestGetByToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		token     string
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:  "Токен существует",
			token: "token",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.token.CreateToken(context.Background(), &model.Token{UserID: "userid", Token: "token"})
			},
			wantError: nil,
		},
		{
			name:      "Токен не существует",
			token:     "token",
			setup:     func(repo *dbRepo) {},
			wantError: model.ErrTokenNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			_, err := repo.token.GetByToken(ctx, tt.token)
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

func TestRevokeToken(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		token     string
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:  "Токен существует",
			token: "token",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.token.CreateToken(context.Background(), &model.Token{UserID: "userid", Token: "token"})
			},
			wantError: nil,
		},
		{
			name:      "Токен не существует",
			token:     "token",
			setup:     func(repo *dbRepo) {},
			wantError: model.ErrTokenNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			err := repo.token.RevokeToken(ctx, tt.token)
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
