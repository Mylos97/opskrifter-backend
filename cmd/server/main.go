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
	env := os.Getenv("ENV")

	isProd := env == "prod"
	db.Init(isProd)

	r := chi.NewRouter()

	if isProd {
		r.Use(middleware.APIKeyAuth)
		log.Println("Starting in PRODUCTION mode...")
	} else {
		log.Println("Starting in LOCAL mode...")
	}

	api.RegisterRoutes(r)

	addr := ":8080"
	if isProd {
		addr = ":80"
	}

	log.Printf("API server running on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
