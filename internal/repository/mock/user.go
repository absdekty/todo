package mock

import (
	"context"
	"todo/internal/model"
)

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
