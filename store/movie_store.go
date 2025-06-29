package store

import (
	"database/sql"
	"strconv"

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

func (s *PostgresMovieStore) GetMovieByID(id int) (*models.Movie, error) {
	movie := &models.Movie{}

	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, popularity, language, poster_url, trailer_url
		FROM movies
		WHERE id = $1
	`

	err := s.DB.QueryRow(query, id).Scan(&movie.ID, &movie.TMDB_ID, &movie.Title, &movie.Tagline, &movie.ReleaseYear, &movie.Overview, &movie.Score, &movie.Popularity, &movie.Language, &movie.PosterURL, &movie.TrailerURL)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *PostgresMovieStore) SearchMoviesByName(name string, order string, genre *int) ([]models.Movie, error) {
	orderBy := "popularity DESC"
	switch order {
	case "score":
		orderBy = "score DESC"
	case "name":
		orderBy = "title"
	case "date":
		orderBy = "release_year DESC"
	}

	genreFilter := ""
	if genre != nil {
		genreFilter = `AND 
			(SELECT COUNT(*) FROM movie_genres WHERE movie_id=movies.id AND genre_id= ` + strconv.Itoa(*genre) + `) = 1`
	}

	query := `
		SELECT id, tmdb_id, title, tagline, release_year, overview, score, popularity, language, poster_url, trailer_url
		FROM movies
		WHERE (title ILIKE $1 OR overview ILIKE $1)` + genreFilter + `
		ORDER BY ` + orderBy + `
		LIMIT $2
	`

	rows, err := s.DB.Query(query, "%"+name+"%", defaultLimit)

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

func (s *PostgresMovieStore) GetAllGenres() ([]models.Genre, error) {
	query := `
		SELECT id, name
		FROM genres
		ORDER BY id
	`

	rows, err := s.DB.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var genres []models.Genre

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.Name)

		if err != nil {
			return nil, err
		}

		genres = append(genres, genre)
	}

	return genres, nil
}
