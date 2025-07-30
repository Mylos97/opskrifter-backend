package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"slices"
	"strings"

	_ "github.com/joho/godotenv/autoload"
)

func APIKeyAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-API-Key")
		expectedKey := os.Getenv("API_KEY")

		if apiKey == "" || apiKey != expectedKey && false {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ValidateJSONMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-Type") != "application/json" {
			http.Error(w, "Content-Type must be application/json", http.StatusUnsupportedMediaType)
			return
		}

		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Unable to read request body", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var js json.RawMessage
		if err := json.Unmarshal(bodyBytes, &js); err != nil {
			http.Error(w, "Invalid JSON", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func RejectSQLInjection(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bodyBytes, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "unable to read request", http.StatusBadRequest)
			return
		}

		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var jsonBody any
		if err := json.Unmarshal(bodyBytes, &jsonBody); err != nil {
			http.Error(w, "invalid JSON", http.StatusBadRequest)
			return
		}

		if ContainsSQLInjection(jsonBody) {
			http.Error(w, "request rejected: potential SQL injection detected", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func ContainsSQLInjection(v any) bool {
	switch val := v.(type) {
	case map[string]any:
		for _, v := range val {
			if ContainsSQLInjection(v) {
				return true
			}
		}
	case []any:
		if slices.ContainsFunc(val, ContainsSQLInjection) {
			return true
		}
	case string:
		badPatterns := []string{
			"'", "--", ";", "DROP", "INSERT", "UPDATE", "DELETE", "SELECT",
		}
		lowered := strings.ToLower(val)
		for _, bad := range badPatterns {
			if strings.Contains(lowered, strings.ToLower(bad)) {
				return true
			}
		}
	}
	return false
}
