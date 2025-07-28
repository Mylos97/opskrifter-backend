package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateHandlerByType(t *testing.T) {
	data, err := json.Marshal(testRecipe)
	require.NoError(t, err, "failed to marshal testRecipe")

	req, rec := testutils.NewJSONPostRequest(data)
	CreateRecipe.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 OK")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	var response Response
	require.NoError(t, json.Unmarshal(body, &response), "failed to unmarshal response JSON")
	require.NotEmpty(t, response.ID, "expected non-empty id field in response")
	require.NoError(t, testutils.AssertCountByType[types.Recipe](1, GetCountByType))

	id := response.ID
	_, err = DeleteByType[types.Recipe](id)
	require.NoError(t, err, "failed to delete recipe")
	require.NoError(t, testutils.AssertCountByType[types.Recipe](0, GetCountByType))
}

func TestDeleteHandlerByType(t *testing.T) {
	id, err := CreateByType(handlerRecipe)
	require.NoError(t, err, "failed to create recipe")
	require.NotEmpty(t, id, "failed to generate id")
	req := httptest.NewRequest(http.MethodDelete, "/recipes/"+id, nil)
	rec := httptest.NewRecorder()
	chiCtx := chi.NewRouteContext()
	chiCtx.URLParams.Add("id", id)
	req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, chiCtx))

	DeleteRecipe.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 OK")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	var response Response
	require.NoError(t, json.Unmarshal(body, &response), "failed to unmarshal response JSON")
	require.NotEmpty(t, response.ID, "expected non-empty id field in response")
	require.NoError(t, testutils.AssertCountByType[types.Recipe](0, GetCountByType))
}

func TestUpdateHandlerByType(t *testing.T) {
	id, err := CreateByType(handlerRecipe)
	require.NoError(t, err, "failed to create recipe")
	require.NotEmpty(t, id, "failed to generate id")

	updatedRecipe := recipeGenerator.Generate()
	updatedRecipe.ID = id

	data, err := json.Marshal(updatedRecipe)
	require.NoError(t, err, "failed to marshal updatedRecipe")

	req, rec := testutils.NewJSONPostRequest(data)
	UpdateRecipe.ServeHTTP(rec, req)

	resp := rec.Result()
	defer resp.Body.Close()

	require.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200 OK")

	body, err := io.ReadAll(resp.Body)
	require.NoError(t, err, "failed to read response body")

	var response Response
	require.NoError(t, json.Unmarshal(body, &response), "failed to unmarshal response JSON")

	require.NotEmpty(t, response.ID, "expected non-empty id field in response")

	require.NoError(t, testutils.AssertCountByType[types.Recipe](1, GetCountByType), "failed to get the count")

	updated, err := GetByType(updatedRecipe)
	require.NoError(t, err, "error fetching updated recipe")

	testutils.EqualByValue(updatedRecipe, updated)

	_, err = DeleteByType[types.Recipe](id)
	require.NoError(t, err, "error deleting recipe")

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
				if recipes[0].ID != testRecipes[5].ID {
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
			require.NoError(t, err, "failed to marshal query")

			req, rec := testutils.NewJSONPostRequest(data)
			GetManyRecipe.ServeHTTP(rec, req)

			resp := rec.Result()
			defer resp.Body.Close()

			if tc.expectError {
				assert.NotEqual(t, http.StatusOK, resp.StatusCode, "expected error status code, got 200")
				return
			}

			require.Equal(t, http.StatusOK, resp.StatusCode, "expected status 200")

			body, err := io.ReadAll(resp.Body)
			require.NoError(t, err, "failed to read response body")

			var objs []types.Recipe
			require.NoError(t, json.Unmarshal(body, &objs), "failed to unmarshal response JSON")

			require.Equal(t, tc.expectedLen, len(objs), "unexpected number of items")

			if tc.validate != nil {
				err = tc.validate(objs)
				assert.NoError(t, err, "validation failed")
			}
		})
	}

	err = DeleteManyByType[types.Recipe](ids)
	require.NoError(t, err, "error deleting recipes")
}
