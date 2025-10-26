package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/Giammarco-Ferranti/image-metadata-api/pkg/responses"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := os.Getenv("API_KEY")

		authHeader := r.Header.Get("Authorization")


		if !strings.HasPrefix(authHeader, "Bearer") {
			responses.RespondWithError(w, 400, "Authorization doesn't contain valid Bearer string")
			return 
		}

		providedApiKey := strings.SplitAfter(authHeader, " ")


		if providedApiKey[1] == "" || providedApiKey[1] != apiKey {
			responses.RespondWithError(w, 400, "Unauthorized")
			return 
		}
		
		next.ServeHTTP(w, r)
	})
}