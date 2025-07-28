package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"opskrifter-backend/internal/types"
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
