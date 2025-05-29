package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(r *chi.Mux) {
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	r.Route("/recipes", func(r chi.Router) {
		r.Post("/", CreateRecipe)
		r.Get("/{id}", GetRecipe)
		r.Put("/{id}", UpdateRecipe)
		r.Delete("/{id}", DeleteRecipe)
		r.Get("/", GetRecipesHandler)
	 })

	r.Route("/cookbooks", func(r chi.Router) { 
		r.Post("/", CreateCookBook)
		r.Get("/{id}", GetCookBook)
		r.Put("/{id}", UpdateCookBook)
		r.Delete("/{id}", DeleteCookBook)
 })
}
