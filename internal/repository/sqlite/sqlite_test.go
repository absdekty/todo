package sqlite

import (
	"database/sql"
	"testing"
)

func setupDB(t *testing.T) *sql.DB {
	t.Helper()

	db, err := New(":memory:")
	if err != nil {
		t.Fatalf("ошибка инициализации БД: %v", err)
	}

	return db
}

func setupTaskDB(t *testing.T) *RepositoryTask {
	t.Helper()
	return NewTask(setupDB(t))
}

func setupUserDB(t *testing.T) *RepositoryUser {
	t.Helper()
	return NewUser(setupDB(t))
}
