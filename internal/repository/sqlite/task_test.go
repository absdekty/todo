package sqlite

import (
	"context"
	"errors"
	"testing"
	"todo/internal/model"
)

func TestCreateTask(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		task      *model.Task
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name: "Валидный таск",
			task: &model.Task{ID: "taskid", UserID: "userid"},
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
			},
			wantError: nil,
		},
		{
			name: "Идентичный таск",
			task: &model.Task{ID: "taskid", UserID: "userid"},
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid", UserID: "userid"})
			},
			wantError: model.ErrTaskAlreadyExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			err := repo.task.CreateTask(ctx, tt.task)
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
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name       string
		userid     string
		setup      func(repo *dbRepo)
		wantError  error
		countTasks int
	}{
		{
			name:   "У пользователя 2 задачи",
			userid: "userid",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid1", UserID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid2", UserID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid", UserID: "userid1"})
			},
			wantError:  nil,
			countTasks: 2,
		},
		{
			name:   "У пользователя нет задач",
			userid: "userid",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid", UserID: "userid1"})
			},
			wantError:  nil,
			countTasks: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			tasks, err := repo.task.GetUserTasks(ctx, tt.userid)
			if tt.wantError != nil {
				if !errors.Is(err, tt.wantError) {
					t.Errorf("expected %v, got %v", tt.wantError, err)
				}
			} else {
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}
			}

			if len(tasks) != tt.countTasks {
				t.Errorf("expected task count=%d, got %d", tt.countTasks, len(tasks))
			}
		})
	}
}

func TestGetTaskByID(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		taskid    string
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:   "Существует",
			taskid: "taskid",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid", UserID: "userid"})
			},
			wantError: nil,
		},
		{
			name:      "Не существует",
			taskid:    "taskid",
			setup:     func(repo *dbRepo) {},
			wantError: model.ErrTaskNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			_, err := repo.task.GetTaskByID(ctx, tt.taskid)
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
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name        string
		updatedTask *model.Task
		setup       func(repo *dbRepo)
		wantError   error
	}{
		{
			name:        "Существует",
			updatedTask: &model.Task{ID: "taskid", Title: "new title"},
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid", UserID: "userid"})
			},
			wantError: nil,
		},
		{
			name:        "Не существует",
			updatedTask: &model.Task{ID: "taskid", Title: "new title"},
			setup:       func(repo *dbRepo) {},
			wantError:   model.ErrTaskNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			err := repo.task.UpdateTask(ctx, tt.updatedTask)
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
	t.Parallel()

	ctx := context.Background()

	tests := []struct {
		name      string
		taskid    string
		setup     func(repo *dbRepo)
		wantError error
	}{
		{
			name:   "Существует",
			taskid: "taskid",
			setup: func(repo *dbRepo) {
				repo.user.CreateUser(context.Background(), &model.User{ID: "userid"})
				repo.task.CreateTask(context.Background(), &model.Task{ID: "taskid", UserID: "userid"})
			},
			wantError: nil,
		},
		{
			name:      "Не существует",
			taskid:    "taskid",
			setup:     func(repo *dbRepo) {},
			wantError: model.ErrTaskNotExist,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := setupDB(t)
			tt.setup(repo)

			err := repo.task.DeleteTask(ctx, tt.taskid)
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
