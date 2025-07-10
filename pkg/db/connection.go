package db

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var DB *sqlx.DB

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

func Init(inMemory bool) {
	var err error
	dsn := "./app.db"
	if inMemory {
		dsn = ":memory:"
	}

	DB, err = sqlx.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	packageRoot, err := findProjectRoot()

	if err != nil {
		log.Fatal(err)
	}

	migrationsDir := filepath.Join(packageRoot, "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("migrations directory does not exist: %s", migrationsDir)
	}

	schemaDir := filepath.Join(packageRoot, "migrations", "schema")
	if err := goose.Up(DB.DB, schemaDir); err != nil {
		log.Fatalf("failed to run schema migrations: %v", err)
	}
}
