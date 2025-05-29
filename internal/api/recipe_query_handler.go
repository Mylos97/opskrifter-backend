package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"opskrifter-backend/pkg/db"
	"opskrifter-backend/internal/types"
)

func GetRecipesHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	index := 0
	max := 10
	sort := "name"

	if val := query.Get("index"); val != "" {
		if i, err := strconv.Atoi(val); err == nil && i >= 0 {
			index = i
		}
	}

	if val := query.Get("max"); val != "" {
		if m, err := strconv.Atoi(val); err == nil && m > 0 && m <= 100 {
			max = m
		}
	}

	validSorts := map[string]bool{"name": true, "rating": true, "minutes": true}
	if val := query.Get("sort"); val != "" {
		if validSorts[val] {
			sort = val
		}
	}

	queryStr := fmt.Sprintf(`SELECT id, name, minutes, rating, description, likes, comments, image 
							 FROM recipes 
							 ORDER BY %s 
							 LIMIT ? OFFSET ?`, sort)

	rows, err := db.DB.Query(queryStr, max, index)
	if err != nil {
		http.Error(w, "Failed to fetch recipes: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	recipes := []types.Recipe{}

	for rows.Next() {
		var rec types.Recipe
		err := rows.Scan(&rec.ID, &rec.Name, &rec.Minutes, &rec.Rating, &rec.Description, &rec.Likes, &rec.Comments, &rec.Image)
		if err != nil {
			http.Error(w, "Failed to scan recipe: "+err.Error(), http.StatusInternalServerError)
			return
		}
		recipes = append(recipes, rec)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(recipes)
}
