package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"testing"

	"github.com/go-chi/chi/v5"
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

	updated, err := GetByType[types.Recipe](updatedRecipe.ID)
	require.NoError(t, err, "error fetching updated recipe")

	testutils.EqualByValue(updatedRecipe, updated)

	_, err = DeleteByType[types.Recipe](id)
	require.NoError(t, err, "error deleting recipe")

}

func TestGetManyHandlerByType(t *testing.T) {
	ids, err := CreateManyByType(testRecipes)
	require.NoError(t, err, "error creating recipes")

	defer func() {
		err = DeleteManyByType[types.Recipe](ids)
		require.NoError(t, err, "error deleting recipes")
	}()

	req := httptest.NewRequest("GET", "/?page=0&per_page=5&order_by=name", nil)
	rec := httptest.NewRecorder()

	GetManyRecipe.ServeHTTP(rec, req)

	res := rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode, "expected status OK")

	var got []types.Recipe
	err = json.NewDecoder(res.Body).Decode(&got)
	require.NoError(t, err, "error decoding response")

	require.Len(t, got, 5, "expected 5 recipes in result")
	testutils.AssertSortedBy(t, got, func(a, b types.Recipe) bool {
		return a.Name <= b.Name
	})

	req = httptest.NewRequest("GET", "/?page=0&per_page=5&order_by=loool", nil)
	rec = httptest.NewRecorder()

	GetManyRecipe.ServeHTTP(rec, req)

	res = rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusBadRequest, res.StatusCode, "expected status OK")

	req = httptest.NewRequest("GET", "/?page=2&per_page=5", nil)
	rec = httptest.NewRecorder()

	GetManyRecipe.ServeHTTP(rec, req)

	res = rec.Result()
	defer res.Body.Close()

	require.Equal(t, http.StatusOK, res.StatusCode, "expected status OK")
	err = json.NewDecoder(res.Body).Decode(&got)
	require.NoError(t, err, "error decoding response")

	require.Len(t, got, 5, "expected 5 recipes in result")
	require.Equal(t, ids[5], got[0].ID)
}
