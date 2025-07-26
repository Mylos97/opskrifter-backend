package api

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"os"
	"testing"
)

var handlerRecipe = recipeGenerator.Generate()

func TestCreateHandlerByType(t *testing.T) {
	data, err := os.ReadFile("../testdata/recipe.json")
	if err != nil {
		t.Fatalf("failed to read input file: %v", err)
	}

	req, rec := testutils.NewJSONPostRequest(data)
	CreateRecipe.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	if response.ID == "" {
		t.Errorf("expected non-empty id field in response")
	}

	err = testutils.AssertCountByType[types.Recipe](1, GetCountByType)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}

	var recipe types.Recipe
	recipe.ID = response.ID
	_, err = DeleteByType(recipe)

	if err != nil {
		t.Errorf("failed to delete recipe")
	}

	err = testutils.AssertCountByType[types.Recipe](0, GetCountByType)

	if err != nil {
		t.Fatalf("failed to get the count %v", err)
	}
}

func TestDeleteHandlerByType(t *testing.T) {
	id, err := CreateByType(handlerRecipe)

	if err != nil {
		t.Fatalf("failed to create recipe")
	}

	if id == "" {
		t.Fatalf("failed to generate id")
	}

	handlerRecipe.ID = id

	data, err := json.Marshal(handlerRecipe)
	if err != nil {
		t.Fatalf("failed to marshal handlerRecipe: %v", err)
	}

	req, rec := testutils.NewJSONPostRequest(data)
	DeleteRecipe.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	if response.ID == "" {
		t.Errorf("expected non-empty id field in response")
	}

	err = testutils.AssertCountByType[types.Recipe](0, GetCountByType)

	if err != nil {
		t.Fatalf("failed get the count %v", err)
	}
}

func TestUpdateHandlerByType(t *testing.T) {
	id, err := CreateByType(handlerRecipe)

	if err != nil {
		t.Fatalf("failed to create recipe")
	}

	if id == "" {
		t.Fatalf("failed to generate id")
	}

	updatedRecipe := recipeGenerator.Generate()
	updatedRecipe.ID = id

	data, err := json.Marshal(updatedRecipe)
	if err != nil {
		t.Fatalf("failed to marshal handlerRecipe: %v", err)
	}

	req, rec := testutils.NewJSONPostRequest(data)
	UpdateRecipe.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatalf("failed to unmarshal response JSON: %v", err)
	}

	if response.ID == "" {
		t.Errorf("expected non-empty id field in response")
	}

	err = testutils.AssertCountByType[types.Recipe](1, GetCountByType)

	if err != nil {
		t.Fatalf("failed get the count %v", err)
	}

	updated, err := GetByType(updatedRecipe)

	if err != nil {
		log.Fatalf("error updating recipe")
	}

	testutils.EqualByValue(updatedRecipe, updated)
}
