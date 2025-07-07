package handlers

import (
	"encoding/json"
	"net/http"
	// "fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/vivekmodar03/go-movies-crud/internal/app/db"
	"github.com/vivekmodar03/go-movies-crud/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"time"
	// "os"
)

// POST /register
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Error hashing password: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Save user to DB
	query := `INSERT INTO user (username, email, password) VALUES (?, ?, ?)`
	_, err = db.DB.Exec(query, user.Username, user.Email, hashedPassword)
	if err != nil {
		http.Error(w, "Error saving user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}



var jwtKey = []byte("jwt@123")

func LoginUser(w http.ResponseWriter, r *http.Request) {
    var creds model.User
    err := json.NewDecoder(r.Body).Decode(&creds)
    if err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }

    var stored model.User
    err = db.DB.QueryRow("SELECT id, email, password FROM user WHERE email = ?", creds.Email).
        Scan(&stored.ID, &stored.Email, &stored.Password)

    if err != nil {
        http.Error(w, "Email not found. Register Please..", http.StatusUnauthorized)
        return
    }

    err = bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(creds.Password))
    if err != nil {
        http.Error(w, "Invalid password", http.StatusUnauthorized)
        return
    }

    // Generate access token (15-minute expiry)
    accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": stored.ID,
        "exp":     time.Now().Add(15 * time.Minute).Unix(),
    })
    accessTokenString, err := accessToken.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error generating access token", http.StatusInternalServerError)
        return
    }

    // Generate refresh token (7-day expiry)
    refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": stored.ID,
        "exp":     time.Now().Add(7 * 24 * time.Hour).Unix(),
    })
    refreshTokenString, err := refreshToken.SignedString(jwtKey)
    if err != nil {
        http.Error(w, "Error generating refresh token", http.StatusInternalServerError)
        return
    }

    // Store refresh token in DB (for invalidation)
    _, err = db.DB.Exec(
        "INSERT INTO user_tokens (user_id, token, expires_at) VALUES (?, ?, ?)",
        stored.ID, refreshTokenString, time.Now().Add(7*24*time.Hour),
    )
    if err != nil {
        http.Error(w, "Failed to save refresh token", http.StatusInternalServerError)
        return
    }

    // Set access token cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "access_token",
        Value:    accessTokenString,
        HttpOnly: true,
        Secure:   true, // HTTPS only
        SameSite: http.SameSiteStrictMode,
        Expires:  time.Now().Add(15 * time.Minute),
        Path:     "/", // Accessible across all paths
    })

    // Set refresh token cookie
    http.SetCookie(w, &http.Cookie{
        Name:     "refresh_token",
        Value:    refreshTokenString,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
        Expires:  time.Now().Add(7 * 24 * time.Hour),
        Path:     "/",
    })

    // Return success message
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]interface{}{
        "message": "Logged in successfully",
        // Tokens are now in cookies, not in response body
    })
}


func RefreshToken(w http.ResponseWriter, r *http.Request) {
    var req struct { RefreshToken string `json:"refresh_token"` }
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }

    // Validate refresh token
    token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
        return jwtKey, nil
    })
    if err != nil || !token.Valid {
        http.Error(w, "Invalid refresh token", http.StatusUnauthorized)
        return
    }

    // Check if token exists in DB (not revoked)
    var exists bool
    err = db.DB.QueryRow(
        "SELECT EXISTS(SELECT 1 FROM user_tokens WHERE token = ? AND expires_at > NOW())",
        req.RefreshToken,
    ).Scan(&exists)
    if err != nil || !exists {
        http.Error(w, "Refresh token revoked/expired", http.StatusUnauthorized)
        return
    }

    // Issue new access token
    claims := token.Claims.(jwt.MapClaims)
    newAccessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": claims["user_id"],
        "exp":     time.Now().Add(15 * time.Minute).Unix(),
    })
    newAccessTokenString, _ := newAccessToken.SignedString(jwtKey)

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newAccessTokenString,
    })
}



func LogoutUser(w http.ResponseWriter, r *http.Request) {
    // Get refresh token from cookie
    cookie, err := r.Cookie("refresh_token")
    if err != nil {
        http.Error(w, "Not logged in", http.StatusBadRequest)
        return
    }

    // Delete refresh token from DB
    _, err = db.DB.Exec("DELETE FROM user_tokens WHERE token = ?", cookie.Value)
    if err != nil {
        http.Error(w, "Logout failed", http.StatusInternalServerError)
        return
    }

    // Clear cookies
    http.SetCookie(w, &http.Cookie{
        Name:    "access_token",
        Value:   "",
        MaxAge:  -1,
    })
    http.SetCookie(w, &http.Cookie{
        Name:    "refresh_token",
        Value:   "",
        MaxAge:  -1,
    })

    w.WriteHeader(http.StatusOK)
}