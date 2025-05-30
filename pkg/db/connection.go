package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose"
)

var DB *sql.DB

func Init(inMemory bool) {
	var err error
	dsn := "./app.db"
	if inMemory {
		dsn = ":memory:"
	}

	DB, err = sql.Open("sqlite3", dsn)
	if err != nil {
		log.Fatal(err)
	}

	if err := goose.SetDialect("sqlite3"); err != nil {
		log.Fatal(err)
	}

	if err := goose.Up(DB, "migrations"); err != nil {
		log.Fatal(err)
	}
}
