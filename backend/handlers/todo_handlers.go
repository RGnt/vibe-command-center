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

// GetTodos gets all top-level todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	projectIDStr := r.URL.Query().Get("project_id")
	
	var rows *sql.Rows
	var err error

	if projectIDStr != "" {
		projectID, _ := strconv.Atoi(projectIDStr)
		rows, err = database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
			FROM todos WHERE parent_id IS NULL AND project_id = ? ORDER BY created_at DESC`, projectID)
	} else {
		rows, err = database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
			FROM todos WHERE parent_id IS NULL ORDER BY created_at DESC`)
	}

	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var todos []models.Todo
	for rows.Next() {
		var todo models.Todo
		var parentID sql.NullInt64
		var projectID sql.NullInt64
		var createdAt string

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
		if err != nil {
			http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
			return
		}

		t, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err == nil {
			todo.CreatedAt = t
		}

		if parentID.Valid {
			pid := int(parentID.Int64)
			todo.ParentID = &pid
		}
		
		if projectID.Valid {
			pid := int(projectID.Int64)
			todo.ProjectID = &pid
		}

		todos = append(todos, todo)
	}

	for i := range todos {
		subtasks, err := getSubtasks(todos[i].ID)
		if err == nil {
			todos[i].Subtasks = subtasks
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodo gets a single todo by ID
func GetTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = ?`, id).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		}
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		todo.CreatedAt = t
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		todo.ParentID = &pid
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		todo.ProjectID = &pid
	}

	subtasks, err := getSubtasks(id)
	if err == nil {
		todo.Subtasks = subtasks
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// getSubtasks is a helper function to get subtasks
func getSubtasks(parentID int) ([]models.Todo, error) {
	rows, err := database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE parent_id = ? ORDER BY created_at DESC`, parentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []models.Todo
	for rows.Next() {
		var todo models.Todo
		var projectID sql.NullInt64
		var createdAt string

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &todo.ParentID, &projectID, &todo.Stage)
		if err != nil {
			return nil, err
		}

		t, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err == nil {
			todo.CreatedAt = t
		}
		
		if projectID.Valid {
			pid := int(projectID.Int64)
			todo.ProjectID = &pid
		}

		subtasks = append(subtasks, todo)
	}

	return subtasks, nil
}

// CreateTodo creates a new todo
func CreateTodo(w http.ResponseWriter, r *http.Request) {
	var todo models.Todo
	if err := json.NewDecoder(r.Body).Decode(&todo); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `INSERT INTO todos (title, content, stage, completed, project_id) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, todo.Title, todo.Content, todo.Stage, todo.Completed, todo.ProjectID)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get todo ID", http.StatusInternalServerError)
		return
	}

	var createdTodo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = ?`, id).Scan(&createdTodo.ID, &createdTodo.Title, &createdTodo.Content, &createdTodo.Completed, &createdAt, &parentID, &projectID, &createdTodo.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch created todo", http.StatusInternalServerError)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		createdTodo.CreatedAt = t
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		createdTodo.ParentID = &pid
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		createdTodo.ProjectID = &pid
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdTodo)
}

// CreateSubtask creates a subtask for a todo
func CreateSubtask(w http.ResponseWriter, r *http.Request) {
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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	var exists bool
	var parentProjectID sql.NullInt64
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM todos WHERE id = ?), project_id FROM todos WHERE id = ?", parentID, parentID).Scan(&exists, &parentProjectID)
	if err != nil || !exists {
		http.Error(w, "Parent todo not found", http.StatusNotFound)
		return
	}

	var pID *int
	if parentProjectID.Valid {
		pid := int(parentProjectID.Int64)
		pID = &pid
	} else {
		pID = subtask.ProjectID
	}

	query := `INSERT INTO todos (title, content, parent_id, stage, completed, project_id) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, subtask.Title, subtask.Content, parentID, subtask.Stage, subtask.Completed, pID)
	if err != nil {
		http.Error(w, "Failed to create subtask", http.StatusInternalServerError)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to get subtask ID", http.StatusInternalServerError)
		return
	}

	var createdSubtask models.Todo
	var projectID sql.NullInt64
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = ?`, id).Scan(&createdSubtask.ID, &createdSubtask.Title, &createdSubtask.Content, &createdSubtask.Completed, &createdAt, &createdSubtask.ParentID, &projectID, &createdSubtask.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch created subtask", http.StatusInternalServerError)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		createdSubtask.CreatedAt = t
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		createdSubtask.ProjectID = &pid
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdSubtask)
}

// UpdateTodo updates a todo
func UpdateTodo(w http.ResponseWriter, r *http.Request) {
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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `UPDATE todos SET title = ?, content = ?, completed = ?, stage = ?, project_id = ? WHERE id = ?`
	_, err = database.DB.Exec(query, updatedTodo.Title, updatedTodo.Content, updatedTodo.Completed, updatedTodo.Stage, updatedTodo.ProjectID, id)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = ?`, id).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch updated todo", http.StatusInternalServerError)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		todo.CreatedAt = t
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		todo.ParentID = &pid
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		todo.ProjectID = &pid
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// DeleteTodo deletes a todo
func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	_, err = database.DB.Exec("DELETE FROM todos WHERE id = ?", id)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleTodo toggles completion status
func ToggleTodo(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid todo ID", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	var currentCompleted bool
	err = database.DB.QueryRow("SELECT completed FROM todos WHERE id = ?", id).Scan(&currentCompleted)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	newCompleted := !currentCompleted
	_, err = database.DB.Exec("UPDATE todos SET completed = ? WHERE id = ?", newCompleted, id)
	if err != nil {
		http.Error(w, "Failed to toggle todo", http.StatusInternalServerError)
		return
	}

	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt string

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = ?`, id).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch updated todo", http.StatusInternalServerError)
		return
	}

	t, err := time.Parse("2006-01-02 15:04:05", createdAt)
	if err == nil {
		todo.CreatedAt = t
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		todo.ParentID = &pid
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		todo.ProjectID = &pid
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// GetTodosByStage gets todos grouped by stage
func GetTodosByStage(w http.ResponseWriter, r *http.Request) {
	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	var stagesStr string
	err := database.DB.QueryRow(`SELECT GROUP_CONCAT(name ORDER BY "order") FROM workflow_stages ws 
		JOIN workflows w ON ws.workflow_id = w.id WHERE w.name = 'Default Workflow'`).Scan(&stagesStr)
	
	stages := []string{"To Do", "In Progress", "Review", "Done"}
	
	stageRows, err := database.DB.Query(`SELECT name FROM workflow_stages ws JOIN workflows w ON ws.workflow_id = w.id WHERE w.name = 'Default Workflow' ORDER BY "order"`)
	if err == nil {
		var dbStages []string
		for stageRows.Next() {
			var stage string
			if err := stageRows.Scan(&stage); err == nil {
				dbStages = append(dbStages, stage)
			}
		}
		stageRows.Close()
		if len(dbStages) > 0 {
			stages = dbStages
		}
	}

	stageMap := make(map[string][]models.Todo)
	for _, stage := range stages {
		stageMap[stage] = []models.Todo{}
	}

	rows, err := database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos ORDER BY created_at DESC`)
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		var parentID sql.NullInt64
		var projectID sql.NullInt64
		var createdAt string

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
		if err != nil {
			http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
			return
		}

		t, err := time.Parse("2006-01-02 15:04:05", createdAt)
		if err == nil {
			todo.CreatedAt = t
		}

		if parentID.Valid {
			pid := int(parentID.Int64)
			todo.ParentID = &pid
		}
		
		if projectID.Valid {
			pid := int(projectID.Int64)
			todo.ProjectID = &pid
		}

		if _, exists := stageMap[todo.Stage]; exists {
			stageMap[todo.Stage] = append(stageMap[todo.Stage], todo)
		} else {
			stageMap["To Do"] = append(stageMap["To Do"], todo)
		}
	}

	var result []map[string]interface{}
	for _, stage := range stages {
		result = append(result, map[string]interface{}{
			"name":  stage,
			"todos": stageMap[stage],
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}
