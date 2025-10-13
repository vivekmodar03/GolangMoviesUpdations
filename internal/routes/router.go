package routes

import (
	"github.com/gorilla/mux"
	"github.com/vivekmodar03/go-movies-crud/internal/handlers"
	"github.com/vivekmodar03/go-movies-crud/internal/middleware"
)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// User registration and login routes
	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	// All movie routes are protected and require a valid Firebase ID token
	r.HandleFunc("/movies", middleware.Auth(handlers.CreateMovie)).Methods("POST")
	r.HandleFunc("/movies", middleware.Auth(handlers.GetMovies)).Methods("GET") // Add this line
	r.HandleFunc("/movies/{id}", middleware.Auth(handlers.GetMovieByID)).Methods("GET")
	r.HandleFunc("/movies/{id}", middleware.Auth(handlers.UpdateMovie)).Methods("PUT")
	r.HandleFunc("/movies/{id}", middleware.Auth(handlers.DeleteMovieById)).Methods("DELETE")
	r.HandleFunc("/movies", middleware.Auth(handlers.DeleteAllMovies)).Methods("DELETE")

	return r
}