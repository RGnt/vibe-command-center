package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"todo-backend/database"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

// GetProjects gets all projects
func GetProjects(w http.ResponseWriter, r *http.Request) {
	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	rows, err := database.DB.Query(`SELECT id, name, description, workflow_id, created_at FROM projects ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "Failed to fetch projects", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		var workflowID sql.NullInt64
		var createdAt string

		err := rows.Scan(&project.ID, &project.Name, &project.Description, &workflowID, &createdAt)
		if err != nil {
			http.Error(w, "Failed to scan project", http.StatusInternalServerError)
			return
		}

		t, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err == nil {
			project.CreatedAt = t
		}

		if workflowID.Valid {
			wid := int(workflowID.Int64)
			project.WorkflowID = &wid
		}

		projects = append(projects, project)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projects)
}

// GetProject gets a specific project
func GetProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	var project models.Project
	var workflowID sql.NullInt64
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, name, description, workflow_id, created_at FROM projects WHERE id = ?`, id).Scan(&project.ID, &project.Name, &project.Description, &workflowID, &createdAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Project not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch project", http.StatusInternalServerError)
		}
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		project.CreatedAt = t
	}

	if workflowID.Valid {
		wid := int(workflowID.Int64)
		project.WorkflowID = &wid
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(project)
}

// CreateProject creates a new project
func CreateProject(w http.ResponseWriter, r *http.Request) {
	var project models.Project
	if err := json.NewDecoder(r.Body).Decode(&project); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `INSERT INTO projects (name, description, workflow_id) VALUES (?, ?, ?)`
	result, err := database.DB.Exec(query, project.Name, project.Description, project.WorkflowID)
	if err != nil {
		http.Error(w, "Failed to create project", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get project ID", http.StatusInternalServerError)
		return
	}

	project.ID = int(id)
	project.CreatedAt = time.Now() // Close enough for response

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(project)
}

// UpdateProject updates a project
func UpdateProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	var updatedProject models.Project
	if err := json.NewDecoder(r.Body).Decode(&updatedProject); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `UPDATE projects SET name = ?, description = ?, workflow_id = ? WHERE id = ?`
	_, err = database.DB.Exec(query, updatedProject.Name, updatedProject.Description, updatedProject.WorkflowID, id)
	if err != nil {
		http.Error(w, "Failed to update project", http.StatusInternalServerError)
		return
	}

	updatedProject.ID = id

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedProject)
}

// DeleteProject deletes a project
func DeleteProject(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	_, err = database.DB.Exec("DELETE FROM projects WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete project", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
