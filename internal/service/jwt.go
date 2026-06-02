package service

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"time"
	"todo/internal/model"
)

/* Интерфейс репозитория рефреш-токенов */
type TokenRepository interface {
	CreateToken(ctx context.Context, token *model.Token) error
	GetByToken(ctx context.Context, token string) (*model.Token, error)
	RevokeToken(ctx context.Context, token string) error
}

/* Конфиг для конструктора JWT-Сервиса */
type JWTConfig struct {
	JWTSecretA     string
	JWTSecretR     string
	JWTExpirationA time.Duration
	JWTExpirationR time.Duration
}

/* JWT-Сервис */
type ServiceJWT struct {
	tokenRepo   TokenRepository
	secretA     []byte
	secretR     []byte
	expirationA time.Duration
	expirationR time.Duration
}

/* Claims токенов */
type ClaimsToken struct {
	jwt.RegisteredClaims
	TokenType string `json:"token_type"`
}

func NewJWT(tokenRepo TokenRepository, cfg JWTConfig) *ServiceJWT {
	return &ServiceJWT{
		tokenRepo:   tokenRepo,
		secretA:     []byte(cfg.JWTSecretA),
		secretR:     []byte(cfg.JWTSecretR),
		expirationA: cfg.JWTExpirationA,
		expirationR: cfg.JWTExpirationR,
	}
}

func (j *ServiceJWT) GenerateTokens(ctx context.Context, userID string) (string, string, error) {
	accessToken, err := j.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := j.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	token := &model.Token{
		UserID:    userID,
		Token:     refreshToken,
		Revoked:   false,
		ExpiresAt: time.Now().UTC().Add(j.expirationR),
		CreatedAt: time.Now().UTC(),
	}

	if err := j.tokenRepo.CreateToken(ctx, token); err != nil {
		return "", "", fmt.Errorf("failed to save refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}

func (j *ServiceJWT) GenerateAccessToken(userID string) (string, error) {
	claims := &ClaimsToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expirationA)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		TokenType: "access",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretA)
	if err != nil {
		return "", fmt.Errorf("failed to generate access token: %w", err)
	}
	return tokenString, nil
}

func (j *ServiceJWT) ValidateAccessToken(tokenStr string) (string, error) {
	claims := &ClaimsToken{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretA, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse access token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid access token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("access token expired")
	}

	if claims.TokenType != "access" {
		return "", fmt.Errorf("invalid token type: expected access, got %s", claims.TokenType)
	}

	return claims.Subject, nil
}

func (j *ServiceJWT) GenerateRefreshToken(userID string) (string, error) {
	claims := &ClaimsToken{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expirationR)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ID:        uuid.New().String(),
		},
		TokenType: "refresh",
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secretR)
	if err != nil {
		return "", fmt.Errorf("failed to generate refresh token: %w", err)
	}
	return tokenString, nil
}

func (j *ServiceJWT) ValidateRefreshToken(ctx context.Context, tokenStr string) (string, error) {
	claims := &ClaimsToken{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretR, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse refresh token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid refresh token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("refresh token expired")
	}

	if claims.TokenType != "refresh" {
		return "", fmt.Errorf("invalid token type: expected refresh, got %s", claims.TokenType)
	}

	dbToken, err := j.tokenRepo.GetByToken(ctx, tokenStr)
	if err != nil {
		return "", err
	}

	if dbToken.Revoked {
		return "", model.ErrTokenRevoked
	}

	return claims.Subject, nil
}

func (j *ServiceJWT) RevokeRefreshToken(ctx context.Context, tokenStr string) error {
	return j.tokenRepo.RevokeToken(ctx, tokenStr)
}

func (j *ServiceJWT) RefreshTokens(ctx context.Context, refreshTokenStr string) (string, string, error) {
	userID, err := j.ValidateRefreshToken(ctx, refreshTokenStr)
	if err != nil {
		return "", "", err
	}

	if err := j.RevokeRefreshToken(ctx, refreshTokenStr); err != nil {
		return "", "", fmt.Errorf("failed to revoke old refresh token: %w", err)
	}

	return j.GenerateTokens(ctx, userID)
}
