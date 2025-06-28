package store

import "github.com/junaidshaikh-js/CineHubServer/models"

type MovieStore interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (models.Movie, error)
	GetMoviesByName(name string) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}
