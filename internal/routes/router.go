package Routes

import (
	// "github.com/vivekmodar03/go-movies-crud/model"
	"github.com/vivekmodar03/go-movies-crud/internal/handlers"
	"github.com/gorilla/mux"
	"github.com/vivekmodar03/go-movies-crud/internal/middleware"
	// "net/http"

)

func SetupRouter() *mux.Router {
	r := mux.NewRouter()

	// BEFORE AUTHENTICATION
	// r.HandleFunc("/movies", handlers.GetMovies).Methods("GET")
	// r.HandleFunc("/movies/{id}", handlers.GetMovieByID).Methods("GET")
	// r.HandleFunc("/movies", handlers.CreateMovie).Methods("POST")
	// r.HandleFunc("/movies/{id}", handlers.UpdateMovie).Methods("PUT")
	// r.HandleFunc("/movies/{id}", handlers.DeleteMovieById).Methods("DELETE")
	// r.HandleFunc("/movies", handlers.DeleteAllMovies).Methods("DELETE")

// Public routes (no auth required)
r.HandleFunc("/register", handlers.RegisterUser).Methods("POST")
r.HandleFunc("/login", handlers.LoginUser).Methods("POST")

// Protected routes (require valid JWT)
r.HandleFunc("/movies", middleware.Auth(handlers.GetMovies)).Methods("GET")
r.HandleFunc("/movies", middleware.Auth(handlers.CreateMovie)).Methods("POST")
r.HandleFunc("/movies/{id}", middleware.Auth(handlers.GetMovieByID)).Methods("GET")
r.HandleFunc("/movies/{id}", middleware.Auth(handlers.UpdateMovie)).Methods("PUT")
r.HandleFunc("/movies/{id}", middleware.Auth(handlers.DeleteMovieById)).Methods("DELETE")

// Token management
r.HandleFunc("/refresh-token", handlers.RefreshToken).Methods("POST")  // New
r.HandleFunc("/logout", middleware.Auth(handlers.LogoutUser)).Methods("POST")  // New

	return r
}
