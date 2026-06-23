package db

import (
	"database/sql"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func Open(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath)
}

func Init(dbPath string) (*sql.DB, error) {
	db, err := Open(dbPath)
	if err != nil {
		return nil, err
	}

	if err := createTables(db); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func createTables(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS password_records (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			archive_path TEXT,
			archive_name TEXT,
			archive_hash TEXT,
			encrypted_password BLOB NOT NULL,
			success_count INTEGER DEFAULT 1,
			last_used_at DATETIME,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS extract_history (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			archive_path TEXT NOT NULL,
			output_dir TEXT NOT NULL,
			success INTEGER NOT NULL,
			used_password INTEGER DEFAULT 0,
			error_message TEXT,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS app_config (
			key TEXT PRIMARY KEY,
			value TEXT,
			updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return err
		}
	}

	return nil
}

func DefaultDBPath(configDir string) string {
	return filepath.Join(configDir, "pocketzip.db")
}
