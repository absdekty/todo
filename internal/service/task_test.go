package service

import (
	"context"
	"errors"
	"strings"
	"testing"
	"todo/internal/model"
	"todo/internal/repository/mock"
)

func TestCreateTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name        string
		userID      string
		title       string
		description string
		setup       func(*mock.MockRepository)
		wantError   error
	}{
		{
			name:        "Успешное создание",
			userID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			title:       "Test Task",
			description: "",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: nil,
		},
		{
			name:        "Невалидный userID",
			userID:      "invalid",
			title:       "Test Task",
			description: "",
			setup:       func(m *mock.MockRepository) {},
			wantError:   model.ErrUserInvalidID,
		},
		{
			name:        "Пользователь не существует",
			userID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			title:       "Test Task",
			description: "",
			setup:       func(m *mock.MockRepository) {},
			wantError:   model.ErrUserNotExist,
		},
		{
			name:        "Короткий титул",
			userID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			title:       "Te",
			description: "",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskTitleTooShort,
		},
		{
			name:        "Длинный титул",
			userID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			title:       strings.Repeat("a", 30),
			description: "",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskTitleTooLong,
		},
		{
			name:        "Длинное описание",
			userID:      "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			title:       "Test Task",
			description: strings.Repeat("a", 250),
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskDescriptionTooLong,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewTask(mockRepo, mockRepo)
			_, err := service.CreateTask(ctx, tt.userID, tt.title, tt.description)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestGetUserTasks(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		setup     func(*mock.MockRepository)
		wantError error
	}{
		{
			name:   "Успешное получение",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: nil,
		},
		{
			name:      "Невалидный userID",
			userID:    "invalid",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserInvalidID,
		},
		{
			name:      "Пользователь не существует",
			userID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewTask(mockRepo, mockRepo)
			_, err := service.GetUserTasks(ctx, tt.userID)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		taskID    string
		setup     func(*mock.MockRepository)
		wantError error
	}{
		{
			name:   "Успешное получение",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
				m.CreateTask(ctx, &model.Task{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", UserID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Title: "Task"})
			},
			wantError: nil,
		},
		{
			name:      "Невалидный ID пользователя",
			userID:    "invalid",
			taskID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserInvalidID,
		},
		{
			name:      "Пользователь не существует",
			userID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNotExist,
		},
		{
			name:   "Невалидный ID задачи",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "..",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskInvalidID,
		},
		{
			name:   "Задача не существует",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskNotExist,
		},
		{
			name:   "Чужая задача",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
				m.CreateTask(ctx, &model.Task{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", UserID: "cccccccc-cccc-cccc-cccc-cccccccccccc", Title: "Another Task"})
			},
			wantError: model.ErrUserForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewTask(mockRepo, mockRepo)
			_, err := service.GetTaskByID(ctx, tt.userID, tt.taskID)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestUpdateTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		task      *model.Task
		setup     func(*mock.MockRepository)
		wantError error
	}{
		{
			name: "Успешное обновление",
			task: &model.Task{
				ID:     "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				UserID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				Title:  "New Title",
			},
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
				m.CreateTask(ctx, &model.Task{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", UserID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Title: "Old Title"})
			},
			wantError: nil,
		},
		{
			name:      "Nil задача",
			task:      nil,
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrTaskIsNil,
		},
		{
			name: "Невалидный ID пользователя",
			task: &model.Task{
				ID:     "task-123",
				UserID: "invalid",
				Title:  "Updated Title",
			},
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserInvalidID,
		},
		{
			name: "Пользователь не существует",
			task: &model.Task{
				ID:     "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				UserID: "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				Title:  "Updated Title",
			},
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNotExist,
		},
		{
			name: "Задача не существует",
			task: &model.Task{
				ID:     "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb",
				UserID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				Title:  "Updated Title",
			},
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskNotExist,
		},
		{
			name: "Обновление чужой задачи",
			task: &model.Task{
				ID:     "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				UserID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
				Title:  "Updated Title",
			},
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
				m.CreateTask(ctx, &model.Task{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", UserID: "cccccccc-cccc-cccc-cccc-cccccccccccc", Title: "Another Task"})
			},
			wantError: model.ErrUserForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewTask(mockRepo, mockRepo)
			err := service.UpdateTask(ctx, tt.task)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	ctx := context.Background()

	tests := []struct {
		name      string
		userID    string
		taskID    string
		setup     func(*mock.MockRepository)
		wantError error
	}{
		{
			name:   "Успешное удаление",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
				m.CreateTask(ctx, &model.Task{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", UserID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Title: "Task"})
			},
			wantError: nil,
		},
		{
			name:      "Невалидный ID пользователя",
			userID:    "invalid",
			taskID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserInvalidID,
		},
		{
			name:      "Пользователь не существует",
			userID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID:    "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup:     func(m *mock.MockRepository) {},
			wantError: model.ErrUserNotExist,
		},
		{
			name:   "Невалидный ID задачи",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "invalid",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskInvalidID,
		},
		{
			name:   "Задача не существует",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
			},
			wantError: model.ErrTaskNotExist,
		},
		{
			name:   "Удаление чужой задачи",
			userID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			taskID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa",
			setup: func(m *mock.MockRepository) {
				m.CreateUser(ctx, &model.User{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", Name: "john"})
				m.CreateTask(ctx, &model.Task{ID: "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa", UserID: "cccccccc-cccc-cccc-cccc-cccccccccccc", Title: "Another Task"})
			},
			wantError: model.ErrUserForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := mock.New()
			tt.setup(mockRepo)

			service := NewTask(mockRepo, mockRepo)
			err := service.DeleteTask(ctx, tt.userID, tt.taskID)

			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}
		})
	}
}
