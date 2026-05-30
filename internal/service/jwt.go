package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JWTConfig struct {
	JWTSecret     string
	JWTExpiration time.Duration
}

type ServiceJWT struct {
	secret     []byte
	expiration time.Duration
}

func NewJWT(cfg JWTConfig) *ServiceJWT {
	return &ServiceJWT{
		secret:     []byte(cfg.JWTSecret),
		expiration: cfg.JWTExpiration,
	}
}

func (j *ServiceJWT) GenerateToken(userID string) (string, error) {
	claims := &jwt.RegisteredClaims{
		Subject:   userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.expiration)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(j.secret)
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %w", err)
	}
	return tokenString, nil
}

func (j *ServiceJWT) ValidateToken(tokenStr string) (string, error) {
	claims := &jwt.RegisteredClaims{}

	token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secret, nil
	})

	if err != nil {
		return "", fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return "", fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt != nil && claims.ExpiresAt.Before(time.Now()) {
		return "", fmt.Errorf("token expired")
	}

	return claims.Subject, nil
}
