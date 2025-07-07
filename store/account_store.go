package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/junaidshaikh-js/CineHubServer/models"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

type PostgresAccountStore struct {
	db *sql.DB
}

func NewPostgresAccountStore(db *sql.DB) *PostgresAccountStore {
	return &PostgresAccountStore{
		db: db,
	}
}

func (s *PostgresAccountStore) Register(name, email, password string) (bool, error) {
	// Validate basic requirements
	if name == "" || email == "" || password == "" {
		return false, ErrRegistrationValidation
	}

	// Check if user already exists
	var exists bool
	err := s.db.QueryRow(`
		SELECT EXISTS(SELECT 1 FROM users WHERE email = $1)
	`, email).Scan(&exists)

	if err != nil {
		return false, err
	}

	if exists {
		return false, ErrUserAlreadyExists
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	if err != nil {
		return false, err
	}

	// Insert new user
	query := `
		INSERT INTO users (name, email, password_hashed, time_created)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`

	var userID int
	err = s.db.QueryRow(
		query,
		name,
		email,
		string(hashedPassword),
		time.Now(),
	).Scan(&userID)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (s *PostgresAccountStore) Authenticate(email string, password string) (bool, error) {
	if email == "" || password == "" {
		return false, ErrAuthenticationValidation
	}

	// Fetch user by email
	var user models.User
	query := `
		SELECT id, name, email, password_hashed
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
	)
	if err == sql.ErrNoRows {
		return false, ErrAuthenticationValidation
	}
	if err != nil {
		return false, err
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return false, ErrAuthenticationValidation
	}

	// Update last login time
	updateQuery := `
		UPDATE users 
		SET last_login = $1
		WHERE id = $2
	`
	_, err = s.db.Exec(updateQuery, time.Now(), user.ID)
	if err != nil {
		// Do nothing
	}

	return true, nil
}

func (s *PostgresAccountStore) GetAccountDetails(email string) (models.User, error) {
	var user models.User
	query := `
		SELECT id, name, email
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`
	err := s.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
	)
	if err == sql.ErrNoRows {
		return models.User{}, ErrUserNotFound
	}
	if err != nil {
		return models.User{}, err
	}

	// Fetch favorites
	favoritesQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'favorite'
	`
	favoriteRows, err := s.db.Query(favoritesQuery, user.ID)
	if err != nil {
		return user, err
	}

	defer favoriteRows.Close()

	for favoriteRows.Next() {
		var m models.Movie
		if err := favoriteRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			return user, err
		}
		user.Favorites = append(user.Favorites, m)
	}

	// Fetch watchlist
	watchlistQuery := `
		SELECT m.id, m.tmdb_id, m.title, m.tagline, m.release_year, 
		       m.overview, m.score, m.popularity, m.language, 
		       m.poster_url, m.trailer_url
		FROM movies m
		JOIN user_movies um ON m.id = um.movie_id
		WHERE um.user_id = $1 AND um.relation_type = 'watchlist'
	`
	watchlistRows, err := s.db.Query(watchlistQuery, user.ID)
	if err != nil {
		return user, err
	}
	defer watchlistRows.Close()

	for watchlistRows.Next() {
		var m models.Movie
		if err := watchlistRows.Scan(
			&m.ID, &m.TMDB_ID, &m.Title, &m.Tagline, &m.ReleaseYear,
			&m.Overview, &m.Score, &m.Popularity, &m.Language,
			&m.PosterURL, &m.TrailerURL,
		); err != nil {
			return user, err
		}
		user.Watchlist = append(user.Watchlist, m)
	}

	return user, nil
}

func (s *PostgresAccountStore) SaveCollection(user models.User, movieID int, collection string) (bool, error) {
	// Validate inputs
	if movieID <= 0 {
		return false, errors.New("invalid movie ID")
	}
	if collection != "favorite" && collection != "watchlist" {
		return false, errors.New("collection must be 'favorite' or 'watchlist'")
	}

	// Get user ID from email
	var userID int
	err := s.db.QueryRow(`
		SELECT id 
		FROM users 
		WHERE email = $1 AND time_deleted IS NULL
	`, user.Email).Scan(&userID)
	if err == sql.ErrNoRows {
		return false, ErrUserNotFound
	}
	if err != nil {
		return false, err
	}

	// Check if the relationship already exists
	var exists bool
	err = s.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 
			FROM user_movies 
			WHERE user_id = $1 
			AND movie_id = $2 
			AND relation_type = $3
		)
	`, userID, movieID, collection).Scan(&exists)
	if err != nil {
		return false, err
	}
	if exists {
		return true, nil // Return true since the movie is already in the collection
	}

	// Insert the new relationship
	query := `
		INSERT INTO user_movies (user_id, movie_id, relation_type, time_added)
		VALUES ($1, $2, $3, $4)
	`
	_, err = s.db.Exec(query, userID, movieID, collection, time.Now())
	if err != nil {
		return false, err
	}

	return true, nil
}

var (
	ErrRegistrationValidation   = errors.New("registration failed")
	ErrAuthenticationValidation = errors.New("authentication failed")
	ErrUserAlreadyExists        = errors.New("user already exists")
	ErrUserNotFound             = errors.New("user not found")
)
