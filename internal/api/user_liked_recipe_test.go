package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
	"testing"

	"github.com/go-chi/chi/v5"
)

var testLike = types.UserLikedRecipe{
	RecipeId: testRecipe.ID,
	UserId:   testUser.ID,
}

func insertTestLike(t *testing.T) {
	_, err := db.DB.Exec(`DELETE FROM user_liked_recipe WHERE user_id = ? AND recipe_id = ?`, testUser.ID, testRecipe.ID)
	if err != nil {
		t.Fatalf("failed to delete test like: %v", err)
	}
	_, err = db.DB.Exec(`INSERT INTO user_liked_recipe (user_id, recipe_id) VALUES (?, ?)`, testUser.ID, testRecipe.ID)
	if err != nil {
		t.Fatalf("failed to insert test like: %v", err)
	}
}

func TestLikeRecipe(t *testing.T) {
	payload, err := json.Marshal(testLike)
	if err != nil {
		t.Fatalf("failed to marshal comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/likerecipe", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	setupTestUser(t)
	setupTestRecipe(t)
	LikeRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	var newEntity types.UserLikedRecipe
	if err := json.NewDecoder(res.Body).Decode(&newEntity); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if newEntity.RecipeId != testLike.RecipeId && newEntity.UserId != testLike.UserId {
		t.Errorf("expected Comment name 'Test Recipe', got '%s'", testUser.Name)
	}
}

func TestUnLikeRecipe(t *testing.T) {
	payload, err := json.Marshal(testLike)
	if err != nil {
		t.Fatalf("failed to marshal comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/likerecipe", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupTestUser(t)
	setupTestRecipe(t)
	insertTestLike(t)

	UnLikeRecipe(w, req)

	res := w.Result()
	defer res.Body.Close()
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}
}

func TestLikedRecipes(t *testing.T) {
	var userId = testUser.ID

	req := httptest.NewRequest(http.MethodGet, "/likedRecipes/"+userId, nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("user_id", userId)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	insertTestLike(t)

	GetLikeDRecipes(w, req)
	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec []types.UserLikedRecipe
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(rec) < 1 {
		t.Errorf("Expected 1 like")
	}
}
