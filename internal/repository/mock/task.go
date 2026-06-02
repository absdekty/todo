package mock

import (
	"context"
	"todo/internal/model"
)

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
