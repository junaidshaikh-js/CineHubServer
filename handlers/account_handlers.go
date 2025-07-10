package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/junaidshaikh-js/CineHubServer/logger"
	"github.com/junaidshaikh-js/CineHubServer/models"
	"github.com/junaidshaikh-js/CineHubServer/store"
	"github.com/junaidshaikh-js/CineHubServer/token"
	"github.com/junaidshaikh-js/CineHubServer/utils"
)

type RegisterRequestPayload struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthRequestPayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Token   string `json:"token"`
}

type AccountHandler struct {
	accountStore store.AccountStore
	logger       *logger.Logger
}

func NewAccountHandler(accountStore store.AccountStore, logger *logger.Logger) *AccountHandler {
	return &AccountHandler{
		accountStore: accountStore,
		logger:       logger,
	}
}

func (h *AccountHandler) handleStorageError(w http.ResponseWriter, err error, context string) bool {
	if err != nil {
		switch err {
		case store.ErrAuthenticationValidation, store.ErrUserAlreadyExists, store.ErrRegistrationValidation:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(AuthResponse{Success: false, Message: err.Error()})
			return true
		case store.ErrUserNotFound:
			http.Error(w, "User not found", http.StatusNotFound)
			return true
		default:
			h.logger.Error(context, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return true
		}
	}
	return false
}

func (h *AccountHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequestPayload
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Error("Failed to decode registration request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success, err := h.accountStore.Register(req.Name, req.Email, req.Password)

	if h.handleStorageError(w, err, "Failed to register user") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "User registered successfully",
		Token:   token.CreateJWT(models.User{Email: req.Email, Name: req.Name}, *h.logger),
	}

	if err := utils.WriteJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully registered user with email: " + req.Email)
	}
}

func (h *AccountHandler) Authenticate(w http.ResponseWriter, r *http.Request) {
	var req AuthRequestPayload
	err := json.NewDecoder(r.Body).Decode(&req)

	if err != nil {
		h.logger.Error("Failed to decode authentication request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	success, err := h.accountStore.Authenticate(req.Email, req.Password)
	if h.handleStorageError(w, err, "Failed to authenticate user") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "User authenticated successfully",
		Token:   token.CreateJWT(models.User{Email: req.Email}, *h.logger),
	}

	if err := utils.WriteJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully authenticated user with email: " + req.Email)
	}
}

func (h *AccountHandler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		if tokenStr == "" {
			http.Error(w, "Missing authorization token", http.StatusUnauthorized)
			return
		}

		// Remove "Bearer " prefix if present
		tokenStr = strings.TrimPrefix(tokenStr, "Bearer ")

		// Parse and validate the token

		token, err := jwt.Parse(tokenStr,
			func(t *jwt.Token) (interface{}, error) {
				// Ensure the signing method is HMAC
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, jwt.ErrSignatureInvalid
				}
				return []byte(token.GetJWTSecret(*h.logger)), nil
			},
		)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract claims from the token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Get the email from claims
		email, ok := claims["email"].(string)
		if !ok {
			http.Error(w, "Email not found in token", http.StatusUnauthorized)
			return
		}

		// Inject email into the request context
		ctx := context.WithValue(r.Context(), "email", email)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *AccountHandler) SaveToCollection(w http.ResponseWriter, r *http.Request) {
	type CollectionRequest struct {
		MovieID    int    `json:"movieId"`
		Collection string `json:"collection"`
	}

	var req CollectionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode collection request", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}

	success, err := h.accountStore.SaveCollection(models.User{Email: email},
		req.MovieID, req.Collection)
	if h.handleStorageError(w, err, "Failed to save to collection") {
		return
	}

	response := AuthResponse{
		Success: success,
		Message: "Movie added to " + req.Collection + " successfully",
	}

	if err := utils.WriteJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully saved movie to " + req.Collection)
	}
}

func (h *AccountHandler) GetFavorites(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	details, err := h.accountStore.GetAccountDetails(email)
	if err != nil {
		http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
		return
	}
	if err := utils.WriteJSONResponse(w, details.Favorites); err == nil {
		h.logger.Info("Successfully sent favorites")
	}
}

func (h *AccountHandler) GetWatchlist(w http.ResponseWriter, r *http.Request) {
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Unable to retrieve email", http.StatusInternalServerError)
		return
	}
	details, err := h.accountStore.GetAccountDetails(email)
	if err != nil {
		http.Error(w, "Unable to retrieve collections", http.StatusInternalServerError)
		return
	}
	if err := utils.WriteJSONResponse(w, details.Watchlist); err == nil {
		h.logger.Info("Successfully sent favorites")
	}
}
