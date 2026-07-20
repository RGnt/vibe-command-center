package handlers

import (
	"encoding/json"
	"net/http"

	"todo-backend/middleware"
	"todo-backend/models"
	"todo-backend/service"
)

// UserHandler ...
type UserHandler struct {
	userService service.UserService
}

// NewUserHandler ...
func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// GetUserSettings ...
func (h *UserHandler) GetUserSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	settings, err := h.userService.GetSettings(userID)
	if err != nil {
		http.Error(w, "Error fetching settings", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(settings)
}

// UpdateUserSettings ...
func (h *UserHandler) UpdateUserSettings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var settings models.UserSettings
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedSettings, err := h.userService.UpdateSettings(userID, settings)
	if err != nil {
		http.Error(w, "Error updating settings", http.StatusInternalServerError)
		return
	}

	_ = json.NewEncoder(w).Encode(updatedSettings)
}
