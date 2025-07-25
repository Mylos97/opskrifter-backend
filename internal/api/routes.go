package api

import (
	"net/http"
	"opskrifter-backend/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func RegisterRoutes(r *chi.Mux) {

	r.Use(httprate.LimitByIP(20, 1*time.Minute))
	r.Use(middleware.APIKeyAuth)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/recipes", func(r chi.Router) {
		r.Post("/", CreateRecipe)
		r.Get("/{id}", GetRecipe)
		r.Put("/{id}", UpdateRecipe)
		r.Delete("/{id}", DeleteRecipe)
	})
}
