package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"github.com/go-chi/chi/v5"
	"opskrifter-backend/pkg/db"
	"opskrifter-backend/internal/types"
	"github.com/google/uuid"
)

// Helper: write JSON response
func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// CreateRecipe handles POST /recipes
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var rec types.Recipe
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	rec.ID = uuid.New().String()

	_, err := db.DB.Exec(`
		INSERT INTO recipes (id, name, minutes, rating, description, likes, comments, image)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.Name, rec.Minutes, rec.Rating, rec.Description, rec.Likes, rec.Comments, rec.Image)

	if err != nil {
		http.Error(w, "Failed to insert recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, rec)
}

// GetRecipe handles GET /recipes/{id}
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	row := db.DB.QueryRow(`
		SELECT id, name, minutes, rating, description, likes, comments, image
		FROM recipes WHERE id = ?`, id)

	var rec types.Recipe
	err := row.Scan(&rec.ID, &rec.Name, &rec.Minutes, &rec.Rating, &rec.Description, &rec.Likes, &rec.Comments, &rec.Image)
	if err == sql.ErrNoRows {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

// UpdateRecipe handles PUT /recipes/{id}
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var rec types.Recipe
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	// Make sure the ID in URL matches ID in JSON body (optional)
	if rec.ID != "" && rec.ID != id {
		http.Error(w, "ID in URL and body do not match", http.StatusBadRequest)
		return
	}

	rec.ID = id

	_, err := db.DB.Exec(`
		UPDATE recipes SET name=?, minutes=?, rating=?, description=?, likes=?, comments=?, image=?
		WHERE id=?`,
		rec.Name, rec.Minutes, rec.Rating, rec.Description, rec.Likes, rec.Comments, rec.Image, rec.ID)

	if err != nil {
		http.Error(w, "Failed to update recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

// DeleteRecipe handles DELETE /recipes/{id}
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := db.DB.Exec(`DELETE FROM recipes WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
