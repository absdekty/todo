package repository

import (
	"context"
	"todo/internal/model"
)

type RepositoryTask interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTasks(ctx context.Context, userID string) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, userID, taskID string) error
}

type RepositoryUser interface {
	CreateUser(ctx context.Context, user *model.User) error
	FindByName(ctx context.Context, name string) (*model.User, error)
	FindByID(ctx context.Context, userID string) (*model.User, error)
}
