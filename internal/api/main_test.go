package api

import (
	"fmt"
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
	var fkEnabled bool
	err = myDB.DB.QueryRow("PRAGMA foreign_keys;").Scan(&fkEnabled)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Foreign keys enabled: %v\n", fkEnabled)

	myDB.DB.Close()
	os.Exit(code)
}
