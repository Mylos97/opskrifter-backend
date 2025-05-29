package middleware

import (
	"context"
	"net/http"
	"strings"
	"google.golang.org/api/idtoken"
	"log"
)

func GoogleAuthMiddleware(audience string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
				return
			}
			idToken := parts[1]

			ctx := context.Background()
			payload, err := idtoken.Validate(ctx, idToken, audience)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			email, ok := payload.Claims["email"].(string)
			if ok {
				log.Println("Authenticated user:", email)
			}

			next.ServeHTTP(w, r)
		})
	}
}
