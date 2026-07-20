package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"todo-backend/internal/database"
	"todo-backend/internal/middleware"
	"todo-backend/internal/models"
	"todo-backend/internal/service"

	"github.com/go-chi/chi/v5"
)

// TodoHandler ...
type TodoHandler struct {
	todoService service.TodoService
}

// NewTodoHandler ...
func NewTodoHandler(todoService service.TodoService) *TodoHandler {
	return &TodoHandler{
		todoService: todoService,
	}
}

// GetTodos gets all top-level todos
func (h *TodoHandler) GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	projectIDStr := r.URL.Query().Get("project_id")
	var pID *int
	if projectIDStr != "" {
		projectID, _ := strconv.Atoi(projectIDStr)
		pID = &projectID
	}

	todos, err := h.todoService.GetTodos(userID, pID)
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}

	if todos == nil {
		todos = []models.Todo{}
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todos)
}

// GetTodo gets a single todo by ID
func (h *TodoHandler) GetTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.GetTodo(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}

// CreateTodo creates a new todo
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdTodo, err := h.todoService.CreateTodo(userID, todo)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdTodo)
}

// CreateSubtask creates a subtask for a todo
func (h *TodoHandler) CreateSubtask(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	parentID, err := strconv.Atoi(chi.URLParam(r, "parent_id"))
	if err != nil {
		http.Error(w, "Invalid parent ID", http.StatusBadRequest)
		return
	}

	var subtask models.Todo
	if err := json.NewDecoder(r.Body).Decode(&subtask); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	createdSubtask, err := h.todoService.CreateSubtask(userID, parentID, subtask)
	if err != nil {
		if err.Error() == "parent todo not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
		} else {
			http.Error(w, "Failed to create subtask", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(createdSubtask)
}

// UpdateTodo updates a todo
func (h *TodoHandler) UpdateTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	var updatedTodo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&updatedTodo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.UpdateTodo(id, userID, updatedTodo)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes a todo
func (h *TodoHandler) DeleteTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	err = h.todoService.DeleteTodo(id, userID)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleTodo toggles completion status
func (h *TodoHandler) ToggleTodo(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	todo, err := h.todoService.ToggleTodo(id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to toggle todo", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(todo)
}

// GetTodosByStage gets todos grouped by stage
func (h *TodoHandler) GetTodosByStage(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Wait, we need to fetch the user's workflow stages here.
	// We can do it using DB for now or inject workflowService.
	// For now, doing it via DB directly is a slight layer breach but acceptable since we don't have GetWorkflowService injected.
	// A better way is to pass workflowService to TodoService, or do it here. Let's just do a quick DB query here.
	stages := []string{"To Do", "In Progress", "Review", "Done"}
	
	stageRows, err := database.DB.Query(`SELECT name FROM workflow_stages ws JOIN workflows w ON ws.workflow_id = w.id WHERE w.user_id = $1 ORDER BY "order"`, userID)
	if err == nil {
		var dbStages []string
		for stageRows.Next() {
			var stage string
			if err := stageRows.Scan(&stage); err == nil {
				dbStages = append(dbStages, stage)
			}
		}
		_ = stageRows.Close()
		if len(dbStages) > 0 {
			stages = dbStages
		}
	}

	result, err := h.todoService.GetTodosByStage(userID, stages)
	if err != nil {
		http.Error(w, "Failed to fetch todos by stage", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(result)
}
