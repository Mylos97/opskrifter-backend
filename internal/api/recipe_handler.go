package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// PUT /recipes/
func CreateRecipe(w http.ResponseWriter, r *http.Request) {
	var rec types.Recipe
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	rec.ID = uuid.New().String()

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO recipes (id, name, minutes, description, likes, comments, image, user_id, recipe_cuisine)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rec.ID, rec.Name, rec.Minutes, rec.Description, rec.Likes, rec.Comments, rec.Image, rec.User.ID, rec.RecipeCuisine)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	for _, ing := range rec.Ingredients {
		_, err := tx.Exec(`
			INSERT INTO ingredients_for_recipe (recipe_id, name)
			VALUES (?, ?)`,
			rec.ID, ing.Name)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to insert ingredients: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	for _, cat := range rec.Categories {
		_, err := tx.Exec(`
			INSERT INTO recipe_categories (recipe_id, category)
			VALUES (?, ?)`,
			rec.ID, cat.Category)

		if err != nil {
			tx.Rollback()
			http.Error(w, "Failed to insert categories: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, rec)
}

// GET /recipes/{id}
func GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	rowRecipe := db.DB.QueryRow(`
		SELECT id, name, minutes, description, likes, comments, image
		FROM recipes WHERE id = ?`, id)

	var rec types.Recipe
	err := rowRecipe.Scan(&rec.ID, &rec.Name, &rec.Minutes, &rec.Description, &rec.Likes, &rec.Comments, &rec.Image)

	if err == sql.ErrNoRows {
		http.Error(w, "Recipe not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.DB.Query(`SELECT name FROM ingredients_for_recipe WHERE recipe_id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to fetch ingredients: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	for rows.Next() {
		var ing types.IngredientsForRecipe
		if err := rows.Scan(&ing.Name); err != nil {
			log.Printf("Failed to scan ingredient: %v", err)
			continue
		}
		rec.Ingredients = append(rec.Ingredients, ing)
	}

	writeJSON(w, http.StatusOK, rec)
}

// PUT /recipes/{id}
func UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var rec types.Recipe
	if err := json.NewDecoder(r.Body).Decode(&rec); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if rec.ID != "" && rec.ID != id {
		http.Error(w, "ID in URL and body do not match", http.StatusBadRequest)
		return
	}
	rec.ID = id

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to start transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Update the recipe itself
	_, err = tx.Exec(`
		UPDATE recipes SET name=?, minutes=?, description=?, likes=?, comments=?, image=?
		WHERE id=?`,
		rec.Name, rec.Minutes, rec.Description, rec.Likes, rec.Comments, rec.Image, rec.ID)
	if err != nil {
		http.Error(w, "Failed to update recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Delete existing ingredients for the recipe
	_, err = tx.Exec(`DELETE FROM ingredients_for_recipe WHERE recipe_id = ?`, rec.ID)
	if err != nil {
		http.Error(w, "Failed to delete existing ingredients: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Insert new ingredients
	stmt, err := tx.Prepare(`INSERT INTO ingredients_for_recipe (recipe_id, name) VALUES (?, ?)`)
	if err != nil {
		http.Error(w, "Failed to prepare ingredient insert: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer stmt.Close()

	for _, ing := range rec.Ingredients {
		if _, err := stmt.Exec(rec.ID, ing.Name); err != nil {
			http.Error(w, "Failed to insert ingredient: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, rec)
}

// DELETE /recipes/{id}
func DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	_, err := db.DB.Exec(`DELETE FROM recipes WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
