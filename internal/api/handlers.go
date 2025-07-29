package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"opskrifter-backend/internal/types"

	"github.com/go-chi/chi/v5"
)

var body struct {
	UserID string `json:"user_id"`
}

var CreateRecipe = HandlerByType(CreateByTypeWithRelations[types.Recipe])
var UpdateRecipe = HandlerByType(UpdateByType[types.Recipe])
var DeleteRecipe = DeleteHandlerByType[types.Recipe](DeleteByType[types.Recipe])
var GetRecipe = GetHandlerByType(GetByType[types.Recipe])
var GetManyRecipe = GetHandlerManyByType(GetManyByType[types.Recipe])
var GetManyIngredients = GetAllHandlerManyByType(GetAllByType[types.Ingredient])

func UnlikeRecipe(w http.ResponseWriter, r *http.Request) {
	recipeID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if recipeID == "" || body.UserID == "" {
		http.Error(w, "missing user_id or recipe_id", http.StatusBadRequest)
		return
	}
	err := DeleteRelationByType[types.UserLikedRecipe](body.UserID, recipeID)
	if err != nil {
		http.Error(w, "could not unlike recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}
	var recipe types.Recipe
	recipe.ID = recipeID
	err = UpdateCountByType(recipe, "likes", "-1")

	if err != nil {
		http.Error(w, "could not decrement likes recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func LikeRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe types.Recipe
	recipeID := chi.URLParam(r, "id")

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		http.Error(w, "invalid body", http.StatusBadRequest)
		return
	}

	if body.UserID == "" || recipeID == "" {
		http.Error(w, "missing user_id or recipe_id", http.StatusBadRequest)
		return
	}
	relations := []types.UserLikedRecipe{
		{
			UserID:   body.UserID,
			RecipeID: recipeID,
		},
	}

	err := CreateOneToManyByType(recipe, body.UserID, relations)

	if err != nil {
		http.Error(w, "could not like recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	recipe.ID = recipeID
	err = UpdateCountByType(recipe, "likes", "+1")

	if err != nil {
		http.Error(w, "could not increment likes recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func UpdateViewRecipe(w http.ResponseWriter, r *http.Request) {
	var recipe types.Recipe
	recipeID := chi.URLParam(r, "id")

	if recipeID == "" {
		http.Error(w, "missing recipe_id", http.StatusBadRequest)
		return
	}
	recipe.ID = recipeID
	err := UpdateCountByType(recipe, "views", "+1")

	if errors.Is(err, ErrNoIdForType) {
		http.Error(w, "no id for the type: "+err.Error(), http.StatusBadRequest)
		return
	}

	if err != nil {
		fmt.Printf("Eec error: %v\n", err)
		http.Error(w, "could not update views recipe: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
