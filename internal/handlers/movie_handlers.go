// package handlers

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// 	"github.com/vivekmodar03/go-movies-crud/internal/app/db"
// 	"github.com/vivekmodar03/go-movies-crud/internal/model"
// )

// // POST /movies - Add a new movie
// func CreateMovie(w http.ResponseWriter, r *http.Request) {
// 	var movie model.Movie

// 	err := json.NewDecoder(r.Body).Decode(&movie)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	query := `INSERT INTO movies (isbn, title, director_firstname, director_lastname) VALUES (?, ?, ?, ?)`
// 	result, err := db.DB.Exec(query, movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname)
// 	if err != nil {
// 		http.Error(w, "Error inserting into database: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	id, _ := result.LastInsertId()
// 	fmt.Fprintf(w, "Movie inserted with ID: %d", id)
// 	json.NewEncoder(w).Encode(movie)

// }

// // GET /movies - Get all movies
// func GetMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	rows, err := db.DB.Query("SELECT * FROM movies")
// 	if err != nil {
// 		http.Error(w, "Failed to fetch movies: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var movies []model.Movie

// 	for rows.Next() {
// 		var movie model.Movie
// 		var firstname, lastname string

// 		err := rows.Scan(&movie.ID, &movie.Isbn, &movie.Title, &firstname, &lastname)
// 		if err != nil {
// 			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		movie.Director = model.Director{
// 			Firstname: firstname,
// 			Lastname:  lastname,
// 		}

// 		movies = append(movies, movie)
// 	}

// 	json.NewEncoder(w).Encode(movies)
// }

// // GET /movies/{id} - Get a movie by ID
// func GetMovieByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	params := mux.Vars(r)
// 	id := params["id"]

// 	var movie model.Movie
// 	var firstname, lastname string

// 	query := "SELECT * FROM movies WHERE id = ?"
// 	err := db.DB.QueryRow(query, id).Scan(&movie.ID, &movie.Isbn, &movie.Title, &firstname, &lastname)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Movie not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	movie.Director = model.Director{
// 		Firstname: firstname,
// 		Lastname:  lastname,
// 	}

// 	json.NewEncoder(w).Encode(movie)
// }

// // PUT /movies/{id} - Update a movie by ID
// func UpdateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	var updatedMovie model.Movie
// 	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	query := `UPDATE movies SET isbn = ?, title = ?, director_firstname = ?, director_lastname = ? WHERE id = ?`
// 	result, err := db.DB.Exec(query, updatedMovie.Isbn, updatedMovie.Title, updatedMovie.Director.Firstname, updatedMovie.Director.Lastname, id)
// 	if err != nil {
// 		http.Error(w, "Failed to update movie: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "No movie found with the given ID", http.StatusNotFound)
// 		return
// 	}

// 	updatedMovie.ID, _ = strconv.Atoi(id)
// 	json.NewEncoder(w).Encode(updatedMovie)
// }

// // DELETE /movies/{id} - Delete a movie by ID
// func DeleteMovieById(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	tx, err := db.DB.Begin()
// 	if err != nil {
// 		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Rollback() // Rollback on any error

// 	// Delete the movie
// 	result, err := tx.Exec("DELETE FROM movies WHERE id = ?", id)
// 	if err != nil {
// 		http.Error(w, "Failed to delete movie: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "No movie found with the given ID", http.StatusNotFound)
// 		return
// 	}

// 	// Re-number the IDs
// 	_, err = tx.Exec("SET @count = 0;")
// 	if err != nil {
// 		http.Error(w, "Failed to renumber IDs: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	_, err = tx.Exec("UPDATE movies SET id = @count:= @count + 1 ORDER BY id;")
// 	if err != nil {
// 		http.Error(w, "Failed to renumber IDs: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Reset auto-increment
// 	var maxID int
// 	err = tx.QueryRow("SELECT MAX(id) FROM movies").Scan(&maxID)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			maxID = 0
// 		} else {
// 			http.Error(w, "Error getting max ID: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 	}

// 	// Correctly format the ALTER TABLE query
// 	query := fmt.Sprintf("ALTER TABLE movies AUTO_INCREMENT = %d", maxID+1)
// 	_, err = tx.Exec(query)
// 	if err != nil {
// 		http.Error(w, "Failed to reset auto-increment: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := tx.Commit(); err != nil {
// 		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully and IDs have been re-ordered"})
// }

// // DELETE /movies - Delete all movies
// func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	tx, err := db.DB.Begin()
// 	if err != nil {
// 		http.Error(w, "Failed to begin transaction: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer tx.Rollback()

// 	_, err = tx.Exec("DELETE FROM movies")
// 	if err != nil {
// 		http.Error(w, "Failed to delete all movies: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	_, err = tx.Exec("ALTER TABLE movies AUTO_INCREMENT = 1")
// 	if err != nil {
// 		http.Error(w, "Failed to reset auto-increment: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	if err := tx.Commit(); err != nil {
// 		http.Error(w, "Failed to commit transaction: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "All movies deleted successfully"})
// }


package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var movie model.Movie
	if err := json.NewDecoder(r.Body).Decode(&movie); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO movies (isbn, title, director_firstname, director_lastname, user_id) VALUES (?, ?, ?, ?, ?)`
	result, err := db.DB.Exec(query, movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname, userUID)
	if err != nil {
		http.Error(w, "Error inserting into database: "+err.Error(), http.StatusInternalServerError)
		return
	}
	
	id, _ := result.LastInsertId()
	movie.ID = int(id)
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