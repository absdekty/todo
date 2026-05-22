package model

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestNewUser(t *testing.T) {
	tests := []struct {
		name     string
		uname    string
		password string
		wantErr  error
	}{
		{
			name:     "Валидные данные",
			uname:    "UserName",
			password: "UserPassword",
			wantErr:  nil,
		},
		{
			name:     "Короткое имя",
			uname:    "..",
			password: "",
			wantErr:  ErrUserNameTooShort,
		},
		{
			name:     "Длинное имя",
			uname:    strings.Repeat("a", 16),
			password: "",
			wantErr:  ErrUserNameTooLong,
		},
		// Валидация пароля происходит в сервисе; структура содержит HASH пароля.
		{
			name:     "Короткий пароль",
			uname:    "UserName",
			password: "..",
			wantErr:  nil,
		},
		{
			name:     "Длинный пароль",
			uname:    "UserName",
			password: strings.Repeat("a", 30),
			wantErr:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewUser(tt.uname, tt.password)
			if err != tt.wantErr {
				t.Errorf("expected err %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserValidateName(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr error
	}{
		{
			name:    "Валидные данные",
			data:    "UserName",
			wantErr: nil,
		},
		{
			name:    "Короткое имя",
			data:    "..",
			wantErr: ErrUserNameTooShort,
		},
		{
			name:    "Длинное имя",
			data:    strings.Repeat("a", 30),
			wantErr: ErrUserNameTooLong,
		},
		{
			name:    "Просто пробелы",
			data:    "     ",
			wantErr: ErrUserNameTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserValidateName(tt.data)
			if err != tt.wantErr {
				t.Errorf("expected err %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr error
	}{
		{
			name:    "Валидные данные",
			data:    "UserPassword",
			wantErr: nil,
		},
		{
			name:    "Пустой пароль",
			data:    "",
			wantErr: ErrUserPasswordTooShort,
		},
		{
			name:    "Длинный пароль",
			data:    strings.Repeat("a", 30),
			wantErr: ErrUserPasswordTooLong,
		},
		{
			name:    "Просто пробелы",
			data:    "     ",
			wantErr: ErrUserPasswordTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserValidatePassword(tt.data)
			if err != tt.wantErr {
				t.Errorf("expected err %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestUserValidateID(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr bool
	}{
		{
			name:    "Валидный ID",
			data:    uuid.New().String(),
			wantErr: false,
		},
		{
			name:    "Только пробелы",
			data:    "  ",
			wantErr: true,
		},
		{
			name:    "Пустой ID",
			data:    "",
			wantErr: true,
		},
		{
			name:    "ID =32 символа",
			data:    strings.Repeat("a", 32),
			wantErr: false,
		},
		{
			name:    "ID =32 символа",
			data:    strings.Repeat("z", 32),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := UserValidateID(tt.data)
			if err != nil && !tt.wantErr {
				t.Errorf("not expected err, but got %v", err)
			}
		})
	}
}
