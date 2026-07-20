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

func setupDiagramRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/diagrams", testutils.AuthContext(1, diagramHandler.GetDiagrams))
	r.Get("/api/diagrams/{id}", testutils.AuthContext(1, diagramHandler.GetDiagram))
	r.Post("/api/diagrams", testutils.AuthContext(1, diagramHandler.CreateDiagram))
	r.Put("/api/diagrams/{id}", testutils.AuthContext(1, diagramHandler.UpdateDiagram))
	r.Delete("/api/diagrams/{id}", testutils.AuthContext(1, diagramHandler.DeleteDiagram))
	return r
}

// TestCreateDiagram ...
func TestCreateDiagram(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupDiagramRouter()

	diagram := models.Diagram{
		Name:        "Test Diagram",
		DiagramType: "graph TD",
		Code:        "A-->B",
	}
	body, _ := json.Marshal(diagram)
	req, _ := http.NewRequest("POST", "/api/diagrams", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var responseDiagram models.Diagram
	_ = json.NewDecoder(rr.Body).Decode(&responseDiagram)
	if responseDiagram.Name != "Test Diagram" {
		t.Errorf("expected name to be 'Test Diagram', got %v", responseDiagram.Name)
	}
}

// TestGetDiagrams ...
func TestGetDiagrams(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupDiagramRouter()

	_, _ = database.DB.Exec(`INSERT INTO diagrams (user_id, name, diagram_type, code) VALUES (1, 'Test Diagram', 'graph TD', 'A-->B')`)

	req, _ := http.NewRequest("GET", "/api/diagrams", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var diagrams []models.Diagram
	_ = json.NewDecoder(rr.Body).Decode(&diagrams)
	if len(diagrams) != 1 {
		t.Errorf("expected 1 diagram, got %v", len(diagrams))
	}
}

// TestGetDiagram ...
func TestGetDiagram(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupDiagramRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO diagrams (user_id, name, diagram_type, code) VALUES (1, 'Test Diagram', 'graph TD', 'A-->B') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("GET", "/api/diagrams/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var d models.Diagram
	_ = json.NewDecoder(rr.Body).Decode(&d)
	if d.Name != "Test Diagram" {
		t.Errorf("expected name Test Diagram, got %v", d.Name)
	}
}

// TestUpdateDiagram ...
func TestUpdateDiagram(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupDiagramRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO diagrams (user_id, name, diagram_type, code) VALUES (1, 'Test Diagram', 'graph TD', 'A-->B') RETURNING id`).Scan(&id)

	diagram := models.Diagram{
		Name:        "Updated Diagram",
		DiagramType: "graph TD",
		Code:        "A-->B-->C",
	}
	body, _ := json.Marshal(diagram)
	req, _ := http.NewRequest("PUT", "/api/diagrams/"+strconv.Itoa(id), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseDiagram models.Diagram
	_ = json.NewDecoder(rr.Body).Decode(&responseDiagram)
	if responseDiagram.Name != "Updated Diagram" {
		t.Errorf("expected name to be 'Updated Diagram', got %v", responseDiagram.Name)
	}
}

// TestDeleteDiagram ...
func TestDeleteDiagram(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupDiagramRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO diagrams (user_id, name, diagram_type, code) VALUES (1, 'Test Diagram', 'graph TD', 'A-->B') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("DELETE", "/api/diagrams/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	var count int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM diagrams WHERE id = $1", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 diagrams, got %v", count)
	}
}
