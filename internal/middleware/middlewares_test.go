package middleware

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zeebo/assert"
)

func TestValidateJSONMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		contentType    string
		body           string
		expectedStatus int
	}{
		{
			name:           "Valid JSON",
			contentType:    "application/json",
			body:           `{"name": "example"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Invalid JSON",
			contentType:    "application/json",
			body:           `{"name": "bad json}"`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "Missing Content-Type",
			contentType:    "",
			body:           `{"name": "test"}`,
			expectedStatus: http.StatusUnsupportedMediaType,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ValidateJSONMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("POST", "/", bytes.NewBufferString(tt.body))
			if tt.contentType != "" {
				req.Header.Set("Content-Type", tt.contentType)
			}
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}

func TestRejectSQLInjection(t *testing.T) {
	tests := []struct {
		name           string
		body           string
		expectedStatus int
	}{
		{
			name:           "Clean JSON",
			body:           `{"recipe_id": "1234", "ingredient": "salt"}`,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "SQL Injection Attempt",
			body:           `{"recipe_id": "1; DROP TABLE users; --"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "SQL Injection with OR 1=1",
			body:           `{"ingredient": "' OR '1'='1"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "SQL Injection with comment sequence",
			body:           `{"amount": "0); -- comment"}`,
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "SQL Injection with multiple statements",
			body:           `{"recipe_id": "abc'; DROP TABLE recipes; SELECT * FROM users; --"}`,
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := RejectSQLInjection(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
			}))

			req := httptest.NewRequest("POST", "/", bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			resp := httptest.NewRecorder()

			handler.ServeHTTP(resp, req)

			assert.Equal(t, tt.expectedStatus, resp.Code)
		})
	}
}
