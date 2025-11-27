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

	err = s.getMovieRelations(movie)

	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (s *PostgresMovieStore) getMovieRelations(movie *models.Movie) error {
	genreQuery := `
		SELECT g.id, g.name 
		FROM genres g
		JOIN movie_genres mg ON g.id = mg.genre_id
		WHERE mg.movie_id = $1
	`

	rows, err := s.DB.Query(genreQuery, movie.ID)

	if err != nil {
		return err
	}

	defer rows.Close()

	for rows.Next() {
		var genre models.Genre
		err := rows.Scan(&genre.ID, &genre.Name)

		if err != nil {
			return err
		}

		movie.Genres = append(movie.Genres, genre)
	}

	actorsQuery := `
		SELECT a.id, a.first_name, a.last_name, a.image_url
		FROM actors a
		JOIN movie_cast mc ON mc.actor_id = a.id
		WHERE mc.movie_id = $1
	`

	actorRows, err := s.DB.Query(actorsQuery, movie.ID)

	if err != nil {
		return err
	}

	defer actorRows.Close()

	for actorRows.Next() {
		var actor models.Actor
		err := actorRows.Scan(&actor.ID, &actor.FirstName, &actor.LastName, &actor.ImageURL)

		if err != nil {
			return err
		}

		movie.Casting = append(movie.Casting, actor)
	}

	keywordsQuery := `
		SELECT k.word
		FROM keywords k
		JOIN movie_keywords mk ON mk.keyword_id = k.id
		WHERE mk.movie_id = $1
	`

	keywordsRows, err := s.DB.Query(keywordsQuery, movie.ID)

	if err != nil {
		return err
	}

	defer keywordsRows.Close()

	for keywordsRows.Next() {
		var keyword string
		err := keywordsRows.Scan(&keyword)

		if err != nil {
			return err
		}

		movie.Keywords = append(movie.Keywords, keyword)
	}

	return nil
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

	query := `
    SELECT id, tmdb_id, title, tagline, release_year, overview, score, popularity, language, poster_url, trailer_url
    FROM movies
    WHERE (title ILIKE $1 OR overview ILIKE $1)
  `
	args := []any{"%" + name + "%", defaultLimit}

	if genre != nil {
		query += `
        AND (SELECT COUNT(*) FROM movie_genres 
             WHERE movie_id=movies.id AND genre_id = $3) = 1
    `
		args = append(args, *genre)
	}

	query += `
    ORDER BY ` + orderBy + `
    LIMIT $2
`

	rows, err := s.DB.Query(query, args...)

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
