package middleware

import (
	"net/http"
	"strings"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
)

func AuthMiddleware(apiKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")

			if !strings.HasPrefix(authHeader, "Bearer") {
				responses.RespondWithError(w, 400, "Authorization doesn't contain valid Bearer string")
				return 
			}

			parts := strings.Split(authHeader, " ")

			if len(parts) < 2 || parts[1] == "" {
				responses.RespondWithError(w, 400, "Invalid Bearer token format")
				return
			}

			token := strings.TrimSpace(parts[1])

			if token != apiKey {
				responses.RespondWithError(w, 400, "Unauthorized")
				return 
			}
			
			next.ServeHTTP(w, r)
		})
	}
}