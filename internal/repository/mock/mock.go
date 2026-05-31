package mock

import (
	"context"
	"todo/internal/model"
)

type MockRepository struct {
	task       map[string]*model.Task
	userByID   map[string]*model.User
	userByName map[string]*model.User
}

func New() *MockRepository {
	return &MockRepository{
		task:       make(map[string]*model.Task),
		userByID:   make(map[string]*model.User),
		userByName: make(map[string]*model.User),
	}
}

// Task methods
func (m *MockRepository) CreateTask(ctx context.Context, task *model.Task) error {
	m.task[task.ID] = task
	return nil
}

func (m *MockRepository) GetUserTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	var tasks []*model.Task
	for _, t := range m.task {
		if t.UserID == userID {
			tasks = append(tasks, t)
		}
	}
	return tasks, nil
}

func (m *MockRepository) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	task, ok := m.task[taskID]
	if !ok {
		return nil, model.ErrTaskNotExist
	}
	return task, nil
}

func (m *MockRepository) UpdateTask(ctx context.Context, task *model.Task) error {
	if _, ok := m.task[task.ID]; !ok {
		return model.ErrTaskNotExist
	}
	m.task[task.ID] = task
	return nil
}

func (m *MockRepository) DeleteTask(ctx context.Context, taskID string) error {
	delete(m.task, taskID)
	return nil
}

// User methods
func (m *MockRepository) FindByName(ctx context.Context, name string) (*model.User, error) {
	user, ok := m.userByName[name]
	if !ok {
		return nil, model.ErrUserNotExist
	}
	return user, nil
}

func (m *MockRepository) FindByID(ctx context.Context, userID string) (*model.User, error) {
	user, ok := m.userByID[userID]
	if !ok {
		return nil, model.ErrUserNotExist
	}
	return user, nil
}

func (m *MockRepository) CreateUser(ctx context.Context, user *model.User) error {
	m.userByID[user.ID] = user
	m.userByName[user.Name] = user
	return nil
}
