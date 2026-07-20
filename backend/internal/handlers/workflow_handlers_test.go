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

func setupWorkflowRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/workflows", testutils.AuthContext(1, workflowHandler.GetWorkflows))
	r.Post("/api/workflows", testutils.AuthContext(1, workflowHandler.CreateWorkflow))
	r.Put("/api/workflows/{id}", testutils.AuthContext(1, workflowHandler.UpdateWorkflow))
	r.Delete("/api/workflows/{id}", testutils.AuthContext(1, workflowHandler.DeleteWorkflow))
	return r
}

// TestCreateWorkflow ...
func TestCreateWorkflow(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWorkflowRouter()

	workflow := models.Workflow{
		Name: "Test Workflow",
		Stages: []models.WorkflowStage{
			{Name: "Stage 1", Order: 1},
			{Name: "Stage 2", Order: 2},
		},
	}
	body, _ := json.Marshal(workflow)
	req, _ := http.NewRequest("POST", "/api/workflows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var responseWorkflow models.Workflow
	_ = json.NewDecoder(rr.Body).Decode(&responseWorkflow)
	if responseWorkflow.Name != "Test Workflow" {
		t.Errorf("expected name to be 'Test Workflow', got %v", responseWorkflow.Name)
	}
	if len(responseWorkflow.Stages) != 2 {
		t.Errorf("expected 2 stages, got %v", len(responseWorkflow.Stages))
	}
}

// TestGetWorkflows ...
func TestGetWorkflows(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWorkflowRouter()

	_, _ = database.DB.Exec(`INSERT INTO workflows (user_id, name) VALUES (1, 'Test Workflow 1')`)

	req, _ := http.NewRequest("GET", "/api/workflows", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var workflows []models.Workflow
	_ = json.NewDecoder(rr.Body).Decode(&workflows)
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %v", len(workflows))
	}
}

// TestUpdateWorkflow ...
func TestUpdateWorkflow(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWorkflowRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO workflows (user_id, name) VALUES (1, 'Test Workflow 1') RETURNING id`).Scan(&id)

	workflow := models.Workflow{
		Name: "Updated Workflow",
		Stages: []models.WorkflowStage{
			{Name: "Stage 1", Order: 1},
		},
	}
	body, _ := json.Marshal(workflow)
	req, _ := http.NewRequest("PUT", "/api/workflows/"+strconv.Itoa(id), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseWorkflow models.Workflow
	_ = json.NewDecoder(rr.Body).Decode(&responseWorkflow)
	if responseWorkflow.Name != "Updated Workflow" {
		t.Errorf("expected name to be 'Updated Workflow', got %v", responseWorkflow.Name)
	}
}

// TestDeleteWorkflow ...
func TestDeleteWorkflow(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWorkflowRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO workflows (user_id, name) VALUES (1, 'Test Workflow 1') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("DELETE", "/api/workflows/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	var count int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM workflows WHERE id = $1", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 workflows, got %v", count)
	}
}
