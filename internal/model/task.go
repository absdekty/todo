package model

import (
	"github.com/google/uuid"
	"strings"
	"time"
	"unicode/utf8"
)

type Task struct {
	ID          string
	UserID      string
	Title       string
	Description string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

/* Конструктор */
func NewTask(userID, title, description string) (*Task, error) {
	if err := TaskValidateTitle(title); err != nil {
		return nil, err
	}

	if err := TaskValidateDescription(description); err != nil {
		return nil, err
	}

	return &Task{
		ID:          uuid.New().String(),
		UserID:      userID,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}, nil
}

/* Валидация */
func TaskValidateTitle(title string) error {
	title = strings.TrimSpace(title)

	runeCount := utf8.RuneCountInString(title)

	if runeCount < 3 {
		return ErrTaskTitleTooShort
	}

	if runeCount > 15 {
		return ErrTaskTitleTooLong
	}

	return nil
}

func TaskValidateDescription(description string) error {
	description = strings.TrimSpace(description)

	runeCount := utf8.RuneCountInString(description)

	if runeCount > 200 {
		return ErrTaskDescriptionTooLong
	}

	return nil
}

func TaskValidateID(id string) error {
	if _, err := uuid.Parse(id); err != nil {
		return ErrTaskInvalidID
	}

	return nil
}

/* Методы */
func (t *Task) SetTitle(title string) error {
	if err := TaskValidateTitle(title); err != nil {
		return err
	}

	t.Title = strings.TrimSpace(title)
	t.UpdatedAt = time.Now().UTC()

	return nil
}

func (t *Task) SetDescription(description string) error {
	if err := TaskValidateDescription(description); err != nil {
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
