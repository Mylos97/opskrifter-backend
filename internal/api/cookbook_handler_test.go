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
	"github.com/google/uuid"
)

var testCookBook = types.CookBook{
	ID:          uuid.New().String(),
	Name:        "Test CookBook",
	Description: "Test cookbook description",
	Likes:       5,
	User:        "test-user-123",
	Recipes:     []types.Recipe{testRecipe},
}

func setupTestCookBook(t *testing.T) {
	_, err := db.DB.Exec("DELETE FROM cookbooks WHERE id = ?", testCookBook.ID)
	if err != nil {
		t.Fatalf("failed to clean test recipe: %v", err)
	}
	insertTestCookBook(t, testCookBook)
}

func insertTestCookBook(t *testing.T, cookbook types.CookBook) {

	_, err := db.DB.Exec(`
        INSERT INTO cookbooks (id, name, description, likes, user_id)
        VALUES (?, ?, ?, ?, ?)`,
		cookbook.ID, cookbook.Name, cookbook.Description, cookbook.Likes, cookbook.User,
	)
	if err != nil {
		t.Fatalf("failed to insert test recipe: %v", err)
	}
}

func TestCreateCookBook(t *testing.T) {
	payload, err := json.Marshal(testCookBook)
	if err != nil {
		t.Fatalf("failed to marshal recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/cookbooks", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateCookBook(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	var createdCookBook types.CookBook
	if err := json.NewDecoder(res.Body).Decode(&createdCookBook); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdCookBook.Name != testCookBook.Name {
		t.Errorf("expected recipe name 'Test Recipe', got '%s'", createdCookBook.Name)
	}
}

func TestGetCookbook(t *testing.T) {
	var cookBookID = testCookBook.ID
	setupTestCookBook(t)
	req := httptest.NewRequest(http.MethodGet, "/cookbooks/"+cookBookID, nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", cookBookID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	GetCookBook(w, req)

	res := w.Result()
	LogAndResetBody(t, res)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec types.Recipe
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if rec.ID != cookBookID {
		t.Errorf("expected recipe id '%s', got '%s'", cookBookID, rec.ID)
	}
}

func TestUpdateCookBook(t *testing.T) {
	cookBookID := testCookBook.ID
	var updatedName = "Updated CookBook Name"
	updatedCookBook := testCookBook
	updatedCookBook.Name = updatedName

	payload, err := json.Marshal(updatedCookBook)
	if err != nil {
		t.Fatalf("failed to marshal updated recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/cookbooks/"+cookBookID, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", cookBookID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	UpdateCookBook(w, req)

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

func TestDeleteCookBook(t *testing.T) {
	cookBookID := testCookBook.ID

	req := httptest.NewRequest(http.MethodDelete, "/cookbooks/"+cookBookID, nil)
	w := httptest.NewRecorder()

	DeleteCookBook(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", res.StatusCode)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/cookbooks/"+cookBookID, nil)
	wGet := httptest.NewRecorder()
	GetRecipe(wGet, reqGet)
	resGet := wGet.Result()
	defer resGet.Body.Close()

	if resGet.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found after delete, got %d", resGet.StatusCode)
	}
}
