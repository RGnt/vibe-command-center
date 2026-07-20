package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"todo-backend/internal/database"
	"todo-backend/internal/models"
	"todo-backend/internal/testutils"

	"github.com/go-chi/chi/v5"
)

func setupProjectRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/projects", testutils.AuthContext(1, projectHandler.GetProjects))
	r.Get("/api/projects/{id}", testutils.AuthContext(1, projectHandler.GetProject))
	r.Post("/api/projects", testutils.AuthContext(1, projectHandler.CreateProject))
	r.Put("/api/projects/{id}", testutils.AuthContext(1, projectHandler.UpdateProject))
	r.Delete("/api/projects/{id}", testutils.AuthContext(1, projectHandler.DeleteProject))
	return r
}

// TestCreateProject ...
func TestCreateProject(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupProjectRouter()

	project := models.Project{
		Name:        "Test Project",
		Description: "Project Description",
	}
	body, _ := json.Marshal(project)
	req, _ := http.NewRequest("POST", "/api/projects", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var responseProject models.Project
	_ = json.NewDecoder(rr.Body).Decode(&responseProject)
	if responseProject.Name != "Test Project" {
		t.Errorf("expected name to be 'Test Project', got %v", responseProject.Name)
	}
}

// TestGetProjects ...
func TestGetProjects(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupProjectRouter()

	_, _ = database.DB.Exec(`INSERT INTO projects (user_id, name, description) VALUES (1, 'Test Project 1', 'Desc')`)

	req, _ := http.NewRequest("GET", "/api/projects", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var projects []models.Project
	_ = json.NewDecoder(rr.Body).Decode(&projects)
	if len(projects) != 1 {
		t.Errorf("expected 1 project, got %v", len(projects))
	}
}

// TestGetProject ...
func TestGetProject(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupProjectRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO projects (user_id, name, description) VALUES (1, 'Test Project 1', 'Desc') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("GET", "/api/projects/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var p models.Project
	_ = json.NewDecoder(rr.Body).Decode(&p)
	if p.Name != "Test Project 1" {
		t.Errorf("expected name Test Project 1, got %v", p.Name)
	}
}

// TestUpdateProject ...
func TestUpdateProject(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupProjectRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO projects (user_id, name, description) VALUES (1, 'Test Project 1', 'Desc') RETURNING id`).Scan(&id)

	project := models.Project{
		Name:        "Updated Project",
		Description: "Updated Desc",
	}
	body, _ := json.Marshal(project)
	req, _ := http.NewRequest("PUT", "/api/projects/"+strconv.Itoa(id), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseProject models.Project
	_ = json.NewDecoder(rr.Body).Decode(&responseProject)
	if responseProject.Name != "Updated Project" {
		t.Errorf("expected name to be 'Updated Project', got %v", responseProject.Name)
	}
}

// TestDeleteProject ...
func TestDeleteProject(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupProjectRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO projects (user_id, name, description) VALUES (1, 'Test Project 1', 'Desc') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("DELETE", "/api/projects/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	var count int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM projects WHERE id = $1", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 projects, got %v", count)
	}
}
