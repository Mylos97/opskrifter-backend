package main

import (
	"log"
	"net/http"
	"opskrifter-backend/internal/api"
	"opskrifter-backend/internal/middleware"
	"opskrifter-backend/pkg/db"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")

	isProd := port != "8080"
	db.Init(isProd)

	if _, err := db.DB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		log.Fatalf("failed to enable foreign key constraints: %v", err)
	}

	r := chi.NewRouter()

	if isProd {
		r.Use(middleware.APIKeyAuth)
		log.Println("Starting in PRODUCTION mode...")
	} else {
		log.Println("Starting in LOCAL mode...")
	}

	api.RegisterRoutes(r)

	log.Printf("API server running on http://localhost%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
