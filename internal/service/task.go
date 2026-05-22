package service

import (
	"context"
	"fmt"
	"todo/internal/model"
	"todo/internal/repository"
)

type ServiceTaskI interface {
	CreateTask(ctx context.Context, userID, title, description string) (*model.Task, error)
	GetTasks(ctx context.Context, userID string) ([]*model.Task, error)
	GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error)
	UpdateTask(ctx context.Context, task *model.Task) error
	DeleteTask(ctx context.Context, userID, taskID string) error
}

var _ ServiceTaskI = (*ServiceTask)(nil)

type ServiceTask struct {
	repo     repository.RepositoryTask
	userRepo repository.RepositoryUser
}

func NewTask(repo repository.RepositoryTask, userRepo repository.RepositoryUser) *ServiceTask {
	return &ServiceTask{repo: repo, userRepo: userRepo}
}

func (s *ServiceTask) CreateTask(ctx context.Context, userID, title, description string) (*model.Task, error) {
	if err := model.UserValidateID(userID); err != nil {
		return nil, fmt.Errorf("create task (validate user): %w", err)
	}

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("create task (repo user): %w", err)
	}

	task, err := model.NewTask(userID, title, description)
	if err != nil {
		return nil, fmt.Errorf("create task (model): %w", err)
	}

	if err := s.repo.CreateTask(ctx, task); err != nil {
		return nil, fmt.Errorf("create task (repo): %w", err)
	}

	return task, nil
}

func (s *ServiceTask) GetTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	if err := model.UserValidateID(userID); err != nil {
		return nil, fmt.Errorf("get tasks (validate user): %w", err)
	}

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("get tasks (repo user): %w", err)
	}

	return s.repo.GetTasks(ctx, userID)
}

func (s *ServiceTask) GetTaskByID(ctx context.Context, userID, taskID string) (*model.Task, error) {
	if err := model.UserValidateID(userID); err != nil {
		return nil, fmt.Errorf("get task by id (validate user): %w", err)
	}

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return nil, fmt.Errorf("get task by id (repo user): %w", err)
	}

	if err := model.TaskValidateID(taskID); err != nil {
		return nil, fmt.Errorf("get task by id (validate task): %w", err)
	}

	return s.repo.GetTaskByID(ctx, userID, taskID)
}

func (s *ServiceTask) UpdateTask(ctx context.Context, task *model.Task) error {
	if task == nil {
		return model.ErrTaskIsNil
	}

	if err := model.UserValidateID(task.UserID); err != nil {
		return fmt.Errorf("update task (validate user): %w", err)
	}

	if _, err := s.userRepo.FindByID(ctx, task.UserID); err != nil {
		return fmt.Errorf("update task (repo user): %w", err)
	}

	return s.repo.UpdateTask(ctx, task)
}

func (s *ServiceTask) DeleteTask(ctx context.Context, userID, taskID string) error {
	if err := model.UserValidateID(userID); err != nil {
		return fmt.Errorf("delete task (validate user): %w", err)
	}

	if _, err := s.userRepo.FindByID(ctx, userID); err != nil {
		return fmt.Errorf("delete task (repo user): %w", err)
	}

	if err := model.TaskValidateID(taskID); err != nil {
		return fmt.Errorf("delete task (validate task): %w", err)
	}

	return s.repo.DeleteTask(ctx, userID, taskID)
}
