package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"todo-backend/middleware"
	"todo-backend/models"
	"todo-backend/service"

	"github.com/go-chi/chi/v5"
)

type WorkflowHandler struct {
	workflowService service.WorkflowService
}

func NewWorkflowHandler(workflowService service.WorkflowService) *WorkflowHandler {
	return &WorkflowHandler{
		workflowService: workflowService,
	}
}

func (h *WorkflowHandler) GetWorkflows(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	workflows, err := h.workflowService.GetWorkflows(userID)
	if err != nil {
		http.Error(w, "Failed to fetch workflows", http.StatusInternalServerError)
		return
	}

	if workflows == nil {
		workflows = []models.Workflow{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

func (h *WorkflowHandler) CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdWorkflow, err := h.workflowService.CreateWorkflow(userID, workflow)
	if err != nil {
		http.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkflow)
}

func (h *WorkflowHandler) UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedWorkflow, err := h.workflowService.UpdateWorkflow(id, userID, workflow)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Workflow not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedWorkflow)
}

func (h *WorkflowHandler) DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	err = h.workflowService.DeleteWorkflow(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete workflow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
