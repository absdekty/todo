package repository

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"todo/internal/model"
	"todo/internal/repository"
)

var _ repository.RepositoryI = (*Repository)(nil)

type Repository struct {
	*sql.DB
}

func New(dbPath string) (*Repository, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if err := InitSchema(db); err != nil {
		return nil, err
	}

	return &Repository{db}, nil
}

func InitSchema(db *sql.DB) error {
	query := `CREATE TABLE IF NOT EXISTS tasks (
			id TEXT PRIMARY KEY,
			title TEXT NOT NULL,
			description TEXT NOT NULL,
			completed INTEGER NOT NULL,
			created DATETIME NOT NULL,
			updated DATETIME NOT NULL);`

	_, err := db.Exec(query)
	return err
}

func (r *Repository) Close() error {
	return r.DB.Close()
}

func (r *Repository) CreateTask(ctx context.Context, task *model.Task) error {
	query := `
		INSERT INTO tasks
		(id, title, description, completed, created, updated)
		VALUES (?, ?, ?, ?, ?, ?)`

	result, err := r.ExecContext(ctx, query,
		task.ID, task.Title, task.Description, task.Completed, task.CreatedAt, task.UpdatedAt)
	if err != nil {
		return fmt.Errorf("failed to create task: %w", err)
	}

	if rows, _ := result.RowsAffected(); rows == 0 {
		return fmt.Errorf("task was not inserted")
	}

	return nil
}

func (r *Repository) GetTaskByID(ctx context.Context, taskID string) (*model.Task, error) {
	query := `
		SELECT
		title, description, completed, created, updated
		FROM tasks
		WHERE id=?`

	var task model.Task
	err := r.QueryRowContext(ctx, query, taskID).Scan(
		&task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, model.ErrTaskNotExist
	}
	if err != nil {
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	task.ID = taskID

	return &task, nil
}

func (r *Repository) GetTasks(ctx context.Context) ([]*model.Task, error) {
	query := `
		SELECT
		id, title, description, completed, created, updated
		FROM tasks`

	rows, err := r.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get tasks: %w", err)
	}
	defer rows.Close()

	var tasks []*model.Task
	for rows.Next() {
		task := &model.Task{}

		err := rows.Scan(&task.ID, &task.Title, &task.Description, &task.Completed, &task.CreatedAt, &task.UpdatedAt)
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

func (r *Repository) UpdateTask(ctx context.Context, task *model.Task) error {
	query := `
		UPDATE tasks SET
		title=?, description=?, completed=?, updated=?
		WHERE id=?`

	result, err := r.ExecContext(ctx, query, task.Title, task.Description, task.Completed, task.UpdatedAt, task.ID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return model.ErrTaskNotExist
	}

	return nil
}

func (r *Repository) DeleteTask(ctx context.Context, taskID string) error {
	query := `
		DELETE FROM tasks
		WHERE id=?`

	result, err := r.ExecContext(ctx, query, taskID)
	if err != nil {
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected == 0 {
		return model.ErrTaskNotExist
	}

	return nil
}
