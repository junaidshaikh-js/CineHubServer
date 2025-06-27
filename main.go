package main

import (
	"log"
	"net/http"
	"time"

	"github.com/junaidshaikh-js/CineHubServer/handlers"
	"github.com/junaidshaikh-js/CineHubServer/logger"
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

	movieHandler := handlers.NewMovieHandler()

	http.HandleFunc("/api/movies/top", movieHandler.GetTopMovies)
	http.HandleFunc("/api/movies/random", movieHandler.GetRandomMovie)

	http.Handle("/", http.FileServer(http.Dir("public")))

	const addr = ":5555"

	s := http.Server{
		Addr:         addr,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("Starting server on %s", addr)

	err := s.ListenAndServe()

	if err != nil {
		logger.Error("Server failed to start: ", err)
	}
}
