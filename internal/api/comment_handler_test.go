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

var testUser = types.User{
	ID:   uuid.New().String(),
	Name: "test-user",
}

var testComment = types.Comment{
	Recipe:  testRecipe,
	User:    testUser,
	Comment: "testcomment",
}

func setupTestComment(t *testing.T) {
	_, err := db.DB.Exec("DELETE FROM comments WHERE id = ?", testComment.ID)
	if err != nil {
		t.Fatalf("failed to clean test comment: %v", err)
	}
	insertTestComment(t)
}

func insertTestComment(t *testing.T) {
	_, err := db.DB.Exec(`
        INSERT INTO comments (recipe_id, user_id, comment)
        VALUES (?, ?, ?)`,
		testComment.Recipe.ID, testUser.ID, testComment.Comment)
	if err != nil {
		t.Fatalf("failed to insert test comment: %v", err)
	}
}

func TestCreateComment(t *testing.T) {
	setupTestComment(t)
	payload, err := json.Marshal(testComment)
	if err != nil {
		t.Fatalf("failed to marshal comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/comments", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateComment(w, req)

	res := w.Result()
	defer res.Body.Close()
	LogAndResetBody(t, res)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	var createdComment types.Comment
	if err := json.NewDecoder(res.Body).Decode(&createdComment); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdComment.Comment != testComment.Comment {
		t.Errorf("expected Comment name 'Test Recipe', got '%s'", createdComment.Comment)
	}
}

func TestGetComments(t *testing.T) {
	var recipeID = testRecipe.ID
	setupTestComment(t)
	req := httptest.NewRequest(http.MethodGet, "/comments/"+testRecipe.ID, nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("recipe_id", recipeID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	GetComments(w, req)

	res := w.Result()
	LogAndResetBody(t, res)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec []types.Comment
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(rec) < 1 {
		t.Errorf("Expected 1 comment")
	}
}

func TestUpdateComment(t *testing.T) {
	cookBookID := testCookBook.ID
	var updatedName = "Updated CookBook Name"
	updatedCookBook := testCookBook
	updatedCookBook.Name = updatedName

	payload, err := json.Marshal(updatedCookBook)
	if err != nil {
		t.Fatalf("failed to marshal updated recipe: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/recipes/"+cookBookID, bytes.NewReader(payload))
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

func TestDeleteComment(t *testing.T) {
	cookBookID := testCookBook.ID

	req := httptest.NewRequest(http.MethodDelete, "/recipes/"+cookBookID, nil)
	w := httptest.NewRecorder()

	DeleteCookBook(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", res.StatusCode)
	}

	reqGet := httptest.NewRequest(http.MethodGet, "/recipes/"+cookBookID, nil)
	wGet := httptest.NewRecorder()
	GetRecipe(wGet, reqGet)
	resGet := wGet.Result()
	defer resGet.Body.Close()

	if resGet.StatusCode != http.StatusNotFound {
		t.Errorf("expected status 404 Not Found after delete, got %d", resGet.StatusCode)
	}
}
