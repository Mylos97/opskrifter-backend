package main

import (
	"log"
	"net/http"
	"opskrifter-backend/internal/api"
	"opskrifter-backend/pkg/myDB"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")

	err := myDB.Init(false)
	if err != nil {
		log.Fatalf("error init DB %v", err)
	}
	r := chi.NewRouter()

	api.RegisterRoutes(r)

	log.Printf("API server running on http://localhost%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
