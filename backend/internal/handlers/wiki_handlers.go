package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"todo-backend/internal/middleware"
	"todo-backend/internal/models"
	"todo-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// WikiHandler ...
type WikiHandler struct {
	wikiService service.WikiService
}

// NewWikiHandler ...
func NewWikiHandler(wikiService service.WikiService) *WikiHandler {
	return &WikiHandler{
		wikiService: wikiService,
	}
}

// GetWikis ...
func (h *WikiHandler) GetWikis(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	wikis, err := h.wikiService.GetWikis(userID)
	if err != nil {
		http.Error(w, "Failed to fetch wikis", http.StatusInternalServerError)
		return
	}

	if wikis == nil {
		wikis = []models.WikiPage{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(wikis)
}

// GetWiki ...
func (h *WikiHandler) GetWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := chi.URLParam(r, "slug")

	wiki, err := h.wikiService.GetWiki(slug, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Wiki not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch wiki", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(wiki)
}

// CreateWiki ...
func (h *WikiHandler) CreateWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var wiki models.WikiPage
	if err := json.NewDecoder(r.Body).Decode(&wiki); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdWiki, err := h.wikiService.CreateWiki(userID, wiki)
	if err != nil {
		http.Error(w, "Failed to create wiki", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdWiki)
}

// UpdateWiki ...
func (h *WikiHandler) UpdateWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid wiki ID", http.StatusBadRequest)
		return
	}

	var wiki models.WikiPage
	if err := json.NewDecoder(r.Body).Decode(&wiki); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedWiki, err := h.wikiService.UpdateWiki(id, userID, wiki)
	if err != nil {
		http.Error(w, "Failed to update wiki", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(updatedWiki)
}

// DeleteWiki ...
func (h *WikiHandler) DeleteWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid wiki ID", http.StatusBadRequest)
		return
	}

	err = h.wikiService.DeleteWiki(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete wiki", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
