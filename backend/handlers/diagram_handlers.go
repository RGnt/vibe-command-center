package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/middleware"
	"todo-backend/models"
	"todo-backend/service"
	"github.com/go-chi/chi/v5"
	"log"
)

type DiagramHandler struct {
	diagramService service.DiagramService
}

func NewDiagramHandler(diagramService service.DiagramService) *DiagramHandler {
	return &DiagramHandler{
		diagramService: diagramService,
	}
}

func (h *DiagramHandler) GetDiagrams(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	diagrams, err := h.diagramService.GetDiagrams(userID)
	if err != nil {
		log.Println("Error fetching diagrams:", err)
		http.Error(w, "Error fetching diagrams", http.StatusInternalServerError)
		return
	}

	if diagrams == nil {
		diagrams = []models.Diagram{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diagrams)
}

func (h *DiagramHandler) GetDiagram(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	d, err := h.diagramService.GetDiagram(id, userID)
	if err != nil {
		http.Error(w, "Diagram not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func (h *DiagramHandler) CreateDiagram(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.Diagram
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	createdDiagram, err := h.diagramService.CreateDiagram(userID, req)
	if err != nil {
		http.Error(w, "Failed to create diagram", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDiagram)
}

func (h *DiagramHandler) UpdateDiagram(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")
	var req models.Diagram
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedDiagram, err := h.diagramService.UpdateDiagram(id, userID, req)
	if err != nil {
		http.Error(w, "Failed to update diagram", http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedDiagram)
}

func (h *DiagramHandler) DeleteDiagram(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	err := h.diagramService.DeleteDiagram(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete diagram", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
