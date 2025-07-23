// Package storage is a nice package
package storage

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
)

type Storage struct {
	db *sql.DB
}

type DeleteOptions struct {
	URL string
	Alias string
}

func NewSQLite(storagePath string) (*Storage, error) {
	const operation = "storage.sqlite.New"

	db, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", operation, err)
	}

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS url(
		id INTEGER PRIMARY KEY,
		alias TEXT NOT NULL UNIQUE,
		url TEXT NOT NULL);
	CREATE INDEX IF NOT EXISTS index_alias ON url(alias);
		`)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", operation, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s, %w", operation, err)
	}

	/* _, err = db.Exec(`CREATE TABLE IF NOT EXISTS users(
	userID INTEGER PRIMARY KEY,
	username TEXT NOT NULL,
	password TEXT NOT NULL);`)
	if err != nil {
		return nil, fmt.Errorf("%s, %w", operation, err)
	} */

	return &Storage{db: db}, nil
}

func (s *Storage) SaveURL(urlToSave, alias string) (int64, error) {
	const operation = "storage.sqlite.GetURL"

	stmt, err := s.db.Prepare("INSERT INTO url(url, alias) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s, %w", operation, err)
	}
	res, err := stmt.Exec(urlToSave, alias)
	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s, %w", operation, ErrURLExist)
		}
		return 0, fmt.Errorf("%s, %w", operation, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failes to get last insert id: %w", operation, err)
	}

	return id, nil
}

// GetURL получает url по его алиасу
func (s *Storage) GetURL(alias string) (string, error) {
	const operation = "storage.sqlite.GetURL"
	var resURL string
	
	err := s.db.QueryRow("SELECT url FROM url WHERE alias = ?", alias).Scan(&resURL)
	if errors.Is(err, sql.ErrNoRows) {
    	return "", ErrURLNotFound
	}
	if err != nil {
    	return "", fmt.Errorf("%s: %w", operation, err)
	}
	return resURL, nil
}

func (s *Storage) DeleteURL(opts DeleteOptions) error {
	const operation = "storage.sqlite.DeleteUrl"
	
	if opts.Alias == "" && opts.URL == "" {
		return fmt.Errorf("%s: required either alias or URL", operation)
	}
	
	
	if opts.URL != "" {
		res, err := s.db.Exec("DELETE FROM url WHERE url = ?", opts.URL)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrURLNotFound
		}
		if err != nil {
			return fmt.Errorf("%s: %w", operation, err)
		}
		if rows, _ := res.RowsAffected(); rows > 0 {
			return nil
		}
	}
	if opts.Alias != "" {
		res, err := s.db.Exec("DELETE FROM url WHERE alias = ?", opts.Alias)
		if errors.Is(err, sql.ErrNoRows) {
			return ErrAliasNotFound
		}
		if err != nil {
			return fmt.Errorf("%s: %w", operation, err)
		}
		if rows, _ := res.RowsAffected(); rows > 0 {
			return nil
		}
	}
	return ErrURLNotFound
}

func (s *Storage) IsAliasExists(alias string) (bool, error) {
	const operation = "storage.sqlite.IsAliasExist"

	var exists bool
	const query = "SELECT"
	
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM url WHERE alias = ?)", alias).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s, %w", operation, err)
	}

	return exists, nil
}


