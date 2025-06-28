package store

import (
	"database/sql"

	"github.com/junaidshaikh-js/CineHubServer/models"
)

type PostgresMovieStore struct {
	DB *sql.DB
}

func NewPostgresMovieStore(db *sql.DB) *PostgresMovieStore {
	return &PostgresMovieStore{
		DB: db,
	}
}

const defaultLimit = 20

func (s *PostgresMovieStore) GetTopMovies() ([]models.Movie, error) {
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY popularity DESC
		LIMIT $1
	`

	return s.getMovies(query)
}

func (s *PostgresMovieStore) GetRandomMovies() ([]models.Movie, error) {
	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, popularity, language, poster_url, trailer_url
		FROM movies
		ORDER BY random()
		LIMIT $1
	`

	return s.getMovies(query)
}

func (s *PostgresMovieStore) getMovies(query string) ([]models.Movie, error) {
	rows, err := s.DB.Query(query, defaultLimit)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var movies []models.Movie

	for rows.Next() {
		var movie models.Movie
		err := rows.Scan(&movie.ID, &movie.TMDB_ID, &movie.Title, &movie.Tagline, &movie.ReleaseYear, &movie.Overview, &movie.Score, &movie.Popularity, &movie.Language, &movie.PosterURL, &movie.TrailerURL)

		if err != nil {
			return nil, err
		}

		movies = append(movies, movie)
	}

	return movies, nil
}
