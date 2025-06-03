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

	r.Route("/comments", func(r chi.Router) {
		r.Post("/", CreateComment)
		r.Get("/{recipe_id}", GetComments)
		r.Put("/{id}", UpdateComment)
		r.Delete("/{id}", DeleteComment)
	})

	r.Route("/users", func(r chi.Router) {
		r.Post("/", CreateUser)
		r.Get("/{id}", GetUser)
		r.Put("/{id}", UpdateUser)
		r.Delete("/{id}", DeleteUser)
	})

	r.Route("/like_recipe", func(r chi.Router) {
		r.Put("/like", LikeRecipe)
		r.Put("/unlike", UnLikeRecipe)
		r.Get("/{recipe_id}", GetLikeDRecipes)
	})
}
