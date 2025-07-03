package middleware

import (
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"


)

var jwtKey = []byte("jwt@123")

func Auth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.Split(authHeader, "Bearer ")
		if len(parts) != 2 {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		tokenStr := parts[1]

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		next.ServeHTTP(w, r)
	}
}

