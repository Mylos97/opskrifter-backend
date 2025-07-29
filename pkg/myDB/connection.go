package myDB

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var DB *sqlx.DB

func Init(inMemory bool) error {
	var err error
	dsn := "file:app.db?mode=rwc&_fk=on&_journal_mode=WAL"
	if inMemory {
		dsn = "file:sharedmemdb?mode=memory&cache=shared&_fk=on"
	}

	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Current working directory:", cwd)

	DB, err = sqlx.Open("sqlite3", dsn)
	if err != nil {
		return err
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping DB: %w", err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	packageRoot, err := findProjectRoot()

	if err != nil {
		return err
	}

	migrationsDir := filepath.Join(packageRoot, "migrations")
	if _, err := os.Stat(migrationsDir); os.IsNotExist(err) {
		log.Fatalf("migrations directory does not exist: %s", migrationsDir)
	}

	if err := goose.Up(DB.DB, migrationsDir); err != nil {
		log.Fatalf("failed to run schema migrations: %v", err)
	}

	err = checkDBSettings(inMemory, cwd)

	if err != nil {
		return err
	}

	return nil
}
