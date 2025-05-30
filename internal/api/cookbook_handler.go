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

// POST /cookbooks/
func CreateCookBook(w http.ResponseWriter, r *http.Request) {
	var cb types.CookBook
	if err := json.NewDecoder(r.Body).Decode(&cb); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	cb.ID = uuid.New().String()

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO cookbooks (id, name, description, likes, user)
		VALUES (?, ?, ?, ?, ?)`,
		cb.ID, cb.Name, cb.Description, cb.Likes, cb.User)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert cookbook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, recipe := range cb.Recipes {
		if recipe.ID == "" {
			tx.Rollback()
			http.Error(w, "All recipes must have an ID to link to cookbook", http.StatusBadRequest)
			return
		}

		_, err = tx.Exec(`
			INSERT INTO cookbook_recipes (cookbook_id, recipe_id)
			VALUES (?, ?)`,
			cb.ID, recipe.ID)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to link recipe to cookbook: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
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

	rows, err := db.DB.Query(`
		SELECT r.id, r.name, r.minutes, r.description, r.likes, r.comments, r.image
		FROM recipes r
		INNER JOIN cookbook_recipes cr ON r.id = cr.recipe_id
		WHERE cr.cookbook_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to get recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var recipes []types.Recipe
	for rows.Next() {
		var r types.Recipe
		err := rows.Scan(&r.ID, &r.Name, &r.Minutes, &r.Description, &r.Likes, &r.Comments, &r.Image)
		if err != nil {
			http.Error(w, "Failed to scan recipe: "+err.Error(), http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, r)
	}

	cb.Recipes = recipes

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

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		UPDATE cookbooks SET name=?, description=?, likes=?, user=?
		WHERE id=?`,
		cb.Name, cb.Description, cb.Likes, cb.User, cb.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to update cookbook: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM cookbook_recipes WHERE cookbook_id = ?`, cb.ID)
	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to clear old cookbook recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, recipe := range cb.Recipes {
		_, err := tx.Exec(`
			INSERT INTO cookbook_recipes (cookbook_id, recipe_id)
			VALUES (?, ?)`,
			cb.ID, recipe.ID)
		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to insert cookbook recipe link: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
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
