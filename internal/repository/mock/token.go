package mock

import (
	"context"
	"fmt"
	"strings"
	"todo/internal/model"
)

func (m *MockRepository) GenerateTokens(ctx context.Context, userID string) (string, string, error) {
	refreshToken := "refresh_" + userID

	m.tokens[refreshToken] = &model.Token{UserID: userID, Token: refreshToken}

	return "access_" + userID, refreshToken, nil
}

func (m *MockRepository) ValidateAccessToken(tokenStr string) (string, error) {
	prefix := "access_"
	if !strings.HasPrefix(tokenStr, prefix) {
		return "", fmt.Errorf("invalid token format")
	}
	return strings.TrimPrefix(tokenStr, prefix), nil
}

func (m *MockRepository) ValidateRefreshToken(tokenStr string) (string, error) {
	prefix := "refresh_"
	if !strings.HasPrefix(tokenStr, prefix) {
		return "", fmt.Errorf("invalid token format")
	}

	if _, ok := m.tokens[tokenStr]; !ok {
		return "", model.ErrTokenNotExist
	}

	return strings.TrimPrefix(tokenStr, prefix), nil
}

func (m *MockRepository) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	if _, ok := m.tokens[tokenStr]; !ok {
		return model.ErrTokenNotExist
	}

	delete(m.tokens, tokenStr)
	return nil
}

func (m *MockRepository) RefreshTokens(ctx context.Context, refreshTokenStr string) (string, string, error) {
	userID, err := m.ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return "", "", err
	}

	if err := m.RevokeRefreshToken(ctx, refreshTokenStr); err != nil {
		return "", "", err
	}

	return m.GenerateTokens(ctx, userID)
}
