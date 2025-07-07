package middleware

import (
	"net/http"
	"strings"
	"github.com/golang-jwt/jwt/v5"
	"fmt"
	"context"
	
)

var jwtKey = []byte("jwt@123")

func Auth(next http.HandlerFunc) http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        // Get token from EITHER cookie OR Authorization header
        var tokenStr string
        
        // First try to get from cookie
        if cookie, err := r.Cookie("access_token"); err == nil {
            tokenStr = cookie.Value
        } else {
            // Fall back to Authorization header
            authHeader := r.Header.Get("Authorization")
            if authHeader == "" {
                http.Error(w, "Authorization required", http.StatusUnauthorized)
                return
            }

            parts := strings.Split(authHeader, " ")
            if len(parts) != 2 || parts[0] != "Bearer" {
                http.Error(w, "Authorization format: 'Bearer <token>'", http.StatusUnauthorized)
                return
            }
            tokenStr = parts[1]
        }

        // Parse the token
        token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
            }
            return jwtKey, nil
        })

        if err != nil {
            http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
            return
        }

        if !token.Valid {
            http.Error(w, "Token expired or invalid", http.StatusUnauthorized)
            return
        }

        // Pass claims to downstream handlers
        if claims, ok := token.Claims.(jwt.MapClaims); ok {
            ctx := context.WithValue(r.Context(), "user_id", claims["user_id"])
            r = r.WithContext(ctx)
        }

        next.ServeHTTP(w, r)
    }
}