package store

import "github.com/junaidshaikh-js/CineHubServer/models"

type MovieStore interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (*models.Movie, error)
	SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error)
}
