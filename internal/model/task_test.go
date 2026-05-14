package model

import (
	"github.com/google/uuid"
	"strings"
	"testing"
)

func TestNewTask(t *testing.T) {
	tests := []struct {
		name        string
		title       string
		description string
		wantErr     error
	}{
		{
			name:        "Валидные данные",
			title:       "Title",
			description: "",
			wantErr:     nil,
		},
		{
			name:        "Короткое название",
			title:       "..",
			description: "",
			wantErr:     ErrTaskTitleTooShort,
		},
		{
			name:        "Длинное название",
			title:       strings.Repeat("a", 16),
			description: "",
			wantErr:     ErrTaskTitleTooLong,
		},
		{
			name:        "Пустое описание",
			title:       "Title",
			description: "",
			wantErr:     nil,
		},
		{
			name:        "Длинное описание",
			title:       "Title",
			description: strings.Repeat("a", 201),
			wantErr:     ErrTaskDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewTask(tt.title, tt.description)
			if err != tt.wantErr {
				t.Errorf("expected err %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateTitle(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr error
	}{
		{
			name:    "Валидные данные",
			data:    "Title",
			wantErr: nil,
		},
		{
			name:    "Короткое название",
			data:    "..",
			wantErr: ErrTaskTitleTooShort,
		},
		{
			name:    "Длинное название",
			data:    strings.Repeat("a", 16),
			wantErr: ErrTaskTitleTooLong,
		},
		{
			name:    "Просто пробелы",
			data:    "     ",
			wantErr: ErrTaskTitleTooShort,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateTitle(tt.data)
			if err != tt.wantErr {
				t.Errorf("expected err %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateDescription(t *testing.T) {
	tests := []struct {
		name    string
		data    string
		wantErr error
	}{
		{
			name:    "Валидные данные",
			data:    "description",
			wantErr: nil,
		},
		{
			name:    "Пустое описание",
			data:    "",
			wantErr: nil,
		},
		{
			name:    "Длинное описание",
			data:    strings.Repeat("a", 201),
			wantErr: ErrTaskDescriptionTooLong,
		},
		{
			name:    "Просто пробелы",
			data:    "     ",
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateDescription(tt.data)
			if err != tt.wantErr {
				t.Errorf("expected err %v, but got %v", tt.wantErr, err)
			}
		})
	}
}

func TestValidateID(t *testing.T) {
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
			err := ValidateID(tt.data)
			if err != nil && !tt.wantErr {
				t.Errorf("not expected err, but got %v", err)
			}
		})
	}
}
