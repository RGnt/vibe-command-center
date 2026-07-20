package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/mail"
	"os"
	"strings"

	"todo-backend/internal/middleware"
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
	"todo-backend/internal/service"
)

// AuthInput ...
type AuthInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// AuthResponse ...
type AuthResponse struct {
	User  models.User `json:"user"` // AUTH-10: Token removed
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

// Register handles user registration.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var input AuthInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Validate email format
	input.Email = strings.TrimSpace(strings.ToLower(input.Email))
	if _, err := mail.ParseAddress(input.Email); err != nil {
		http.Error(w, "Invalid email address", http.StatusBadRequest)
		return
	}

	// Validate password length
	if len(input.Password) < 8 {
		http.Error(w, "Password must be at least 8 characters", http.StatusBadRequest)
		return
	}
	if len(input.Password) > 128 {
		http.Error(w, "Password must not exceed 128 characters", http.StatusBadRequest)
		return
	}

	user, accessToken, refreshToken, err := h.authService.Register(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, repository.ErrEmailExists) {
			// AUTH-05: Prevent account enumeration
			w.WriteHeader(http.StatusAccepted)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "If this email is not already registered, your account has been created.",
			})
			return
		}
		log.Printf("Registration error: %v\n", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, accessToken, refreshToken)

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

	user, accessToken, refreshToken, err := h.authService.Login(input.Email, input.Password)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, accessToken, refreshToken)
	_ = json.NewEncoder(w).Encode(AuthResponse{User: user})
}

// Refresh issues a new access token and refresh token.
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("refresh_token")
	if err != nil || cookie.Value == "" {
		http.Error(w, "Missing refresh token", http.StatusUnauthorized)
		return
	}

	accessToken, newRefreshToken, err := h.authService.Refresh(cookie.Value)
	if err != nil {
		if errors.Is(err, repository.ErrInvalidRefreshToken) {
			http.Error(w, "Invalid or expired refresh token", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	h.setAuthCookies(w, accessToken, newRefreshToken)
	w.WriteHeader(http.StatusOK)
}

// Logout clears the authentication cookies and revokes the current session.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	// AUTH-01: Revoke the token if it exists
	if cookie, err := r.Cookie("token"); err == nil && cookie.Value != "" {
		_ = h.authService.RevokeTokenByString(cookie.Value)
	}

	if rtCookie, err := r.Cookie("refresh_token"); err == nil && rtCookie.Value != "" {
		_ = h.authService.Logout(rtCookie.Value)
	}

	secureCookie := os.Getenv("IS_HTTPS") == "true"
	
	// Clear access token
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie, // AUTH-03
		SameSite: http.SameSiteStrictMode, // AUTH-08
	})

	// Clear refresh token
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	})

	w.WriteHeader(http.StatusOK)
}

// GetMe returns the currently authenticated user.
func (h *AuthHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	user, err := h.authService.GetUserByID(userID)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			log.Printf("GetMe DB error for userID %d: %v", userID, err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	_ = json.NewEncoder(w).Encode(user)
}

func (h *AuthHandler) setAuthCookies(w http.ResponseWriter, accessToken, refreshToken string) {
	secureCookie := os.Getenv("IS_HTTPS") == "true"
	
	// Access token (15m)
	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    accessToken,
		Path:     "/",
		MaxAge:   900, // 15 minutes
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	})

	// Refresh token (30 days)
	http.SetCookie(w, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		Path:     "/",
		MaxAge:   2592000, // 30 days
		HttpOnly: true,
		Secure:   secureCookie,
		SameSite: http.SameSiteStrictMode,
	})
}
