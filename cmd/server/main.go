package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"opskrifter-backend/internal/api"
	"opskrifter-backend/pkg/myDB"
	"os"

	"github.com/go-chi/chi/v5"
)

func main() {
	port := os.Getenv("PORT")
	env := flag.String("env", "dev", "Application environment: dev or prod")
	flag.Parse()
	fmt.Printf("Running in %s mode\n", *env)

	err := myDB.Init(false)
	if err != nil {
		log.Fatalf("error init DB %v", err)
	}
	r := chi.NewRouter()

	api.RegisterRoutes(r, *env)

	log.Printf("API server running on http://localhost%s", port)
	if err := http.ListenAndServe(port, r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
