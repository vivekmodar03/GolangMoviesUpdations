package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	// "errors"
	"io/ioutil"
	"net/http"
	// "os"
	// "strconv"

	"firebase.google.com/go/v4/auth"
	// "github.com/gorilla/mux"
	// "github.com/vivekmodar03/go-movies-crud/internal/app/db"
	"github.com/vivekmodar03/go-movies-crud/internal/app/firebase"
	"github.com/vivekmodar03/go-movies-crud/internal/model"
)

// POST /register - Register a new user
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user model.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	params := (&auth.UserToCreate{}).
		Email(user.Email).
		Password(user.Password)

	u, err := firebase.AuthClient.CreateUser(context.Background(), params)
	if err != nil {
		http.Error(w, "Error creating user: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully", "uid": u.UID})
}


// POST /login - Log in a user
func LoginUser(w http.ResponseWriter, r *http.Request) {
	email, password, ok := r.BasicAuth()
	if !ok {
		http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
		return
	}

	// You need to replace this with your Firebase Web API Key
	// Directly assign the key to the variable
		apiKey := "AIzaSyAQgAtVvpk0ui3x8WCZGVtXtyO28fYwlbI"

		if apiKey == "" { // This check is a bit redundant now but good practice
			http.Error(w, "Firebase API key is not set", http.StatusInternalServerError)
			return
		}

	requestBody, err := json.Marshal(map[string]interface{}{
		"email":             email,
		"password":          password,
		"returnSecureToken": true,
	})

	if err != nil {
		http.Error(w, "Error creating request body", http.StatusInternalServerError)
		return
	}

	resp, err := http.Post("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key="+apiKey, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		http.Error(w, "Error logging in: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Error reading response body", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "Login failed: "+string(body), resp.StatusCode)
		return
	}

	var result map[string]interface{}
	json.Unmarshal(body, &result)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}