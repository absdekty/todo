package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"todo/internal/model"
)

type RepositoryUser struct {
	*sql.DB
}

func NewUser(db *sql.DB) *RepositoryUser {
	return &RepositoryUser{db}
}

func (r *RepositoryUser) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users
		(id, name, password, role, created)
		VALUES (?, ?, ?, ?, ?)`

	_, err := r.ExecContext(ctx, query, user.ID, user.Name, user.Password, user.Role, user.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return model.ErrUserAlreadyExist
		}
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *RepositoryUser) FindByName(ctx context.Context, name string) (*model.User, error) {
	query := `
		SELECT id, name, password, role, created FROM users
		WHERE name=?`

	var user model.User
	err := r.QueryRowContext(ctx, query, name).Scan(
		&user.ID, &user.Name, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotExist
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}

func (r *RepositoryUser) FindByID(ctx context.Context, userID string) (*model.User, error) {
	query := `
		SELECT id, name, password, role, created FROM users
		WHERE id=?`

	var user model.User
	err := r.QueryRowContext(ctx, query, userID).Scan(
		&user.ID, &user.Name, &user.Password, &user.Role, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotExist
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
