package model

import (
	"github.com/google/uuid"
	"strings"
	"time"
)

type Task struct {
	ID          string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

/* Конструктор */
func NewTask(title, description string) (*Task, error) {
	if err := ValidateTitle(title); err != nil {
		return nil, err
	}

	if err := ValidateDescription(description); err != nil {
		return nil, err
	}

	return &Task{
		ID:          uuid.New().String(),
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

/* Валидация */
func ValidateTitle(title string) error {
	title = strings.TrimSpace(title)

	if len(title) < 3 {
		return ErrTaskTitleTooShort
	}

	if len(title) > 15 {
		return ErrTaskTitleTooLong
	}

	return nil
}

func ValidateDescription(description string) error {
	description = strings.TrimSpace(description)

	if len(description) > 200 {
		return ErrTaskDescriptionTooLong
	}

	return nil
}

func ValidateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrTaskInvalidID
	}

	return nil
}

/* Методы */
func (t *Task) SetTitle(title string) error {
	if err := ValidateTitle(title); err != nil {
		return err
	}

	t.Title = strings.TrimSpace(title)
	t.UpdatedAt = time.Now().UTC()

	return nil
}

func (t *Task) SetDescription(description string) error {
	if err := ValidateDescription(description); err != nil {
		return err
	}

	t.Description = strings.TrimSpace(description)
	t.UpdatedAt = time.Now().UTC()

	return nil
}

func (t *Task) SetCompleted(status bool) {
	t.Completed = status
	t.UpdatedAt = time.Now().UTC()
}
