package api

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
	"testing"

	"github.com/go-chi/chi/v5"
)

var testIngredient = types.Ingredient{
	Name: "Tomat",
}

var testRecipe = types.Recipe{
	ID:          "test-id-123",
	Name:        "Test Recipe",
	Minutes:     20,
	Description: "Test description",
	Likes:       0,
	Comments:    0,
	Image:       "http://example.com/image.jpg",
	Ingredients: []types.IngredientAmount{
		{
			Ingredient: testIngredient,
			Amount:     "2 stk",
		},
	},
	Categories: []types.RecipeCategory{},
	RecipeCuisine: types.RecipeCuisine{
		ID:   "it",
		Name: "Italian",
	},
	User: types.User{
		ID: "google-user-123",
	},
}

func setupTestRecipe(t *testing.T) {
	_, err := db.DB.Exec("DELETE FROM recipes WHERE id = ?", testRecipe.ID)
	if err != nil {
		t.Fatalf("failed to clean test recipe: %v", err)
	}
	insertTestRecipe(t, testRecipe)
}

func insertTestRecipe(t *testing.T, recipe types.Recipe) {
	dbRec := types.RecipeDB{
		ID:            recipe.ID,
		Name:          recipe.Name,
		Minutes:       recipe.Minutes,
		Description:   recipe.Description,
		Likes:         recipe.Likes,
		Comments:      recipe.Comments,
		Image:         recipe.Image,
		RecipeCuisine: recipe.RecipeCuisine.ID,
		User:          recipe.User.ID,
	}

	_, err := db.DB.Exec(`
        INSERT INTO recipes (id, name, minutes, description, likes, comments, image, recipe_cuisine, user)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		dbRec.ID, dbRec.Name, dbRec.Minutes, dbRec.Description,
		dbRec.Likes, dbRec.Comments, dbRec.Image, dbRec.RecipeCuisine, dbRec.User)
	if err != nil {
		t.Fatalf("failed to insert test recipe: %v", err)
	}
}

func TestCreateRecipe(t *testing.T) {
	setupTestRecipe(t)
	payload, err := json.Marshal(testRecipe)
	if err != nil {
		t.Fatalf("failed to marshal recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)
	t.Logf("Response body: %s", string(bodyBytes))
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	var createdRecipe types.Recipe
	if err := json.NewDecoder(res.Body).Decode(&createdRecipe); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdRecipe.Name != "Test Recipe" {
		t.Errorf("expected recipe name 'Test Recipe', got '%s'", createdRecipe.Name)
	}
}

func TestGetRecipe(t *testing.T) {
	var recipeID = testRecipe.ID
	setupTestRecipe(t)
	req := httptest.NewRequest(http.MethodGet, "/recipes/"+recipeID, nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", recipeID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	GetRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()
	bodyBytes, _ := io.ReadAll(res.Body)
	t.Logf("Response body: %s", string(bodyBytes))
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec types.Recipe
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if rec.ID != recipeID {
		t.Errorf("expected recipe id '%s', got '%s'", recipeID, rec.ID)
	}
}

func TestUpdateRecipe(t *testing.T) {
	recipeID := testRecipe.ID
	var updatedName = "Updated Recipe Name"
	updatedRecipe := testRecipe
	updatedRecipe.Name = updatedName

	payload, err := json.Marshal(updatedRecipe)
	if err != nil {
		t.Fatalf("failed to marshal updated recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/recipes/"+recipeID, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", recipeID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	UpdateRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	t.Logf("Response body: %s", string(bodyBytes))
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec types.Recipe
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if rec.Name != updatedName {
		t.Errorf("expected updated recipe name '%s', got '%s'", updatedName, rec.Name)
	}
}

func TestDeleteRecipe(t *testing.T) {
	recipeID := testRecipe.ID

	req := httptest.NewRequest(http.MethodDelete, "/recipes/"+recipeID, nil)
	w := httptest.NewRecorder()

	DeleteRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", res.StatusCode)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/recipes/"+recipeID, nil)
	wGet := httptest.NewRecorder()
	GetRecipe(wGet, reqGet)
	resGet := wGet.Result()
	defer resGet.Body.Close()

	if resGet.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found after delete, got %d", resGet.StatusCode)
	}
}
