package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/junaidshaikh-js/CineHubServer/models"
)

type MovieHandler struct{}

func NewMovieHandler() *MovieHandler {
	return &MovieHandler{}
}

func (h *MovieHandler) writeJSON(w http.ResponseWriter, data interface{}) {
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (h *MovieHandler) GetTopMovies(w http.ResponseWriter, r *http.Request) {
	movies := []models.Movie{
		{
			ID:          1,
			TMDB_ID:     101,
			Title:       "The Hacker",
			ReleaseYear: 2022,
			Genres: []models.Genre{{
				ID:   1,
				Name: "Thriller",
			}},
			Keywords: []string{
				"hacking",
				"cybercrime",
			},
			Casting: []models.Actor{{
				ID:        1,
				FirstName: "Jane",
				LastName:  "Doe",
			}},
		},
		{
			ID:          2,
			TMDB_ID:     102,
			Title:       "Space Dreams",
			ReleaseYear: 2020,
			Genres: []models.Genre{{
				ID:   2,
				Name: "Sci-Fi",
			}},
			Keywords: []string{
				"space",
				"exploration",
			},
			Casting: []models.Actor{{
				ID:        2,
				FirstName: "John",
				LastName:  "Star",
			}},
		},
		{
			ID:          3,
			TMDB_ID:     103,
			Title:       "The Lost City",
			ReleaseYear: 2019,
			Genres: []models.Genre{{
				ID:   3,
				Name: "Adventure",
			}},
			Keywords: []string{
				"jungle",
				"treasure",
			},
			Casting: []models.Actor{{
				ID:        3,
				FirstName: "Lara",
				LastName:  "Hunt",
			}},
		},
	}

	h.writeJSON(w, movies)
}

func (h *MovieHandler) GetRandomMovie(w http.ResponseWriter, r *http.Request) {
	movie := models.Movie{
		ID:          2,
		TMDB_ID:     102,
		Title:       "Space Dreams",
		ReleaseYear: 2020,
		Genres: []models.Genre{{
			ID:   2,
			Name: "Sci-Fi",
		}},
		Keywords: []string{
			"space",
			"exploration",
		},
		Casting: []models.Actor{{
			ID:        2,
			FirstName: "John",
			LastName:  "Star",
		}},
	}

	h.writeJSON(w, movie)
}
