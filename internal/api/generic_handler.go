package api

import (
	"encoding/json"
	"log"
	"net/http"
	"opskrifter-backend/internal/types"
)

type Response struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func HandlerByType[T types.Identifiable](crudFunc CrudFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj T

		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		id, err := crudFunc(obj)
		if err != nil {
			http.Error(w, "operation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		resp := Response{
			ID:      id,
			Message: "operation succeeded",
		}
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			log.Printf("failed to encode response: %v", err)
		}
	}
}

func GetHandlerByType[T types.Identifiable](getFunc GetFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var obj T

		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
			return
		}

		if obj.GetID() == "" {
			http.Error(w, "missing ID field", http.StatusBadRequest)
			return
		}

		result, err := getFunc(obj)
		if err != nil {
			http.Error(w, "operation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
