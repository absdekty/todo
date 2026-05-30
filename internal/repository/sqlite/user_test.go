package sqlite

import (
	"context"
	"errors"
	"testing"
	"todo/internal/model"
)

func TestCreateUser(t *testing.T) {
	db := setupUserDB(t)

	user, err := model.NewUser("...", "123456789")
	if err != nil {
		t.Skipf("new user: %v", err)
	}

	tests := []struct {
		name    string
		user    *model.User
		wantErr bool
	}{
		{
			name:    "обычный пользователь",
			user:    user,
			wantErr: false,
		},
		{
			name:    "Тот-же name",
			user:    user,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.CreateUser(context.Background(), user); err != nil && !tt.wantErr {
				t.Errorf("not expected error, but got: %v", err)
			}
		})
	}
}

func TestFindByName(t *testing.T) {
	if _, err := setupUserDB(t).FindByName(context.Background(), ""); err != nil {
		if !errors.Is(err, model.ErrUserNotExist) {
			t.Errorf("not expected error, but got: %v", err)
		}
	}
}
