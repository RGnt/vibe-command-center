package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"todo-backend/database"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

func TestMain(m *testing.M) {
	database.InitTestDB()
	code := m.Run()
	database.DB.Close()
	os.Exit(code)
}

func clearTodosTable() {
	database.DB.Exec("DELETE FROM todos")
}

func setupRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/todos", GetTodos)
	r.Get("/api/todos/{id}", GetTodo)
	r.Post("/api/todos", CreateTodo)
	r.Put("/api/todos/{id}", UpdateTodo)
	r.Delete("/api/todos/{id}", DeleteTodo)
	r.Patch("/api/todos/{id}/toggle", ToggleTodo)
	return r
}

func TestCreateTodo(t *testing.T) {
	clearTodosTable()
	router := setupRouter()

	todo := models.Todo{
		Title:   "Test Todo",
		Content: "Test Content",
		Stage:   "To Do",
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("POST", "/api/todos", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
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
	clearTodosTable()
	router := setupRouter()

	database.DB.Exec(`INSERT INTO todos (title, content, stage, completed) VALUES ('Test Todo 1', 'Test Content', 'To Do', false)`)

	req, _ := http.NewRequest("GET", "/api/todos", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var todos []models.Todo
	json.NewDecoder(rr.Body).Decode(&todos)
	if len(todos) != 1 {
		t.Errorf("expected 1 todo, got %v", len(todos))
	}
}

func TestUpdateTodo(t *testing.T) {
	clearTodosTable()
	router := setupRouter()

	res, _ := database.DB.Exec(`INSERT INTO todos (title, stage, completed) VALUES ('Test Todo 1', 'To Do', false)`)
	id, _ := res.LastInsertId()

	todo := models.Todo{
		Title: "Updated Todo",
		Stage: "In Progress",
	}
	body, _ := json.Marshal(todo)
	req, _ := http.NewRequest("PUT", "/api/todos/"+strconv.Itoa(int(id)), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseTodo models.Todo
	json.NewDecoder(rr.Body).Decode(&responseTodo)
	if responseTodo.Title != "Updated Todo" {
		t.Errorf("expected title to be 'Updated Todo', got %v", responseTodo.Title)
	}
}

func TestToggleTodo(t *testing.T) {
	clearTodosTable()
	router := setupRouter()

	res, _ := database.DB.Exec(`INSERT INTO todos (title, stage, completed) VALUES ('Test Todo 1', 'To Do', false)`)
	id, _ := res.LastInsertId()

	req, _ := http.NewRequest("PATCH", "/api/todos/"+strconv.Itoa(int(id))+"/toggle", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var responseTodo models.Todo
	json.NewDecoder(rr.Body).Decode(&responseTodo)
	if responseTodo.Completed != true {
		t.Errorf("expected completed to be true, got %v", responseTodo.Completed)
	}
}

func TestDeleteTodo(t *testing.T) {
	clearTodosTable()
	router := setupRouter()

	res, _ := database.DB.Exec(`INSERT INTO todos (title, stage, completed) VALUES ('Test Todo 1', 'To Do', false)`)
	id, _ := res.LastInsertId()

	req, _ := http.NewRequest("DELETE", "/api/todos/"+strconv.Itoa(int(id)), nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM todos WHERE id = ?", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 todos, got %v", count)
	}
}
