package store

import "github.com/junaidshaikh-js/CineHubServer/models"

type MovieStore interface {
	GetTopMovies() ([]models.Movie, error)
	GetRandomMovies() ([]models.Movie, error)
	GetMovieByID(id int) (*models.Movie, error)
	SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error)
	GetAllGenres() ([]models.Genre, error)
}

type AccountStore interface {
	Authenticate(email, password string) (bool, error)
	Register(name, email, password string) (bool, error)
	GetAccountDetails(email string) (models.User, error)
	SaveCollection(user models.User, movieId int, collection string) (bool, error)
}
