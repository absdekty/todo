package mock

import (
	"todo/internal/model"
)

type MockRepository struct {
	task       map[string]*model.Task
	userByID   map[string]*model.User
	userByName map[string]*model.User
	tokens     map[string]*model.Token
}

func New() *MockRepository {
	return &MockRepository{
		task:       make(map[string]*model.Task),
		userByID:   make(map[string]*model.User),
		userByName: make(map[string]*model.User),
		tokens:     make(map[string]*model.Token),
	}
}
