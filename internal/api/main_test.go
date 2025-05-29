package api

import (
	"log"
	"os"
	"testing"
	_ "github.com/mattn/go-sqlite3"
	"opskrifter-backend/pkg/db"
)

func TestMain(m *testing.M) {
    var err error
    db.Init()
    
		if err != nil {
        log.Fatalf("failed to initialize test DB: %v", err)
    }

    code := m.Run()

    db.DB.Close()
    os.Exit(code)
}
