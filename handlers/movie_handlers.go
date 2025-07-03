package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vivekmodar03/go-movies-crud/db"
	"github.com/vivekmodar03/go-movies-crud/model"
)

// POST /movies - Add a new movie
func CreateMovie(w http.ResponseWriter, r *http.Request) {
	var movie model.Movie

	err := json.NewDecoder(r.Body).Decode(&movie)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO movies (isbn, title, director_firstname, director_lastname) VALUES (?, ?, ?, ?)`
	result, err := db.DB.Exec(query, movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname)
	if err != nil {
		http.Error(w, "Error inserting into database: "+err.Error(), http.StatusInternalServerError)
		return
	}

	id, _ := result.LastInsertId()
	fmt.Fprintf(w, "Movie inserted with ID: %d", id)
	json.NewEncoder(w).Encode(movie)

}

// GET /movies - Get all movies
func GetMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	rows, err := db.DB.Query("SELECT * FROM movies")
	if err != nil {
		http.Error(w, "Failed to fetch movies: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var movies []model.Movie

	for rows.Next() {
		var movie model.Movie
		var firstname, lastname string

		err := rows.Scan(&movie.ID, &movie.Isbn, &movie.Title, &firstname, &lastname)
		if err != nil {
			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
			return
		}

		movie.Director = model.Director{
			Firstname: firstname,
			Lastname:  lastname,
		}

		movies = append(movies, movie)
	}

	json.NewEncoder(w).Encode(movies)
}

// GET /movies/{id} - Get a movie by ID
func GetMovieByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	params := mux.Vars(r)
	id := params["id"]

	var movie model.Movie
	var firstname, lastname string

	query := "SELECT * FROM movies WHERE id = ?"
	err := db.DB.QueryRow(query, id).Scan(&movie.ID, &movie.Isbn, &movie.Title, &firstname, &lastname)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Movie not found", http.StatusNotFound)
		} else {
			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
		}
		return
	}

	movie.Director = model.Director{
		Firstname: firstname,
		Lastname:  lastname,
	}

	json.NewEncoder(w).Encode(movie)
}

// PUT /movies/{id} - Update a movie by ID
func UpdateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	var updatedMovie model.Movie
	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
	if err != nil {
		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `UPDATE movies SET isbn = ?, title = ?, director_firstname = ?, director_lastname = ? WHERE id = ?`
	result, err := db.DB.Exec(query, updatedMovie.Isbn, updatedMovie.Title, updatedMovie.Director.Firstname, updatedMovie.Director.Lastname, id)
	if err != nil {
		http.Error(w, "Failed to update movie: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No movie found with the given ID", http.StatusNotFound)
		return
	}

	updatedMovie.ID, _ = strconv.Atoi(id)
	json.NewEncoder(w).Encode(updatedMovie)
}

// DELETE /movies/{id} - Delete a movie by ID
func DeleteMovieById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	id := params["id"]

	query := "DELETE FROM movies WHERE id = ?"
	result, err := db.DB.Exec(query, id)
	if err != nil {
		http.Error(w, "Failed to delete movie: "+err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		http.Error(w, "No movie found with the given ID", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
}

// DELETE /movies - Delete all movies
func DeleteAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	query := "DELETE FROM movies"
	_, err := db.DB.Exec(query)
	if err != nil {
		http.Error(w, "Failed to delete all movies: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "All movies deleted successfully"})
}
