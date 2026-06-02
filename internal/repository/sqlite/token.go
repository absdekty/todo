package sqlite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"todo/internal/model"
)

type RepositoryToken struct {
	*sql.DB
}

func NewToken(db *sql.DB) *RepositoryToken {
	return &RepositoryToken{db}
}

func (r *RepositoryToken) CreateToken(ctx context.Context, token *model.Token) error {
	query := `INSERT INTO tokens
			(userid, token, revoked, expires, created)
			VALUES(?, ?, ?, ?, ?)`

	res, err := r.ExecContext(ctx, query, token.UserID, token.Token, token.Revoked, token.ExpiresAt, token.CreatedAt)
	if err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			return model.ErrTokenAlreadyExist
		}
		return fmt.Errorf("failed to create token: %w", err)
	}

	if rows, _ := res.RowsAffected(); rows == 0 {

	}

	return nil
}

func (r *RepositoryToken) GetByToken(ctx context.Context, token string) (*model.Token, error) {
	query := `SELECT
		userid, token, revoked, expires, created
		FROM tokens
		WHERE token=?`

	var Token model.Token
	err := r.QueryRowContext(ctx, query, token).Scan(
		&Token.UserID, &Token.Token, &Token.Revoked, &Token.ExpiresAt, &Token.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrTokenNotExist
		}
		return nil, fmt.Errorf("failed to get token: %w", err)
	}

	return &Token, nil
}

func (r *RepositoryToken) RevokeToken(ctx context.Context, token string) error {
	query := `UPDATE tokens
		SET revoked=1
		WHERE token=?`

	res, err := r.ExecContext(ctx, query, token)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return model.ErrTokenNotExist
	}

	return nil
}
