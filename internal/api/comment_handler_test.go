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
	setupTestUser(t)
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
	payload, err := json.Marshal(testComment)
	if err != nil {
		t.Fatalf("failed to marshal comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/comments/", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	setupTestRecipe(t)

	var initialComments int
	err = db.DB.QueryRow(`SELECT comments FROM recipes WHERE id = ?`, testLike.RecipeId).Scan(&initialComments)
	if err != nil {
		t.Fatalf("failed to get initial comments: %v", err)
	}

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

	var updatedComments int
	err = db.DB.QueryRow(`SELECT comments FROM recipes WHERE id = ?`, testLike.RecipeId).Scan(&updatedComments)
	if err != nil {
		t.Fatalf("failed to get updated likes: %v", err)
	}

	if updatedComments != 1 {
		t.Errorf("expected comments to be %d, got %d", 0, updatedComments)
	}
}

func TestGetComments(t *testing.T) {
	var recipeID = testRecipe.ID
	req := httptest.NewRequest(http.MethodGet, "/comments/"+testRecipe.ID, nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("recipe_id", recipeID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setupTestComment(t)
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
	commentID := testComment.ID
	var updated = "Updated Comment"
	updatedComment := testComment
	updatedComment.Comment = updated

	payload, err := json.Marshal(updatedComment)
	if err != nil {
		t.Fatalf("failed to marshal updated comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/comments/"+commentID, bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", commentID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	setupTestComment(t)
	UpdateComment(w, req)

	res := w.Result()
	defer res.Body.Close()

	bodyBytes, _ := io.ReadAll(res.Body)
	t.Logf("Response body: %s", string(bodyBytes))
	res.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec types.Comment
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if rec.Comment != updated {
		t.Errorf("expected updated comment '%s', got '%s'", updated, rec.Comment)
	}
}

func TestDeleteComment(t *testing.T) {
	commentID := testComment.ID

	req := httptest.NewRequest(http.MethodDelete, "/comments/"+commentID, nil)
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", commentID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))
	w := httptest.NewRecorder()

	setupTestCookBook(t)
	insertTestComment(t)
	DeleteComment(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Errorf("expected status 204 No Content, got %d", res.StatusCode)
	}

	var updatedComments int
	err := db.DB.QueryRow(`SELECT comments FROM recipes WHERE id = ?`, testLike.RecipeId).Scan(&updatedComments)
	if err != nil {
		t.Fatalf("failed to get updated likes: %v", err)
	}

	if updatedComments != 0 {
		t.Errorf("expected likes to be %d, got %d", 0, updatedComments)
	}
}
