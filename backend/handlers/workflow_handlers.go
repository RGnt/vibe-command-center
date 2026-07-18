package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"todo-backend/database"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

// GetWorkflows gets all workflows
func GetWorkflows(w http.ResponseWriter, r *http.Request) {
	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	rows, err := database.DB.Query(`SELECT id, name, created_at FROM workflows ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "Failed to fetch workflows", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var workflow models.Workflow
		var createdAt string

		err := rows.Scan(&workflow.ID, &workflow.Name, &createdAt)
		if err != nil {
			http.Error(w, "Failed to scan workflow", http.StatusInternalServerError)
			return
		}

		t, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err == nil {
			workflow.CreatedAt = t
		}

		stageRows, err := database.DB.Query(`SELECT id, name, "order" FROM workflow_stages WHERE workflow_id = ? ORDER BY "order"`, workflow.ID)
		if err != nil {
			http.Error(w, "Failed to fetch workflow stages", http.StatusInternalServerError)
			return
		}
		defer stageRows.Close()

		var stages []models.WorkflowStage
		for stageRows.Next() {
			var stage models.WorkflowStage
			err := stageRows.Scan(&stage.ID, &stage.Name, &stage.Order)
			if err != nil {
				http.Error(w, "Failed to scan workflow stage", http.StatusInternalServerError)
				return
			}
			stages = append(stages, stage)
		}
		workflow.Stages = stages

		workflows = append(workflows, workflow)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflows)
}

// CreateWorkflow creates a new workflow
func CreateWorkflow(w http.ResponseWriter, r *http.Request) {
	var workflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&workflow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `INSERT INTO workflows (name) VALUES (?)`
	result, err := database.DB.Exec(query, workflow.Name)
	if err != nil {
		http.Error(w, "Failed to create workflow", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get workflow ID", http.StatusInternalServerError)
		return
	}

	for i, stage := range workflow.Stages {
		stageQuery := `INSERT INTO workflow_stages (workflow_id, name, "order") VALUES (?, ?, ?)`
		_, err := database.DB.Exec(stageQuery, id, stage.Name, i+1)
		if err != nil {
			http.Error(w, "Failed to create workflow stage", http.StatusInternalServerError)
			return
		}
	}

	var createdWorkflow models.Workflow
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, name, created_at FROM workflows WHERE id = ?`, id).Scan(&createdWorkflow.ID, &createdWorkflow.Name, &createdAt)
	if err != nil {
		http.Error(w, "Failed to fetch created workflow", http.StatusInternalServerError)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		createdWorkflow.CreatedAt = t
	}

	stageRows, err := database.DB.Query(`SELECT id, name, "order" FROM workflow_stages WHERE workflow_id = ? ORDER BY "order"`, id)
	if err != nil {
		http.Error(w, "Failed to fetch workflow stages", http.StatusInternalServerError)
		return
	}
	defer stageRows.Close()

	var stages []models.WorkflowStage
	for stageRows.Next() {
		var stage models.WorkflowStage
		err := stageRows.Scan(&stage.ID, &stage.Name, &stage.Order)
		if err != nil {
			http.Error(w, "Failed to scan workflow stage", http.StatusInternalServerError)
			return
		}
		stages = append(stages, stage)
	}
	createdWorkflow.Stages = stages

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdWorkflow)
}

// UpdateWorkflow updates a workflow
func UpdateWorkflow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	var updatedWorkflow models.Workflow
	if err := json.NewDecoder(r.Body).Decode(&updatedWorkflow); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	_, err = database.DB.Exec("UPDATE workflows SET name = ? WHERE id = ?", updatedWorkflow.Name, id)
	if err != nil {
		http.Error(w, "Failed to update workflow", http.StatusInternalServerError)
		return
	}

	_, err = database.DB.Exec("DELETE FROM workflow_stages WHERE workflow_id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete workflow stages", http.StatusInternalServerError)
		return
	}

	for i, stage := range updatedWorkflow.Stages {
		stageQuery := `INSERT INTO workflow_stages (workflow_id, name, "order") VALUES (?, ?, ?)`
		_, err := database.DB.Exec(stageQuery, id, stage.Name, i+1)
		if err != nil {
			http.Error(w, "Failed to create workflow stage", http.StatusInternalServerError)
			return
		}
	}

	var workflow models.Workflow
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, name, created_at FROM workflows WHERE id = ?`, id).Scan(&workflow.ID, &workflow.Name, &createdAt)
	if err != nil {
		http.Error(w, "Failed to fetch updated workflow", http.StatusInternalServerError)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		workflow.CreatedAt = t
	}

	stageRows, err := database.DB.Query(`SELECT id, name, "order" FROM workflow_stages WHERE workflow_id = ? ORDER BY "order"`, id)
	if err != nil {
		http.Error(w, "Failed to fetch workflow stages", http.StatusInternalServerError)
		return
	}
	defer stageRows.Close()

	var stages []models.WorkflowStage
	for stageRows.Next() {
		var stage models.WorkflowStage
		err := stageRows.Scan(&stage.ID, &stage.Name, &stage.Order)
		if err != nil {
			http.Error(w, "Failed to scan workflow stage", http.StatusInternalServerError)
			return
		}
		stages = append(stages, stage)
	}
	workflow.Stages = stages

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workflow)
}

// DeleteWorkflow deletes a workflow
func DeleteWorkflow(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid workflow ID", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	_, err = database.DB.Exec("DELETE FROM workflows WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete workflow", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
