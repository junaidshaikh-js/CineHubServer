package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/junaidshaikh-js/CineHubServer/logger"
	"github.com/junaidshaikh-js/CineHubServer/store"
)

type MovieHandler struct {
	movieStore store.MovieStore
	logger     *logger.Logger
}

func NewMovieHandler(movieStore store.MovieStore, logger *logger.Logger) *MovieHandler {
	return &MovieHandler{
		movieStore: movieStore,
		logger:     logger,
	}
}

func (h *MovieHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		h.logger.Error("Failed to encode JSON", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.movieStore.GetTopMovies()

	if err != nil {
		h.logger.Error("Failed to get top movies", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, movies)
}

func (h *MovieHandler) GetRandomMovies(w http.ResponseWriter, r *http.Request) {
	movies, err := h.movieStore.GetRandomMovies()

	if err != nil {
		h.logger.Error("Failed to get random movies", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	h.writeJSON(w, movies)
}
