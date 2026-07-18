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
	"todo-backend/testutils"

	"github.com/go-chi/chi/v5"
)

func setupTodoRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/todos", testutils.AuthContext(1, todoHandler.GetTodos))
	r.Get("/api/todos/{id}", testutils.AuthContext(1, todoHandler.GetTodo))
	r.Post("/api/todos", testutils.AuthContext(1, todoHandler.CreateTodo))
	r.Post("/api/todos/{parent_id}/subtasks", testutils.AuthContext(1, todoHandler.CreateSubtask))
	r.Put("/api/todos/{id}", testutils.AuthContext(1, todoHandler.UpdateTodo))
	r.Delete("/api/todos/{id}", testutils.AuthContext(1, todoHandler.DeleteTodo))
	r.Patch("/api/todos/{id}/toggle", testutils.AuthContext(1, todoHandler.ToggleTodo))
	r.Get("/api/todos/stage", testutils.AuthContext(1, todoHandler.GetTodosByStage))
	return r
}

func setupTestUser() {
	database.DB.Exec("INSERT INTO users (id, email, password_hash) VALUES (1, 'test@example.com', 'hash') ON CONFLICT DO NOTHING")
}

func TestCreateTodo(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupTodoRouter()

	todo := models.Todo{
		Title:   "Test Todo",
		Content: "Test Content",
		Stage:   "To Do",
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var responseTodo models.Todo
	json.NewDecoder(rr.Body).Decode(&responseTodo)
	if responseTodo.Title != "Test Todo" {
		t.Errorf("expected title to be 'Test Todo', got %v", responseTodo.Title)
	}
	if responseTodo.Content != "Test Content" {
		t.Errorf("expected content to be 'Test Content', got %v", responseTodo.Content)
	}
}

func TestGetTodos(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupTodoRouter()

	database.DB.Exec(`INSERT INTO todos (user_id, title, content, stage, completed) VALUES (1, 'Test Todo 1', 'Test Content', 'To Do', false)`)

	req, _ := http.NewRequest("GET", "/api/todos", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var todos []models.Todo
	json.NewDecoder(rr.Body).Decode(&todos)
	if len(todos) != 1 {
		t.Errorf("expected 1 todo, got %v", len(todos))
	}
}

func TestGetTodos_Unauthorized(t *testing.T) {
	// If AuthContext is not used, it should return 401
	r := chi.NewRouter()
	r.Get("/api/todos", todoHandler.GetTodos)

	req, _ := http.NewRequest("GET", "/api/todos", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rr.Code)
	}
}

func TestUpdateTodo(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupTodoRouter()

	var id int
	database.DB.QueryRow(`INSERT INTO todos (user_id, title, stage, completed) VALUES (1, 'Test Todo 1', 'To Do', false) RETURNING id`).Scan(&id)

	todo := models.Todo{
		Title: "Updated Todo",
		Stage: "In Progress",
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("PUT", "/api/todos/"+strconv.Itoa(id), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseTodo models.Todo
	json.NewDecoder(rr.Body).Decode(&responseTodo)
	if responseTodo.Title != "Updated Todo" {
		t.Errorf("expected title to be 'Updated Todo', got %v", responseTodo.Title)
	}
}

func TestToggleTodo(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupTodoRouter()

	var id int
	database.DB.QueryRow(`INSERT INTO todos (user_id, title, stage, completed) VALUES (1, 'Test Todo 1', 'To Do', false) RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("PATCH", "/api/todos/"+strconv.Itoa(id)+"/toggle", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseTodo models.Todo
	json.NewDecoder(rr.Body).Decode(&responseTodo)
	if responseTodo.Completed != true {
		t.Errorf("expected completed to be true, got %v", responseTodo.Completed)
	}
}

func TestDeleteTodo(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupTodoRouter()

	var id int
	database.DB.QueryRow(`INSERT INTO todos (user_id, title, stage, completed) VALUES (1, 'Test Todo 1', 'To Do', false) RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("DELETE", "/api/todos/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id = $1", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 todos, got %v", count)
	}
}

func TestCreateSubtask(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupTodoRouter()

	var parentID int
	database.DB.QueryRow(`INSERT INTO todos (user_id, title, stage, completed) VALUES (1, 'Parent', 'To Do', false) RETURNING id`).Scan(&parentID)

	todo := models.Todo{
		Title: "Subtask",
		Stage: "To Do",
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/api/todos/"+strconv.Itoa(parentID)+"/subtasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var responseTodo models.Todo
	json.NewDecoder(rr.Body).Decode(&responseTodo)
	if *responseTodo.ParentID != parentID {
		t.Errorf("expected parent ID to be %v, got %v", parentID, responseTodo.ParentID)
	}
}
