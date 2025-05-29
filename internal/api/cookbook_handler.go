package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// POST /cookbooks/{id}
func CreateCookBook(w http.ResponseWriter, r *http.Request) {
	var cb types.CookBook
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	cb.ID = uuid.New().String()

	_, err := db.DB.Exec(`
		INSERT INTO cookbooks (id, name, description, likes, user)
		VALUES (?, ?, ?, ?, ?)`,
		cb.ID, cb.Name, cb.Description, cb.Likes, cb.User)

	if err != nil {
		http.Error(w, "Failed to insert cookbook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, cb)
}

// GET /cookbooks/{id}
func GetCookBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	row := db.DB.QueryRow(`
		SELECT id, name, description, likes, user
		FROM cookbooks WHERE id = ?`, id)

	var cb types.CookBook
	err := row.Scan(&cb.ID, &cb.Name, &cb.Description, &cb.Likes, &cb.User)
	if err == sql.ErrNoRows {
		http.Error(w, "Cookbook not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, cb)
}

// PUT /cookbooks/{id}
func UpdateCookBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var cb types.CookBook
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if cb.ID != "" && cb.ID != id {
		http.Error(w, "ID in URL and body do not match", http.StatusBadRequest)
		return
	}
	cb.ID = id

	_, err := db.DB.Exec(`
		UPDATE cookbooks SET name=?, description=?, likes=?, user=?
		WHERE id=?`,
		cb.Name, cb.Description, cb.Likes, cb.User, cb.ID)

	if err != nil {
		http.Error(w, "Failed to update cookbook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, cb)
}

// DELETE /cookbooks/{id}
func DeleteCookBook(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := db.DB.Exec(`DELETE FROM cookbooks WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete cookbook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
