package api

import (
	"encoding/json"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"

	"github.com/go-chi/chi/v5"
)

// PUT /like_recipe/like
func LikeRecipe(w http.ResponseWriter, r *http.Request) {
	var o types.UserLikedRecipe

	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO user_liked_recipe (user_id, recipe_id)
		VALUES (?, ?)`,
		o.RecipeId, o.UserId)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to like recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, o)
}

// PUT /like_recipe/unlike
func UnLikeRecipe(w http.ResponseWriter, r *http.Request) {
	var o types.UserLikedRecipe

	if err := json.NewDecoder(r.Body).Decode(&o); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`DELETE FROM user_liked_recipe WHERE user_id = ? AND recipe_id = ?`, o.RecipeId, o.UserId)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete liked recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, o)
}

// GET like_recipe/{user_id}
func GetLikeDRecipes(w http.ResponseWriter, r *http.Request) {
	user_id := chi.URLParam(r, "user_id")

	rows, err := db.DB.Query(`SELECT user_id, recipe_id FROM user_liked_recipe WHERE user_id = ?`, user_id)
	if err != nil {
		http.Error(w, "Failed to get liked recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()

	likes := []types.UserLikedRecipe{}
	for rows.Next() {
		var ulr types.UserLikedRecipe
		if err := rows.Scan(&ulr.UserId, &ulr.RecipeId); err != nil {
			http.Error(w, "DB error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		likes = append(likes, ulr)
	}

	writeJSON(w, http.StatusOK, likes)
}
