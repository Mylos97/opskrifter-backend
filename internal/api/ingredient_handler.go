package api

import (
	"encoding/json"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
)

// PUT /ingredient/
func CreatIngredient(w http.ResponseWriter, r *http.Request) {
	var ingredient types.Ingredient
	if err := json.NewDecoder(r.Body).Decode(&ingredient); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO ingredients (name)
		VALUES (?)`,
		ingredient.Name)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert ingredient: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, ingredient)
}

// GET /ingredients/
func GetIngredients(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query(`
		SELECT id, name
		FROM ingredients`)
	if err != nil {
		http.Error(w, "Failed to query ingredients: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var ingredients []types.Ingredient

	for rows.Next() {
		var ingredient types.Ingredient
		if err := rows.Scan(&ingredient.ID, &ingredient.Name); err != nil {
			http.Error(w, "Failed to scan ingredient: "+err.Error(), http.StatusInternalServerError)
			return
		}
		ingredients = append(ingredients, ingredient)
	}

	if err = rows.Err(); err != nil {
		http.Error(w, "Row iteration error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, ingredients)
}
