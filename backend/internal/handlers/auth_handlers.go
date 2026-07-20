package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"

	"todo-backend/internal/middleware"
	"todo-backend/internal/models"
	"todo-backend/internal/service"
)

// AuthInput ...
type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse ...
type AuthResponse struct {
	Token string      `json:"token"`
	User  models.User `json:"user"`
}

// AuthHandler ...
type AuthHandler struct {
	authService service.AuthService
}

// NewAuthHandler ...
func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

// Register ...
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.Register(input.Email, input.Password)
	if err != nil {
		log.Printf("Registration error: %v\n", err)
		if strings.Contains(err.Error(), "email might already exist") { // Could use strongly typed error here too, exported from repository layer
			http.Error(w, "Email might already exist", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   604800, // 1 week
		HttpOnly: true,
		Secure:   false, // Set to true in production if using HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(AuthResponse{User: user})
}

// Login ...
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	user, token, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		Path:     "/",
		MaxAge:   604800, // 1 week
		HttpOnly: true,
		Secure:   false, // Set to true in production if using HTTPS
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)

	_ = json.NewEncoder(w).Encode(AuthResponse{User: user})
}

// Logout ...
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusOK)
}

// GetMe ...
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	_ = json.NewEncoder(w).Encode(user)
}
