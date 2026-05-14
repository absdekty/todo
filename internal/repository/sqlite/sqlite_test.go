package repository

import (
	"context"
	"testing"
	"todo/internal/model"
)

func setupDB(t *testing.T) *Repository {
	t.Helper()

	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("ошибка инициализации БД: %v", err)
	}

	return db
}

func TestCreateTask(t *testing.T) {
	db := setupDB(t)

	task, err := model.NewTask("...", "")
	if err != nil {
		t.Skipf("ошибка создания задачи: %v", err)
	}

	tests := []struct {
		name    string
		task    *model.Task
		wantErr bool
	}{
		{
			name:    "обычное создание",
			task:    task,
			wantErr: false,
		},
		{
			name:    "повторное создание(Тот-же ID)",
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
	_, err := setupDB(t).GetTasks(context.Background())
	if err != nil {
		t.Errorf("not expected error, but got: %v", err)
	}
}

func TestGetTaskByID(t *testing.T) {
	if _, err := setupDB(t).GetTaskByID(context.Background(), ""); err != nil {
		t.Errorf("not expected error, but got: %v", err)
	}
}

func TestUpdateTask(t *testing.T) {
	db := setupDB(t)

	task1, err1 := model.NewTask("...", "")
	task2, err2 := model.NewTask("...", "")
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
	db := setupDB(t)

	task, err := model.NewTask("...", "")
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
			if err := db.DeleteTask(context.Background(), tt.taskID); err != nil && !tt.wantErr {
				t.Errorf("not expected error, but got: %v", err)
			}
		})
	}
}
