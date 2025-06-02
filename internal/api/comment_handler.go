package api

import (
	"encoding/json"
	"net/http"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

// POST /cookbooks/
func CreateComment(w http.ResponseWriter, r *http.Request) {
	var comment types.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	comment.ID = uuid.New().String()

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(`
		INSERT INTO comments (recipe_id, user_id, comment)
		VALUES (?, ?, ?)`,
		comment.Recipe.ID, comment.User.ID, comment.Comment)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, comment)
}

// GET /comments/{recipe_id}?index=0&max=10
func GetComments(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "recipe_id")
	indexStr := r.URL.Query().Get("index")
	maxStr := r.URL.Query().Get("max")

	index := 0
	max := 10
	var err error

	if indexStr != "" {
		index, err = strconv.Atoi(indexStr)
		if err != nil || index < 0 {
			http.Error(w, "Invalid index parameter", http.StatusBadRequest)
			return
		}
	}
	if maxStr != "" {
		max, err = strconv.Atoi(maxStr)
		if err != nil || max <= 0 {
			http.Error(w, "Invalid max parameter", http.StatusBadRequest)
			return
		}
	}

	// Query multiple comments with LIMIT and OFFSET
	rows, err := db.DB.Query(`SELECT
				comments.id,
				comments.comment,
				users.id,
				users.name
			FROM comments
			JOIN users ON comments.user_id = users.id
			WHERE comments.recipe_id = ?
			LIMIT ? OFFSET ?`, recipeID, max, index)

	if err != nil {
		http.Error(w, "DB error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	comments := []types.Comment{}

	for rows.Next() {
		var c types.Comment
		var u types.User

		if err := rows.Scan(&c.ID, &c.Comment, &u.ID, &u.Name); err != nil {
			http.Error(w, "DB error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		c.User = u
		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "DB error after rows: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, comments)
}

// PUT /cookbooks/{id}
func UpdateComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var comment types.Comment
	if err := json.NewDecoder(r.Body).Decode(&comment); err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if comment.ID != "" && comment.ID != id {
		http.Error(w, "ID in URL and body do not match", http.StatusBadRequest)
		return
	}
	comment.ID = id

	_, err = db.DB.Exec(`
		UPDATE comments SET comment = ? WHERE id = ?`,
		comment.Comment, comment.ID)

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to insert comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, comment)
}

// DELETE /cookbooks/{id}
func DeleteComment(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	tx, err := db.DB.Begin()
	if err != nil {
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec(`DELETE FROM comments WHERE id = ?`, id)
	if err != nil {
		http.Error(w, "Failed to delete comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err != nil {
		tx.Rollback()
		http.Error(w, "Failed to delete comment: "+err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(); err != nil {
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
