package sqlite

import (
	"context"
	"errors"
	"testing"
	"todo/internal/model"
)

func TestCreateTask(t *testing.T) {
	db := setupTaskDB(t)

	task, err := model.NewTask("550e8400-e29b-41d4-a716-446655440000", "...", "")
	if err != nil {
		t.Skipf("new task: %v", err)
	}

	tests := []struct {
		name    string
		task    *model.Task
		wantErr bool
	}{
		{
			name:    "обычная задача",
			task:    task,
			wantErr: false,
		},
		{
			name:    "ID прошлой задачи",
			task:    task,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.CreateTask(context.Background(), task); err != nil && !tt.wantErr {
				t.Errorf("not expected error, but got: %v", err)
			}
		})
	}
}

func TestGetTasks(t *testing.T) {
	_, err := setupTaskDB(t).GetTasks(context.Background(), "")
	if err != nil {
		t.Errorf("not expected error, but got: %v", err)
	}
}

func TestGetTaskByID(t *testing.T) {
	if _, err := setupTaskDB(t).GetTaskByID(context.Background(), "", ""); err != nil {
		if !errors.Is(err, model.ErrTaskNotExist) {
			t.Errorf("not expected error, but got: %v", err)
		}
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupTaskDB(t)

	task1, err1 := model.NewTask("550e8400-e29b-41d4-a716-446655440000", "...", "")
	task2, err2 := model.NewTask("550e8400-e29b-41d4-a716-446655440000", "...", "")
	if err1 != nil {
		t.Skipf("ошибка создания задачи: %v", err1)
	}
	if err2 != nil {
		t.Skipf("ошибка создания задачи: %v", err2)
	}
	if err := db.CreateTask(context.Background(), task1); err != nil {
		t.Skipf("ошибка создания задачи(db): %v", err)
	}

	tests := []struct {
		name    string
		task    *model.Task
		wantErr bool
	}{
		{
			name:    "задача существует",
			task:    task1,
			wantErr: false,
		},
		{
			name:    "задача не существует",
			task:    task2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.UpdateTask(context.Background(), tt.task); err != nil && !tt.wantErr {
				t.Errorf("not expected error, but got: %v", err)
			}
		})
	}
}

func TestDeleteTask(t *testing.T) {
	db := setupTaskDB(t)

	task, err := model.NewTask("550e8400-e29b-41d4-a716-446655440000", "...", "")
	if err != nil {
		t.Skipf("ошибка создания задачи: %v", err)
	}
	if err := db.CreateTask(context.Background(), task); err != nil {
		t.Skipf("ошибка создания задачи(db): %v", err)
	}

	tests := []struct {
		name    string
		taskID  string
		wantErr bool
	}{
		{
			name:    "задача существует",
			taskID:  task.ID,
			wantErr: false,
		},
		{
			name:    "задача не существует",
			taskID:  "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.DeleteTask(context.Background(), "550e8400-e29b-41d4-a716-446655440000", tt.taskID); err != nil && !tt.wantErr {
				t.Errorf("not expected error, but got: %v", err)
			}
		})
	}
}
