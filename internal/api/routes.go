package api

import (
	"net/http"
	"github.com/go-chi/chi/v5"
)

// RegisterRoutes registers all API routes on the router.
func RegisterRoutes(r *chi.Mux) {
	// Example route, replace or add your actual routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	// Add your other route registrations here, e.g.:
	//r.Route("/recipes", func(r chi.Router) { ... })
	// r.Route("/cookbooks", func(r chi.Router) { ... })
}
