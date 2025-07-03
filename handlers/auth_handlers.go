package handlers

import (
	"encoding/json"
	"net/http"
	// "fmt"

	"golang.org/x/crypto/bcrypt"
	"github.com/vivekmodar03/go-movies-crud/db"
	"github.com/vivekmodar03/go-movies-crud/model"
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
	// fmt.Println("Login attempt with email:", creds.Email)

	err = db.DB.QueryRow("SELECT email, password FROM user WHERE email = ?", creds.Email).
    Scan(&stored.Email, &stored.Password)

	// fmt.Printf("SQL Query: SELECT id, email, password FROM users WHERE email = '%s'\n", creds.Email)


	if err != nil {
		http.Error(w, "Email not found. Register Please..", http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(stored.Password), []byte(creds.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": stored.ID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}


	json.NewEncoder(w).Encode(map[string]string{ "Message" : "Logged in Successfully.."	})

	json.NewEncoder(w).Encode(map[string]string{ "COPY THE TOKEN AND PASTE INTO AUTHORIZATION -> BEARER TOKEN"})

	json.NewEncoder(w).Encode(map[string]string{
		"token": tokenString,
	})
}

