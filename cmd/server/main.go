package main

import (
	"log"
	"net/http"
	"opskrifter-backend/internal/api"
	"opskrifter-backend/internal/middleware"
	"opskrifter-backend/pkg/myDB"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")

	isProd := port != "8080"
	err := myDB.Init(isProd)
	if err != nil {
		log.Fatalf("error init DB %v", err)
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
