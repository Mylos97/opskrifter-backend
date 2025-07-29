package api

import (
	"net/http"
	"opskrifter-backend/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
)

func RegisterRoutes(r *chi.Mux, env string) {

	if env == "prod" {
		r.Use(httprate.LimitByIP(20, 1*time.Minute))
		r.Use(middleware.APIKeyAuth)
	}

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	setupRouter(r)
}

func setupRouter(r *chi.Mux) {
	r.Route("/recipes", func(r chi.Router) {
		r.Post("/", CreateRecipe)
		r.Get("/{id}", GetRecipe)
		r.Get("/", GetManyRecipe)
		r.Put("/", UpdateRecipe)
		r.Delete("/{id}", DeleteRecipe)

		r.Post("/{id}/like", LikeRecipe)
		r.Delete("/{id}/like", UnlikeRecipe)
		r.Post("/{id}/views", UpdateViewRecipe)
	})

	r.Route("/ingredients", func(r chi.Router) {
		r.Get("/", GetManyIngredients)
	})
}
