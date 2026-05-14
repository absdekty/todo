package repository

import (
	"context"
	"todo/internal/model"
)

type RepositoryI interface {
	CreateTask(ctx context.Context, task *model.Task) error
	GetTasks(ctx context.Context) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, taskID string) error
	Close() error
}
