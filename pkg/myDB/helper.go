package myDB

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}

func checkDBSettings(inMemory bool, cwd string) error {
	var fkEnabled int
	err := DB.QueryRow("PRAGMA foreign_keys;").Scan(&fkEnabled)
	if err != nil {
		return fmt.Errorf("failed to query foreign_keys: %w", err)
	}

	if fkEnabled != 1 {
		return fmt.Errorf("expected foreign_keys to be 1, got %d", fkEnabled)
	}

	var journalMode string
	err = DB.QueryRow("PRAGMA journal_mode;").Scan(&journalMode)
	if err != nil {
		return fmt.Errorf("failed to query journal_mode: %w", err)
	}

	if inMemory && journalMode != "memory" {
		return fmt.Errorf("expected journal_mode to be memory, got %s", journalMode)
	}

	if !inMemory && strings.ToLower(journalMode) != "wal" {
		return fmt.Errorf("expected journal_mode to be WAL, got %s", journalMode)
	}

	if !inMemory {
		if _, err := os.Stat("./app.db"); os.IsNotExist(err) {
			return fmt.Errorf("database file was not created at %s/app.db", cwd)
		}
		fmt.Println("Database file created at:", path.Join(cwd, "app.db"))
	}

	return nil
}
