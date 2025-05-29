package api

import (
	"opskrifter-backend/pkg/db"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

func TestMain(m *testing.M) {
	db.Init(true)
	code := m.Run()

	db.DB.Close()
	os.Exit(code)
}
