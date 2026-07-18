package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"todo-backend/database"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

func clearWorkflowsTable() {
	database.DB.Exec("DELETE FROM workflow_stages")
	database.DB.Exec("DELETE FROM workflows")
}

func setupWorkflowRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/workflows", GetWorkflows)
	r.Post("/api/workflows", CreateWorkflow)
	r.Put("/api/workflows/{id}", UpdateWorkflow)
	r.Delete("/api/workflows/{id}", DeleteWorkflow)
	return r
}

func TestCreateWorkflow(t *testing.T) {
	clearWorkflowsTable()
	router := setupWorkflowRouter()

	workflow := models.Workflow{
		Name: "Test Workflow",
		Stages: []models.WorkflowStage{
			{Name: "Stage 1"},
			{Name: "Stage 2"},
		},
	}
	body, _ := json.Marshal(workflow)
	req, _ := http.NewRequest("POST", "/api/workflows", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var responseWorkflow models.Workflow
	json.NewDecoder(rr.Body).Decode(&responseWorkflow)
	if responseWorkflow.Name != "Test Workflow" {
		t.Errorf("expected name to be 'Test Workflow', got %v", responseWorkflow.Name)
	}
	if len(responseWorkflow.Stages) != 2 {
		t.Errorf("expected 2 stages, got %v", len(responseWorkflow.Stages))
	}
}

func TestGetWorkflows(t *testing.T) {
	clearWorkflowsTable()
	router := setupWorkflowRouter()

	res, _ := database.DB.Exec(`INSERT INTO workflows (name) VALUES ('Test Workflow')`)
	id, _ := res.LastInsertId()
	database.DB.Exec(`INSERT INTO workflow_stages (workflow_id, name, "order") VALUES (?, 'Stage 1', 1)`, id)

	req, _ := http.NewRequest("GET", "/api/workflows", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var workflows []models.Workflow
	json.NewDecoder(rr.Body).Decode(&workflows)
	if len(workflows) != 1 {
		t.Errorf("expected 1 workflow, got %v", len(workflows))
	}
}

func TestUpdateWorkflow(t *testing.T) {
	clearWorkflowsTable()
	router := setupWorkflowRouter()

	res, _ := database.DB.Exec(`INSERT INTO workflows (name) VALUES ('Old Workflow')`)
	id, _ := res.LastInsertId()

	workflow := models.Workflow{
		Name: "Updated Workflow",
		Stages: []models.WorkflowStage{
			{Name: "New Stage 1"},
		},
	}
	body, _ := json.Marshal(workflow)
	req, _ := http.NewRequest("PUT", "/api/workflows/"+strconv.Itoa(int(id)), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseWorkflow models.Workflow
	json.NewDecoder(rr.Body).Decode(&responseWorkflow)
	if responseWorkflow.Name != "Updated Workflow" {
		t.Errorf("expected name to be 'Updated Workflow', got %v", responseWorkflow.Name)
	}
}

func TestDeleteWorkflow(t *testing.T) {
	clearWorkflowsTable()
	router := setupWorkflowRouter()

	res, _ := database.DB.Exec(`INSERT INTO workflows (name) VALUES ('Test Workflow')`)
	id, _ := res.LastInsertId()

	req, _ := http.NewRequest("DELETE", "/api/workflows/"+strconv.Itoa(int(id)), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM workflows WHERE id = ?", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 workflows, got %v", count)
	}
}
