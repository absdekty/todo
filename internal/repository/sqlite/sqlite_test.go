package sqlite

import (
	"testing"
)

type dbRepo struct {
	user  *RepositoryUser
	task  *RepositoryTask
	token *RepositoryToken
}

func setupDB(t *testing.T) *dbRepo {
	t.Helper()

	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("ошибка инициализации БД: %v", err)
	}

	return &dbRepo{
		user:  NewUser(db),
		task:  NewTask(db),
		token: NewToken(db),
	}
}
