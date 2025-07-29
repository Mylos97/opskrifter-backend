package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/testutils"
	"opskrifter-backend/internal/types"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRouteCreateRecipe(t *testing.T) {
	body, _ := json.Marshal(testRecipe)
	req := httptest.NewRequest("POST", "/recipes/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	var response Response
	err := json.Unmarshal(resp.Body.Bytes(), &response)
	require.NoError(t, err, "failed to unmarshal response")

	testRecipe.ID = response.ID

	count, err := GetCountByType(testRecipe)
	require.NoError(t, err, "unexpected error from GetCountByType")
	assert.Equal(t, 1, count, "expected count to be 1")
	_, err = DeleteByType[types.Recipe](response.ID)
	require.NoError(t, err, "error deleting recipe")
}

func TestRouteUpdateRecipe(t *testing.T) {
	body, _ := json.Marshal(testRecipe)
	req := httptest.NewRequest("POST", "/recipes/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	require.Equal(t, http.StatusOK, resp.Code)

	var createResp Response
	err := json.Unmarshal(resp.Body.Bytes(), &createResp)
	require.NoError(t, err)
	testRecipe.ID = createResp.ID

	testRecipe = testRecipes[1]
	testRecipe.ID = createResp.ID
	updatedBody, _ := json.Marshal(testRecipe)

	req = httptest.NewRequest("PUT", "/recipes/", bytes.NewBuffer(updatedBody))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
	_, err = DeleteByType[types.Recipe](createResp.ID)
	require.NoError(t, err, "error deleting recipe")
}

func TestRouteDeleteRecipe(t *testing.T) {
	body, _ := json.Marshal(testRecipe)
	req := httptest.NewRequest("POST", "/recipes/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var createResp Response
	err := json.Unmarshal(resp.Body.Bytes(), &createResp)
	require.NoError(t, err)

	testRecipe.ID = createResp.ID
	updatedBody, _ := json.Marshal(testRecipe)

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/recipes/%s", testRecipe.ID), nil)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	req = httptest.NewRequest("GET", fmt.Sprintf("/recipes/%s", testRecipe.ID), bytes.NewBuffer(updatedBody))
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusNotFound, resp.Code)
}

func TestRouteGetRecipe(t *testing.T) {
	body, _ := json.Marshal(testRecipe)
	req := httptest.NewRequest("POST", "/recipes/", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	require.Equal(t, http.StatusOK, resp.Code)

	var createResp Response
	err := json.Unmarshal(resp.Body.Bytes(), &createResp)
	require.NoError(t, err)

	testRecipe.ID = createResp.ID
	updatedBody, _ := json.Marshal(testRecipe)

	req = httptest.NewRequest("GET", fmt.Sprintf("/recipes/%s", testRecipe.ID), bytes.NewBuffer(updatedBody))
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusOK, resp.Code)

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/recipes/%s", testRecipe.ID), nil)
	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestRouteGetManyRecipe(t *testing.T) {
	ids, err := CreateManyByType(testRecipes)
	require.NoError(t, err, "error creating recipes")

	defer func() {
		err := DeleteManyByType[types.Recipe](ids)
		require.NoError(t, err, "error deleting recipes")
	}()

	req := httptest.NewRequest("GET", "/recipes/?page=0&per_page=2&order_by=name", nil)
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "expected status 200 OK")

	var got []types.Recipe
	err = json.NewDecoder(resp.Body).Decode(&got)
	require.NoError(t, err, "error decoding response body")

	require.Len(t, got, 2, "expected 2 recipes in result")
	testutils.AssertSortedBy(t, got, func(a, b types.Recipe) bool {
		return a.Name <= b.Name
	})
}

func TestRouteLikeRecipe(t *testing.T) {
	id, err := CreateByType(testRecipe)
	require.NoError(t, err, "error creating recipes")

	defer func() {
		_, err := DeleteByType[types.Recipe](id)
		require.NoError(t, err, "error deleting recipes")
	}()

	body := map[string]string{"user_id": adminUser.ID}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", fmt.Sprintf("/recipes/%s/like", id), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, "expected status 201")
	updatedRecipe, err := GetByType[types.Recipe](id)
	require.NoError(t, err, "error getting the recipe")
	assert.Equal(t, updatedRecipe.Likes, testRecipe.Likes+1)
	r, err := GetRelationByType[types.UserLikedRecipe](adminUser.ID, id)
	require.NoError(t, err, "error getting the relation")
	assert.Equal(t, r.RecipeID, id)
	assert.Equal(t, r.UserID, adminUser.ID)

	resp = httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)
	assert.Equal(t, http.StatusBadRequest, resp.Code)
}

func TestRouteUnLikeRecipe(t *testing.T) {
	id, err := CreateByType(testRecipe)
	require.NoError(t, err, "error creating recipes")

	defer func() {
		_, err := DeleteByType[types.Recipe](id)
		require.NoError(t, err, "error deleting recipes")
	}()

	body := map[string]string{"user_id": adminUser.ID}
	bodyBytes, _ := json.Marshal(body)

	req := httptest.NewRequest("POST", fmt.Sprintf("/recipes/%s/like", id), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusCreated, resp.Code, "expected status 201")

	req = httptest.NewRequest("DELETE", fmt.Sprintf("/recipes/%s/like", id), bytes.NewReader(bodyBytes))
	req.Header.Set("Content-Type", "application/json")
	resp = httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	updatedRecipe, err := GetByType[types.Recipe](id)
	require.NoError(t, err, "error getting the recipe")
	assert.Equal(t, updatedRecipe.Likes, testRecipe.Likes)

	assert.Equal(t, http.StatusNoContent, resp.Code, "expected status 204")
	_, err = GetRelationByType[types.UserLikedRecipe](adminUser.ID, id)
	require.ErrorIs(t, err, sql.ErrNoRows)
}

func TestRouteUpdateViewsRecipe_Concurrent(t *testing.T) {
	id, err := CreateByType(testRecipe)
	require.NoError(t, err, "error creating recipe")

	defer func() {
		_, err := DeleteByType[types.Recipe](id)
		require.NoError(t, err, "error deleting recipe")
	}()

	const parallelRequests = 1000
	var wg sync.WaitGroup
	wg.Add(parallelRequests)

	for range parallelRequests {
		go func() {
			defer wg.Done()
			req := httptest.NewRequest("POST", fmt.Sprintf("/recipes/%s/views", id), nil)
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()
			testRouter.ServeHTTP(resp, req)

			assert.Equal(t, http.StatusOK, resp.Code)
		}()
	}

	wg.Wait()

	req := httptest.NewRequest("GET", fmt.Sprintf("/recipes/%s", id), nil)
	resp := httptest.NewRecorder()
	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code)

	var updated types.Recipe
	err = json.NewDecoder(resp.Body).Decode(&updated)
	require.NoError(t, err, "error decoding recipe response")

	assert.Equal(t, testRecipe.Views+parallelRequests, updated.Views, "expected view count to match number of requests")
}

func TestGetAllIngredients(t *testing.T) {
	req := httptest.NewRequest("GET", "/ingredients/", nil)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	testRouter.ServeHTTP(resp, req)

	assert.Equal(t, http.StatusOK, resp.Code, "expected status 201")
	var data []types.Ingredient
	err := json.Unmarshal(resp.Body.Bytes(), &data)
	require.NoError(t, err, "failed to decode response JSON")
	count, err := GetCountByTable(data[0].TableName())
	require.NoError(t, err, "failed to get count response JSON")

	assert.Equal(t, len(data), count, "expected ingredient to be the same ")
}
