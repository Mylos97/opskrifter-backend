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

var testIngredient2 = types.Ingredient{
	ID:   "test-id2",
	Name: "test-ingredient2",
}

func setupTestIngredients(t *testing.T, i types.Ingredient) {
	deleteTestIngredient(t, i)
	insertTestIngredient(t, i)
}

func insertTestIngredient(t *testing.T, i types.Ingredient) {
	_, err := db.DB.Exec(`
        INSERT INTO ingredients (name)
        VALUES (?)`, i.Name)
	if err != nil {
		t.Fatalf("failed to insert ingredient: %v", err)
	}
}

func deleteTestIngredient(t *testing.T, i types.Ingredient) {
	_, err := db.DB.Exec("DELETE FROM ingredients WHERE id = ?", i.ID)
	if err != nil {
		t.Fatalf("failed to clean test user: %v", err)
	}
}
func TestCreateIngredient(t *testing.T) {
	payload, err := json.Marshal(testIngredient)
	if err != nil {
		t.Fatalf("failed to marshal comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	CreatIngredient(w, req)

	res := w.Result()
	defer res.Body.Close()
	LogAndResetBody(t, res)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	var createdIngredient types.Ingredient
	if err := json.NewDecoder(res.Body).Decode(&createdIngredient); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdIngredient.Name != testIngredient.Name {
		t.Errorf("expected Comment name 'Test Recipe', got '%s'", testUser.Name)
	}
}

func TestGetIngredients(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/ingredients/", nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	insertTestIngredient(t, testIngredient)
	insertTestIngredient(t, testIngredient2)

	GetIngredients(w, req)

	res := w.Result()
	LogAndResetBody(t, res)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec []types.Ingredient
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(rec) < 1 {
		t.Errorf("Expected more than 1 ingredient")
	}

	deleteTestIngredient(t, testIngredient)
	deleteTestIngredient(t, testIngredient2)
}
