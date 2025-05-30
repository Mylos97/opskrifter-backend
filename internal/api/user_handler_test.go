package api

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/types"
	"opskrifter-backend/pkg/db"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
)

var testUser = types.User{
	ID:        "test-id",
	Name:      "test-user",
	CreatedAt: time.Now().String(),
}

func setupTestUser(t *testing.T) {
	_, err := db.DB.Exec("DELETE FROM users WHERE id = ?", testUser.ID)
	if err != nil {
		t.Fatalf("failed to clean test user: %v", err)
	}
	insertTestUser(t)
	checkUserExists(t)
}

func insertTestUser(t *testing.T) {
	_, err := db.DB.Exec(`
        INSERT INTO users (id, name, created_at)
        VALUES (?, ?, ?)`,
		testUser.ID, testUser.Name, testUser.CreatedAt)
	if err != nil {
		t.Fatalf("failed to insert test user: %v", err)
	}
}

func checkUserExists(t *testing.T) {
	var exists bool
	err := db.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM users WHERE id=?)", testUser.ID).Scan(&exists)
	if err != nil || !exists {
		t.Fatalf("User with id %s does not exist in test DB", testUser.ID)
	}
}

func TestCreateUser(t *testing.T) {
	payload, err := json.Marshal(testUser)
	if err != nil {
		t.Fatalf("failed to marshal comment: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/users", bytes.NewReader(payload))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateUser(w, req)

	res := w.Result()
	defer res.Body.Close()
	LogAndResetBody(t, res)
	if res.StatusCode != http.StatusCreated {
		t.Errorf("expected status 201 Created, got %d", res.StatusCode)
	}

	var createdUser types.User
	if err := json.NewDecoder(res.Body).Decode(&createdUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if createdUser.Name != testUser.Name {
		t.Errorf("expected Comment name 'Test Recipe', got '%s'", testUser.Name)
	}
}

func TestGetUser(t *testing.T) {
	var ID = testUser.ID
	setupTestUser(t)
	req := httptest.NewRequest(http.MethodGet, "/users/"+ID, nil)
	w := httptest.NewRecorder()
	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	GetUser(w, req)

	res := w.Result()
	LogAndResetBody(t, res)
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("expected status 200 OK, got %d", res.StatusCode)
	}

	var rec types.User
	if err := json.NewDecoder(res.Body).Decode(&rec); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
}

func TestUpdateUser(t *testing.T) {
	setupTestUser(t)

	updatedUser := types.User{
		ID:   testUser.ID,
		Name: "Updated Name",
	}

	bodyBytes, err := json.Marshal(updatedUser)
	if err != nil {
		t.Fatalf("failed to marshal user: %v", err)
	}

	req := httptest.NewRequest(http.MethodPut, "/users/"+testUser.ID, bytes.NewReader(bodyBytes))
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", testUser.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	UpdateUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200 OK, got %d", res.StatusCode)
	}

	var respUser types.User
	if err := json.NewDecoder(res.Body).Decode(&respUser); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if respUser.Name != updatedUser.Name {
		t.Errorf("expected name %q, got %q", updatedUser.Name, respUser.Name)
	}

	var dbUser types.User
	err = db.DB.QueryRow("SELECT id, name FROM users WHERE id = ?", testUser.ID).Scan(&dbUser.ID, &dbUser.Name)
	if err != nil {
		t.Fatalf("failed to query updated user: %v", err)
	}
	if dbUser.Name != updatedUser.Name {
		t.Errorf("DB update mismatch: expected %q, got %q", updatedUser.Name, dbUser.Name)
	}
}

func TestDeleteUser(t *testing.T) {
	setupTestUser(t)

	req := httptest.NewRequest(http.MethodDelete, "/users/"+testUser.ID, nil)
	w := httptest.NewRecorder()

	ctx := chi.NewRouteContext()
	ctx.URLParams.Add("id", testUser.ID)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, ctx))

	DeleteUser(w, req)

	res := w.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusNoContent {
		t.Fatalf("expected status 204 No Content, got %d", res.StatusCode)
	}

	var id string
	err := db.DB.QueryRow("SELECT id FROM users WHERE id = ?", testUser.ID).Scan(&id)
	if err != sql.ErrNoRows {
		t.Errorf("expected no rows after deletion, but found user with id: %v, err: %v", id, err)
	}
}
