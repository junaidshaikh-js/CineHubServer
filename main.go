package main

import (
	"log"
	"net/http"
	"time"

	"github.com/joho/godotenv"
	"github.com/junaidshaikh-js/CineHubServer/handlers"
	"github.com/junaidshaikh-js/CineHubServer/logger"
	"github.com/junaidshaikh-js/CineHubServer/store"
)

func initializeLogger() *logger.Logger {
	logger, err := logger.NewLogger("movie.log")

	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	return logger
}

func main() {
	logger := initializeLogger()

	// Environment variables
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Failed to load environment variables: ", err)
	}

	DB, err := store.Open()

	if err != nil {
		log.Fatal("Failed to open database connection: ", err)
	}

	defer DB.Close()

	movieStore := store.NewPostgresMovieStore(DB)
	movieHandler := handlers.NewMovieHandler(movieStore, logger)

	accountStore := store.NewPostgresAccountStore(DB)
	accountHandler := handlers.NewAccountHandler(accountStore, logger)

	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/random", movieHandler.GetRandomMovies)
	http.HandleFunc("/api/movies/", movieHandler.GetMovieByID)
	http.HandleFunc("/api/movies/search", movieHandler.SearchMoviesByName)
	http.HandleFunc("/api/genres", movieHandler.GetGenres)

	http.HandleFunc("/api/account/register", accountHandler.Register)
	http.HandleFunc("/api/account/authenticate", accountHandler.Authenticate)

	http.Handle("/api/account/favorites",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetFavorites)))

	http.Handle("/api/account/watchlist",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.GetWatchlist)))

	http.Handle("/api/account/save-to-collection",
		accountHandler.AuthMiddleware(http.HandlerFunc(accountHandler.SaveToCollection)))

	catchAllClientRouteHandler := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./public/index.html")
	}

	http.HandleFunc("/movies", catchAllClientRouteHandler)
	http.HandleFunc("/movies/", catchAllClientRouteHandler)
	http.HandleFunc("/account", catchAllClientRouteHandler)
	http.HandleFunc("/account/", catchAllClientRouteHandler)

	http.Handle("/", http.FileServer(http.Dir("public")))

	const addr = ":5555"

	s := http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", addr)

	err = s.ListenAndServe()

	if err != nil {
		logger.Error("Server failed to start: ", err)
	}
}
