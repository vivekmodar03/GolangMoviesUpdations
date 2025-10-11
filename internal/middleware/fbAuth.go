package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/vivekmodar03/go-movies-crud/internal/app/firebase"
)

// Auth middleware verifies the Firebase ID token
func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header must be provided", http.StatusUnauthorized)
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || strings.ToLower(tokenParts[0]) != "bearer" {
			http.Error(w, "Authorization header format must be 'Bearer {token}'", http.StatusUnauthorized)
			return
		}

		idToken := tokenParts[1]
		token, err := firebase.AuthClient.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user's Firebase UID to the request context
		ctx := context.WithValue(r.Context(), "userUID", token.UID)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}