package sqlite

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/empfaze/golang_url_reducer/internal/storage"
	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

const OPERATION_TRACE_NEW = "internal.storage.sqlite.New"
const OPERATION_TRACE_SAVE_URL = "internal.storage.sqlite.SaveURL"
const OPERATION_TRACE_GET_URL = "internal.storage.sqlite.GetURL"

func New(storagePath string) (*Storage, error) {
	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", OPERATION_TRACE_NEW, err)
	}

	query, err := db.Prepare(`
		CREATE TABLE IF NOT EXISTS urls(
			id INTEGER PRIMARY KEY,
			url TEXT NOT NULL,
			alias TEXT NOT NULL UNIQUE);
		CREATE INDEX IF NOT EXISTS idx_alias ON url(alias);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", OPERATION_TRACE_NEW, err)
	}

	if _, err := query.Exec(); err != nil {
		return nil, fmt.Errorf("%s: %w", OPERATION_TRACE_NEW, err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	query, err := s.db.Prepare("INSERT INTO urls(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_SAVE_URL, err)
	}

	result, err := query.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintCheck {
			return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_SAVE_URL, storage.ErrURLExists)
		}

		return 0, fmt.Errorf("%s: %w", OPERATION_TRACE_NEW, err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", "Failed to get last insert id", err)
	}

	return id, nil
}

func (s *Storage) GetURL(alias string) (string, error) {
	query, err := s.db.Prepare("SELECT url FROM urls WHERE alias = ?")
	if err != nil {
		return "", fmt.Errorf("%s: %w", OPERATION_TRACE_GET_URL, err)
	}

	var resultURL string

	err = query.QueryRow(alias).Scan(&resultURL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", storage.ErrURLNotFound
		}

		return "", fmt.Errorf("%s: %w", OPERATION_TRACE_GET_URL, err)
	}

	return resultURL, nil
}
