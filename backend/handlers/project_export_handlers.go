package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"todo-backend/middleware"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

func (h *ProjectHandler) ExportProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectIDStr := chi.URLParam(r, "id")
	projectID, err := strconv.Atoi(projectIDStr)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	payload, err := h.projectService.ExportProject(projectID, userID)
	if err != nil {
		http.Error(w, "Failed to export project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", `attachment; filename="project_export.json"`)
	json.NewEncoder(w).Encode(payload)
}

func (h *ProjectHandler) ImportProject(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var payload models.ProjectExportPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	newProject, err := h.projectService.ImportProject(userID, payload)
	if err != nil {
		http.Error(w, "Failed to import project: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newProject)
}
