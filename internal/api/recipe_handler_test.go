package api

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
    "opskrifter-backend/internal/types"
)

func TestCreateRecipe(t *testing.T) {
    recipe := types.Recipe{
        ID:          "test-id-123",
        Name:        "Test Recipe",
        Minutes:     20,
        Description: "Test description",
        Likes:       0,
        Comments:    0,
        Image:       "http://example.com/image.jpg",
        Ingredients: []types.IngredientAmount{},
        Categories:  []types.RecipeCategory{},
        RecipeCuisine: types.RecipeCuisine{
            ID:   "it",
            Name: "Italian",
        },
        User: types.User{
            ID: "google-user-123",
        },
    }

    payload, err := json.Marshal(recipe)
    if err != nil {
        t.Fatalf("failed to marshal recipe: %v", err)
    }

    req := httptest.NewRequest(http.MethodPost, "/recipes", bytes.NewReader(payload))
    req.Header.Set("Content-Type", "application/json")
    w := httptest.NewRecorder()

    CreateRecipe(w, req)

    res := w.Result()
    defer res.Body.Close()

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
    recipeID := "test-id-123"

    req := httptest.NewRequest(http.MethodGet, "/recipes/"+recipeID, nil)
    w := httptest.NewRecorder()

    GetRecipe(w, req)

    res := w.Result()
    defer res.Body.Close()

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
