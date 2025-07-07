package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/junaidshaikh-js/CineHubServer/logger"
	"github.com/junaidshaikh-js/CineHubServer/store"
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
	}

	if err := utils.WriteJSONResponse(w, response); err == nil {
		h.logger.Info("Successfully authenticated user with email: " + req.Email)
	}
}
