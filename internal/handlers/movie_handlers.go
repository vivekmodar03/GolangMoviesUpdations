
package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"log"

	"github.com/gorilla/mux"
	"github.com/vivekmodar03/go-movies-crud/internal/app/db"
	"github.com/vivekmodar03/go-movies-crud/internal/model"
)

// Helper to get user UID from context
func getUserUIDFromContext(r *http.Request) (string, error) {
	userUID, ok := r.Context().Value("userUID").(string)
	if !ok || userUID == "" {
		return "", errors.New("could not retrieve user UID from context")
	}
	return userUID, nil
}

// POST /movies - Add a new movie for the logged-in user
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	userUID, err := getUserUIDFromContext(r)
	if err != nil {
		log.Println("Error getting user UID from context:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		log.Println("Error decoding request body:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	movie.UserID = userUID

	// 1. Begin a transaction
	tx, err := db.DB.Begin()
	if err != nil {
		log.Println("Failed to begin transaction:", err)
		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	query := `INSERT INTO movies (isbn, title, director_firstname, director_lastname, user_id) VALUES (?, ?, ?, ?, ?)`
	
	// 2. Execute the query on the transaction (tx), not the direct database connection (db.DB)
	result, err := tx.Exec(query, movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname, movie.UserID)
	if err != nil {
		// If there's an error, roll back the transaction
		tx.Rollback()
		log.Println("Error inserting into database:", err)
		http.Error(w, "Error inserting into database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// 3. Commit the transaction to permanently save the data
	if err := tx.Commit(); err != nil {
		log.Println("Failed to commit transaction:", err)
		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	movie.ID = int(id)
	log.Println("Successfully inserted and committed movie with ID:", id) // Updated log message
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(movie)
}

// GET /movies - Get all movies for the logged-in user
func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userUID, err := getUserUIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := db.DB.Query("SELECT id, isbn, title, director_firstname, director_lastname FROM movies WHERE user_id = ?", userUID)
	if err != nil {
		http.Error(w, "Failed to fetch movies: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []model.Movie
	for rows.Next() {
		var movie model.Movie
		err := rows.Scan(&movie.ID, &movie.Isbn, &movie.Title, &movie.Director.Firstname, &movie.Director.Lastname)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}
		movies = append(movies, movie)
	}
	json.NewEncoder(w).Encode(movies)
}

// GET /movies/{id} - Get a specific movie owned by the logged-in user
func GetMovieByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userUID, err := getUserUIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := mux.Vars(r)["id"]
	var movie model.Movie
	query := "SELECT id, isbn, title, director_firstname, director_lastname FROM movies WHERE id = ? AND user_id = ?"
	err = db.DB.QueryRow(query, id, userUID).Scan(&movie.ID, &movie.Isbn, &movie.Title, &movie.Director.Firstname, &movie.Director.Lastname)
	
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found or you don't have permission", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}
	json.NewEncoder(w).Encode(movie)
}

// PUT /movies/{id} - Update a movie owned by the logged-in user
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userUID, err := getUserUIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	id := mux.Vars(r)["id"]
	var updatedMovie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&updatedMovie); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	query := `UPDATE movies SET isbn = ?, title = ?, director_firstname = ?, director_lastname = ? WHERE id = ? AND user_id = ?`
	result, err := db.DB.Exec(query, updatedMovie.Isbn, updatedMovie.Title, updatedMovie.Director.Firstname, updatedMovie.Director.Lastname, id, userUID)
	if err != nil {
		http.Error(w, "Failed to update movie: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No movie found with the given ID or you don't have permission", http.StatusNotFound)
		return
	}
	
	updatedMovie.ID, _ = strconv.Atoi(id)
	json.NewEncoder(w).Encode(updatedMovie)
}

// DELETE /movies/{id} - Delete a movie owned by the logged-in user
func DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userUID, err := getUserUIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	id := mux.Vars(r)["id"]
	query := "DELETE FROM movies WHERE id = ? AND user_id = ?"
	result, err := db.DB.Exec(query, id, userUID)
	if err != nil {
		http.Error(w, "Failed to delete movie: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No movie found with the given ID or you don't have permission", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
}

// DELETE /movies - Delete all movies for the logged-in user
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userUID, err := getUserUIDFromContext(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, err = db.DB.Exec("DELETE FROM movies WHERE user_id = ?", userUID)
	if err != nil {
		http.Error(w, "Failed to delete movies: "+err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"message": "All your movies have been deleted successfully"})
}