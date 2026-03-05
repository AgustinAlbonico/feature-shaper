package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

func DBPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("cannot get home dir: %w", err)
	}
	return filepath.Join(home, ".feature-shaper", "features.db"), nil
}

func Migrate() (*sql.DB, error) {
	dbPath, err := DBPath()
	if err != nil {
		return nil, err
	}

	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, fmt.Errorf("cannot create dir %s: %w", dir, err)
	}

	database, err := sql.Open("sqlite", dbPath+"?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)")
	if err != nil {
		return nil, fmt.Errorf("cannot open db: %w", err)
	}

	if _, err := database.Exec(SchemaSQL); err != nil {
		database.Close()
		return nil, fmt.Errorf("cannot execute schema: %w", err)
	}

	return database, nil
}
