package main

import (
	"log"
	"net/http"
	"github.com/go-chi/chi/v5"
	"opskrifter-backend/pkg/db"
	"opskrifter-backend/internal/api"
)
//"opskrifter-backend/internal/middleware"

func main() {
	db.Init()

	r := chi.NewRouter()
	//r.Use(middleware.GoogleAuthMiddleware("YOUR_GOOGLE_CLIENT_ID.apps.googleusercontent.com"))
	api.RegisterRoutes(r)

	log.Println("API server running on http://localhost:8080")
	http.ListenAndServe(":8080", r)
}
