package main

import (
	"log"
	"net/http"
	"opskrifter-backend/internal/api"
	"opskrifter-backend/pkg/db"

	"github.com/go-chi/chi/v5"
)

//"opskrifter-backend/internal/middleware"

func main() {
	db.Init(false)

	r := chi.NewRouter()
	api.RegisterRoutes(r)

	log.Println("API server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
