package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"opskrifter-backend/internal/types"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type Response struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type (
	DeleteFunc[T types.Identifiable]  func(id string) (string, error)
	CrudFunc[T types.Identifiable]    func(T) (string, error)
	GetFunc[T types.Identifiable]     func(id string) (T, error)
	GetManyFunc[T types.Identifiable] func(q QueryOptions) ([]T, error)
)

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

func DeleteHandlerByType[T types.Identifiable](deleteFunc DeleteFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		_, err := deleteFunc(id)

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

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
		id := chi.URLParam(r, "id")
		if id == "" {
			http.Error(w, "missing id", http.StatusBadRequest)
			return
		}

		result, err := getFunc(id)

		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}

		if err != nil {
			http.Error(w, "operation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}

func GetHandlerManyByType[T types.Identifiable](getManyFunc GetManyFunc[T]) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		page, err := strconv.ParseInt(query.Get("page"), 32, 32)
		if err != nil {
			http.Error(w, "err parsing page: "+err.Error(), http.StatusBadRequest)
			return
		}

		perPage, err := strconv.ParseInt(query.Get("per_page"), 32, 32)
		if err != nil {
			http.Error(w, "err parsing page: "+err.Error(), http.StatusBadRequest)
			return
		}

		ops := QueryOptions{
			Page:    int(page),
			PerPage: int(perPage),
			OrderBy: query.Get("order_by"),
		}

		result, err := getManyFunc(ops)

		if err == ErrNotValidOrderBy {
			http.Error(w, "operation failed "+err.Error(), http.StatusBadRequest)
		}
		if err != nil {
			http.Error(w, "operation failed: "+err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(result)
	}
}
