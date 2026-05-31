package rest

import (
	"context"
	"todo/internal/model"
)

type ServiceTask interface {
	CreateTask(ctx context.Context, userID, title, description string) (*model.Task, error)
	GetUserTasks(ctx context.Context, userID string) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, userID, taskID string) error
}

type ServiceUser interface {
	CreateUser(ctx context.Context, name, password string) error
	FindByName(ctx context.Context, name string) (*model.User, error)
	Login(ctx context.Context, name, password string) (string, error)
}

type ServiceToken interface {
	GenerateToken(userID string) (string, error)
	ValidateToken(token string) (string, error)
}
