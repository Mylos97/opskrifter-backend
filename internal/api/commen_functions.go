package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"opskrifter-backend/pkg/db"
)

func BeginTx(w http.ResponseWriter) (*sql.Tx, error) {
	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

func CommitTx(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	return nil
}
