package Routes

import (
	"github.com/vivekmodar03/go-movies-crud/internal/handlers"
	"github.com/gorilla/mux"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// Public routes for movie CRUD operations
	r.HandleFunc("/movies", handlers.GetMovies).Methods("GET")
	r.HandleFunc("/movies/{id}", handlers.GetMovieByID).Methods("GET")
	r.HandleFunc("/movies", handlers.CreateMovie).Methods("POST")
	r.HandleFunc("/movies/{id}", handlers.UpdateMovie).Methods("PUT")
	r.HandleFunc("/movies/{id}", handlers.DeleteMovieById).Methods("DELETE")
	r.HandleFunc("/movies", handlers.DeleteAllMovies).Methods("DELETE")

	return r
}