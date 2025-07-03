package Routes

import (
	// "github.com/vivekmodar03/go-movies-crud/model"
	"github.com/vivekmodar03/go-movies-crud/handlers"
	"github.com/gorilla/mux"
	"github.com/vivekmodar03/go-movies-crud/middleware"
	// "net/http"

)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// r.HandleFunc("/movies", handlers.GetMovies).Methods("GET")
	// r.HandleFunc("/movies/{id}", handlers.GetMovieByID).Methods("GET")
	// r.HandleFunc("/movies", handlers.CreateMovie).Methods("POST")
	// r.HandleFunc("/movies/{id}", handlers.UpdateMovie).Methods("PUT")
	// r.HandleFunc("/movies/{id}", handlers.DeleteMovieById).Methods("DELETE")
	// r.HandleFunc("/movies", handlers.DeleteAllMovies).Methods("DELETE")

	r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
	r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

	r.HandleFunc("/movies", middleware.Auth(handlers.GetMovies)).Methods("GET")
    r.HandleFunc("/movies", middleware.Auth(handlers.CreateMovie)).Methods("POST")
    r.HandleFunc("/movies/{id}", middleware.Auth(handlers.GetMovieByID)).Methods("GET")
    r.HandleFunc("/movies/{id}", middleware.Auth(handlers.UpdateMovie)).Methods("PUT")
    r.HandleFunc("/movies/{id}", middleware.Auth(handlers.DeleteMovieById)).Methods("DELETE")
    r.HandleFunc("/movies", middleware.Auth(handlers.DeleteAllMovies)).Methods("DELETE")

	return r
}
