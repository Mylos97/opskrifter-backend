package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
	"time"

	"github.com/go-chi/chi/v5"
)

// POST /user/
func CreateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User
	user.CreatedAt = time.Now().String()
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO users (name, created_at)
		VALUES (?, ?)`,
		user.Name, user.CreatedAt)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

// GET /users/{id}
func GetUser(w http.ResponseWriter, r *http.Request) {
	ID := chi.URLParam(r, "id")

	row := db.DB.QueryRow(`
		SELECT id, name, created_at
		FROM users WHERE id = ?`, ID)

	var user types.User
	err := row.Scan(&user.ID, &user.Name, &user.CreatedAt)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// PUT /users/{id}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var user types.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if user.ID != "" && user.ID != id {
		http.Error(w, "ID in URL and body do not match", http.StatusBadRequest)
		return
	}
	user.ID = id

	_, err := db.DB.Exec(`
		UPDATE users SET name = ? WHERE id = ?`,
		user.Name, user.ID)

	if err != nil {
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

// DELETE /user/{id}
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	result, err := db.DB.Exec(`DELETE FROM users WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to get affected rows: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
