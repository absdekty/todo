package model

import (
	"github.com/google/uuid"
	"strings"
	"time"
	"unicode/utf8"
)

const defaultRole = "user"

type User struct {
	ID        string
	Name      string
	Password  string
	Role      string
	CreatedAt time.Time
}

func NewUser(name, password string) (*User, error) {
	if err := UserValidateName(name); err != nil {
		return nil, err
	}

	return &User{
		ID:        uuid.New().String(),
		Name:      strings.TrimSpace(name),
		Password:  strings.TrimSpace(password),
		Role:      defaultRole,
		CreatedAt: time.Now().UTC(),
	}, nil
}

/* Валидация */
func UserValidateName(name string) error {
	name = strings.TrimSpace(name)

	runeCount := utf8.RuneCountInString(name)

	if runeCount < 3 {
		return ErrUserNameTooShort
	}

	if runeCount > 10 {
		return ErrUserNameTooLong
	}

	return nil
}

// Проверка ПЕРЕД хешированием | Структура уже содержит Hash
func UserValidatePassword(password string) error {
	password = strings.TrimSpace(password)

	runeCount := utf8.RuneCountInString(password)

	if runeCount < 8 {
		return ErrUserPasswordTooShort
	}

	if runeCount > 16 {
		return ErrUserPasswordTooLong
	}

	return nil
}

func UserValidateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrUserInvalidID
	}

	return nil
}
