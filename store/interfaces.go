package store

import "github.com/junaidshaikh-js/CineHubServer/models"

type MovieStore interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
}
