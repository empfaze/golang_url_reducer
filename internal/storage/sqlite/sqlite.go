package sqlite

import (
	"database/sql"
	"fmt"
)

type Storage struct {
	db *sql.DB
}

const OPERATION_TRACE = "storage.sqlite.New"

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", OPERATION_TRACE, err)
	}

	query, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls(
			id INTEGER PRIMARY KEY,
			alias TEXT NOT NULL UNIQUE,
			url TEXT NOT NULL);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("#{OPERATION_TRACE}: #{err}")
	}

	if _, err := query.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", OPERATION_TRACE, err)
	}

	return &Storage{db: db}, nil
}
