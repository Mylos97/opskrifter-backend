package api

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"
)

func TestCreateHandlerByType(t *testing.T) {
	data, err := json.Marshal(testRecipe)
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

func TestGetManyHandlerByType(t *testing.T) {
	ids, err := CreateManyByType(testRecipes)
	for i := range ids {
		testRecipes[i].ID = ids[i]
	}

	if err != nil {
		t.Fatalf("error creating recipes")
	}

	testCases := []struct {
		name        string
		query       QueryOptions
		expectedLen int
		expectError bool
		validate    func([]types.Recipe) error
	}{
		{
			name: "basic pagination",
			query: QueryOptions{
				Page:    1,
				PerPage: 3,
			},
			expectedLen: 3,
		},
		{
			name: "second page",
			query: QueryOptions{
				Page:    2,
				PerPage: 5,
			},
			expectedLen: 5,
			validate: func(recipes []types.Recipe) error {
				if recipes[0].ID != testRecipes[4].ID {
					return fmt.Errorf("unexpected first item on page 1 got id %s expected %s", recipes[0].ID, testRecipes[4].ID)
				}
				return nil
			},
		},
		{
			name: "invalid per_page",
			query: QueryOptions{
				PerPage: -1,
			},
			expectError: true,
		},
		{
			name: "ordering by name",
			query: QueryOptions{
				PerPage: 5,
				OrderBy: "name",
			},
			expectedLen: 5,
			validate: func(recipes []types.Recipe) error {
				for i := 0; i < len(recipes)-1; i++ {
					if recipes[i].Name > recipes[i+1].Name {
						return fmt.Errorf("items not ordered by name")
					}
				}
				return nil
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			data, err := json.Marshal(tc.query)
			if err != nil {
				t.Fatalf("failed to marshal query: %v", err)
			}

			req, rec := testutils.NewJSONPostRequest(data)
			GetManyRecipe.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if tc.expectError {
				if resp.StatusCode == http.StatusOK {
					t.Error("expected error status code, got 200")
				}
				return
			}

			if resp.StatusCode != http.StatusOK {
				t.Fatalf("expected status 200, got %d", resp.StatusCode)
			}

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Fatalf("failed to read response body: %v", err)
			}

			var objs []types.Recipe
			if err := json.Unmarshal(body, &objs); err != nil {
				t.Fatalf("failed to unmarshal response JSON: %v", err)
			}

			if tc.expectedLen != len(objs) {
				t.Fatalf("expected %d items, got %d", tc.expectedLen, len(objs))
			}

			if tc.validate != nil {
				if err := tc.validate(objs); err != nil {
					t.Errorf("validation failed: %v", err)
				}
			}
		})
	}

	err = DeleteManyByType(testRecipes)

	if err != nil {
		t.Fatalf("error deleting recipes %v", err)
	}
}
