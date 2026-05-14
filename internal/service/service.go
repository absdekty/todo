package service

import (
	"context"
	"fmt"
	"todo/internal/model"
	"todo/internal/repository"
)

type ServiceI interface {
	CreateTask(ctx context.Context, title, description string) (*model.Task, error)
	GetTasks(ctx context.Context) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, taskID string) error
}

var _ ServiceI = (*Service)(nil)

type Service struct {
	repo repository.RepositoryI
}

func New(repo repository.RepositoryI) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(ctx context.Context, title, description string) (*model.Task, error) {
	task, err := model.NewTask(title, description)
	if err != nil {
		return nil, fmt.Errorf("create task (model): %w", err)
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("create task (repo): %w", err)
	}

	return task, nil
}

func (s *Service) GetTasks(ctx context.Context) ([]*model.Task, error) {
	return s.repo.GetTasks(ctx)
}

func (s *Service) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	if err := model.ValidateID(taskID); err != nil {
		return nil, fmt.Errorf("get task by id (validate): %w", err)
	}

	return s.repo.GetTaskByID(ctx, taskID)
}

func (s *Service) UpdateTask(ctx context.Context, task *model.Task) error {
	if task == nil {
		return model.ErrTaskIsNil
	}

	return s.repo.UpdateTask(ctx, task)
}

func (s *Service) DeleteTask(ctx context.Context, taskID string) error {
	if err := model.ValidateID(taskID); err != nil {
		return fmt.Errorf("delete task (validate): %w", err)
	}

	return s.repo.DeleteTask(ctx, taskID)
}
