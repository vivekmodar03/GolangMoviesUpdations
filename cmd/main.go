// package main

// import (
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"math/rand"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// )

// type Movie struct {
// 	ID string   `json:"id"`
// 	Isbn     string   `json:"isbn"`
// 	Title    string   `json:"title"`
// 	Director *Director `json:"director"`
// }

// type Director struct {
// 	Firstname string `json:"firstname"`
// 	Lastname  string `json:"lastname"`
// }

// var movies []Movie

// func getMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(movies)
// }

// func deleteMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for index, item := range movies {
// 		if item.ID == params["id"] {
// 			movies = append(movies[:index], movies[index+1:]...)
// 			break
// 		}
// 	}
// 	json.NewEncoder(w).Encode(movies)
// }

// func getMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	for _, item := range movies {
// 		if item.ID == params["id"] {
// 			json.NewEncoder(w).Encode(item)
// 			return
// 		}
// 	}
// }

// func createMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	var movie Movie
// 	_ = json.NewDecoder(r.Body).Decode(&movie)
// 	movie.ID = strconv.Itoa(rand.Intn(100000000))
// 	movies = append(movies, movie)
// 	json.NewEncoder(w).Encode(movie)
// }

// func updateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)

// 	for index, item := range movies {
// 		if item.ID == params["id"] {
// 			movies = append(movies[:index], movies[index+1:]...)
// 			var movie Movie
// 			_ = json.NewDecoder(r.Body).Decode(&movie)
// 			movie.ID = params["id"]
// 			movies = append(movies, movie)
// 			json.NewEncoder(w).Encode(movie)
// 			return
// 		}
// 	}
// }

// func main() {
// 	r := mux.NewRouter()

// 	movies = append(movies, Movie{ID: "1", Isbn: "438227", Title: "Movie One", Director: &Director{Firstname: "John", Lastname: "Doe"}})
// 	movies = append(movies, Movie{ID: "2", Isbn: "454555", Title: "Movie Two", Director: &Director{Firstname: "Steve", Lastname: "Smith"}})

// 	r.HandleFunc("/movies", getMovies).Methods("GET")
// 	r.HandleFunc("/movies/{id}", getMovie).Methods("GET")
// 	r.HandleFunc("/movies", createMovie).Methods("POST")
// 	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
// 	r.HandleFunc("/movies/{id}", deleteMovie).Methods("DELETE")

// 	fmt.Println("Starting server at port 8080...")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }



// package main

// import (
// 	"database/sql"
// 	"encoding/json"
// 	"fmt"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gorilla/mux"
// 	_ "github.com/go-sql-driver/mysql"
// )

// // Movie struct for incoming JSON
// type Movie struct {
// 	ID       int      `json:"id"`
// 	Isbn     string   `json:"isbn"`
// 	Title    string   `json:"title"`
// 	Director Director `json:"director"`
// }

// type Director struct {
// 	Firstname string `json:"firstname"`
// 	Lastname  string `json:"lastname"`
// }

// // Global DB connection
// var db *sql.DB

// func main() {
// 	var err error
// 	// 1. Connect to MySQL
// 	db, err = sql.Open("mysql", "root:system@0987@tcp(127.0.0.1:3305)/go_movie_db")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer db.Close()

// 	// Check connection
// 	if err = db.Ping(); err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Println("Connected to MySQL!")

// 	// 2. Setup router
// 	r := mux.NewRouter()
// 	r.HandleFunc("/movies", createMovie).Methods("POST")
// 	r.HandleFunc("/movies", getMovies).Methods("GET")
// 	r.HandleFunc("/movies/{id}", getMovieByID).Methods("GET")
// 	r.HandleFunc("/movies/{id}", updateMovie).Methods("PUT")
// 	r.HandleFunc("/movies/{id}", deleteMovieById).Methods("DELETE")
// 	r.HandleFunc("/movies/{id}", deleteAllMovies).Methods("DELETE")
	
	

// 	// 3. Start server
// 	fmt.Println("Server started at :8080")
// 	log.Fatal(http.ListenAndServe(":8080", r))
// }

// 4. Create movie handler
// func createMovie(w http.ResponseWriter, r *http.Request) {
// 	var movie Movie
// 	err := json.NewDecoder(r.Body).Decode(&movie)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// 5. Insert into MySQL
// 	query := `INSERT INTO movies (isbn, title, director_firstname, director_lastname) VALUES (?, ?, ?, ?)`
// 	result, err := db.Exec(query, movie.Isbn, movie.Title, movie.Director.Firstname, movie.Director.Lastname)
// 	if err != nil {
// 		http.Error(w, "Error inserting into database: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// 6. Get inserted ID and respond
// 	id, _ := result.LastInsertId()
// 	fmt.Fprintf(w, "Movie inserted with ID: %d", id)
// }


// func getMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	rows, err := db.Query("SELECT * FROM movies")
// 	if err != nil {
// 		http.Error(w, "Failed to fetch movies: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	defer rows.Close()

// 	var movies []Movie

// 	for rows.Next() {
// 		var movie Movie
// 		var firstname, lastname string

// 		err := rows.Scan(&movie.ID, &movie.Isbn, &movie.Title, &firstname, &lastname)
// 		if err != nil {
// 			http.Error(w, "Error scanning row: "+err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		movie.Director = Director{
// 			Firstname: firstname,
// 			Lastname:  lastname,
// 		}

// 		movies = append(movies, movie)
// 	}

// 	json.NewEncoder(w).Encode(movies)
// }


// func getMovieByID(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	params := mux.Vars(r)
// 	id := params["id"]

// 	var movie Movie
// 	var firstname, lastname string

// 	query := "SELECT * FROM movies WHERE id = ?"
// 	err := db.QueryRow(query, id).Scan(&movie.ID, &movie.Isbn, &movie.Title, &firstname, &lastname)
// 	if err != nil {
// 		if err == sql.ErrNoRows {
// 			http.Error(w, "Movie not found", http.StatusNotFound)
// 		} else {
// 			http.Error(w, "Database error: "+err.Error(), http.StatusInternalServerError)
// 		}
// 		return
// 	}

// 	movie.Director = Director{
// 		Firstname: firstname,
// 		Lastname:  lastname,
// 	}

// 	json.NewEncoder(w).Encode(movie)
// }


// func updateMovie(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	var updatedMovie Movie
// 	err := json.NewDecoder(r.Body).Decode(&updatedMovie)
// 	if err != nil {
// 		http.Error(w, "Invalid JSON: "+err.Error(), http.StatusBadRequest)
// 		return
// 	}

// 	// Execute UPDATE query
// 	query := `UPDATE movies SET isbn = ?, title = ?, director_firstname = ?, director_lastname = ? WHERE id = ?`
// 	result, err := db.Exec(query, updatedMovie.Isbn, updatedMovie.Title, updatedMovie.Director.Firstname, updatedMovie.Director.Lastname, id)
// 	if err != nil {
// 		http.Error(w, "Failed to update movie: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Check if a row was actually updated
// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "No movie found with the given ID", http.StatusNotFound)
// 		return
// 	}

// 	// Set the updated ID to return in response
// 	updatedMovie.ID, _ = strconv.Atoi(id)
// 	json.NewEncoder(w).Encode(updatedMovie)
// }


// func deleteMovieById(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")
// 	params := mux.Vars(r)
// 	id := params["id"]

// 	query := "DELETE FROM movies WHERE id = ?"
// 	result, err := db.Exec(query, id)
// 	if err != nil {
// 		http.Error(w, "Failed to delete movie: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	rowsAffected, _ := result.RowsAffected()
// 	if rowsAffected == 0 {
// 		http.Error(w, "No movie found with the given ID", http.StatusNotFound)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "Movie deleted successfully"})
// }


// func deleteAllMovies(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "application/json")

// 	query := "DELETE FROM movies"
// 	_, err := db.Exec(query)
// 	if err != nil {
// 		http.Error(w, "Failed to delete all movies: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	json.NewEncoder(w).Encode(map[string]string{"message": "All movies deleted successfully"})
// }


package main

import (
	"log"
	"net/http"
	"github.com/vivekmodar03/go-movies-crud/internal/app/db"
	"github.com/vivekmodar03/go-movies-crud/internal/routes"
)

func main() {
	db.Init() // Connect to MySQL
	r := Routes.SetupRouter() // Register routes
	log.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
