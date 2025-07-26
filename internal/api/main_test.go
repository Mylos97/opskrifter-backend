package api

import (
	"log"
	"opskrifter-backend/pkg/myDB"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	err := myDB.Init(true)
	if err != nil {
		log.Fatalf("error init DB %s", err)
	}
	code := m.Run()
	myDB.DB.Close()
	os.Exit(code)
}
