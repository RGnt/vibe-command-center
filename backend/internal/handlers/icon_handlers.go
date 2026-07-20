package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/internal/middleware"
	"todo-backend/internal/models"
	"todo-backend/internal/service"
	"github.com/go-chi/chi/v5"
	"log"
)

// IconHandler ...
type IconHandler struct {
	iconService service.IconService
}

// NewIconHandler ...
func NewIconHandler(iconService service.IconService) *IconHandler {
	return &IconHandler{
		iconService: iconService,
	}
}

// GetIcons ...
func (h *IconHandler) GetIcons(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	icons, err := h.iconService.GetIcons(userID)
	if err != nil {
		log.Println("Error fetching icons:", err)
		http.Error(w, "Error fetching icons", http.StatusInternalServerError)
		return
	}

	if icons == nil {
		icons = []models.Icon{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(icons)
}

// UpdateIcon ...
func (h *IconHandler) UpdateIcon(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.Icon
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedIcon, err := h.iconService.UpdateIcon(id, userID, req)
	if err != nil {
		http.Error(w, "Failed to update icon", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedIcon)
}

// DeleteIcon ...
func (h *IconHandler) DeleteIcon(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err := h.iconService.DeleteIcon(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete icon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
