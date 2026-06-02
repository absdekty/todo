package sqlite

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

func New(dbPath string) (*sql.DB, error) {
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

	return db, nil
}

func InitSchema(db *sql.DB) error {
	query := []string{
		`PRAGMA foreign_keys = ON`,

		`CREATE TABLE IF NOT EXISTS users (
					id TEXT PRIMARY KEY,
					name TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL,
					role TEXT NOT NULL,
					created DATETIME NOT NULL
				);`,

		`CREATE TABLE IF NOT EXISTS tasks (
					id TEXT NOT NULL,
					userid TEXT NOT NULL,
					title TEXT NOT NULL,
					description TEXT NOT NULL,
					completed INTEGER NOT NULL,
					created DATETIME NOT NULL,
					updated DATETIME NOT NULL,
					
					FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE,
					UNIQUE(userid, id)
				);`,

		`CREATE TABLE IF NOT EXISTS tokens (
					userid TEXT NOT NULL,
					token TEXT PRIMARY KEY,
					revoked INTEGER NOT NULL DEFAULT 0,
					expires DATETIME NOT NULL,
					created DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
					
					FOREIGN KEY (userid) REFERENCES users(id) ON DELETE CASCADE
				);`,

		`CREATE INDEX IF NOT EXISTS idx_tasks_userid ON tasks(userid);`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_userid ON tokens(userid);`,
		`CREATE INDEX IF NOT EXISTS idx_tokens_expires ON tokens(expires);`,
	}

	for _, q := range query {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}

func Close(db *sql.DB) error {
	return db.Close()
}
