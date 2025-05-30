package main

import (
	"log"
	"net/http"
	"opskrifter-backend/internal/api"
	"opskrifter-backend/internal/middleware"
	"opskrifter-backend/pkg/db"

	"github.com/go-chi/chi/v5"
)

func main() {
	db.Init(false)

	r := chi.NewRouter()
	r.Use(middleware.APIKeyAuth)
	api.RegisterRoutes(r)
	log.Println("API server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
