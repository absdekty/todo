package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"todo/internal/model"
)

type RepositoryTask struct {
	*sql.DB
}

func NewTask(db *sql.DB) *RepositoryTask {
	return &RepositoryTask{db}
}

func (r *RepositoryTask) CreateTask(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks
		(id, userid, title, description, completed, created, updated)
		VALUES (?, ?, ?, ?, ?, ?, ?)`

	result, err := r.ExecContext(ctx, query,
		task.ID, task.UserID, task.Title, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("task was not inserted")
	}

	return nil
}

func (r *RepositoryTask) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	query := `
		SELECT
		id, userid, title, description, completed, created, updated
		FROM tasks
		WHERE id=?`

	var task model.Task
	err := r.QueryRowContext(ctx, query, taskID).Scan(
		&task.ID, &task.UserID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrTaskNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (r *RepositoryTask) GetUserTasks(ctx context.Context, userID string) ([]*model.Task, error) {
	query := `
		SELECT
		id, userid, title, description, completed, created, updated
		FROM tasks WHERE userid=?`

	rows, err := r.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}

		err := rows.Scan(&task.ID, &task.UserID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		tasks = append(tasks, task)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *RepositoryTask) UpdateTask(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE tasks SET
		title=?, description=?, completed=?, updated=?
		WHERE id=?`

	result, err := r.ExecContext(ctx, query, task.Title, task.Description, task.Completed, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return model.ErrTaskNotExist
	}

	return nil
}

func (r *RepositoryTask) DeleteTask(ctx context.Context, taskID string) error {
	query := `
		DELETE FROM tasks
		WHERE id=?`

	result, err := r.ExecContext(ctx, query, taskID)
	if err != nil {
		return err
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return model.ErrTaskNotExist
	}

	return nil
}
