package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"todo-backend/database"
	"todo-backend/middleware"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

// GetTodos gets all top-level todos
func GetTodos(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	projectIDStr := r.URL.Query().Get("project_id")
	
	var rows *sql.Rows
	var err error

	if projectIDStr != "" {
		projectID, _ := strconv.Atoi(projectIDStr)
		rows, err = database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
			FROM todos WHERE parent_id IS NULL AND project_id = $1 AND user_id = $2 ORDER BY created_at DESC`, projectID, userID)
	} else {
		rows, err = database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
			FROM todos WHERE parent_id IS NULL AND user_id = $1 ORDER BY created_at DESC`, userID)
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
		var createdAt time.Time

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
		if err != nil {
			http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
			return
		}

		todo.UserID = userID
		todo.CreatedAt = createdAt

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
		subtasks, err := getSubtasks(todos[i].ID, userID)
		if err == nil {
			todos[i].Subtasks = subtasks
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todos)
}

// GetTodo gets a single todo by ID
func GetTodo(w http.ResponseWriter, r *http.Request) {
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

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt time.Time

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = $1 AND user_id = $2`, id, userID).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Todo not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch todo", http.StatusInternalServerError)
		}
		return
	}

	todo.UserID = userID
	todo.CreatedAt = createdAt

	if parentID.Valid {
		pid := int(parentID.Int64)
		todo.ParentID = &pid
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		todo.ProjectID = &pid
	}

	subtasks, err := getSubtasks(id, userID)
	if err == nil {
		todo.Subtasks = subtasks
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(todo)
}

// getSubtasks is a helper function to get subtasks
func getSubtasks(parentID int, userID int) ([]models.Todo, error) {
	rows, err := database.DB.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE parent_id = $1 AND user_id = $2 ORDER BY created_at DESC`, parentID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subtasks []models.Todo
	for rows.Next() {
		var todo models.Todo
		var projectID sql.NullInt64
		var createdAt time.Time

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &todo.ParentID, &projectID, &todo.Stage)
		if err != nil {
			return nil, err
		}

		todo.UserID = userID
		todo.CreatedAt = createdAt
		
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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	var id int64
	query := `INSERT INTO todos (user_id, title, content, stage, completed, project_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := database.DB.QueryRow(query, userID, todo.Title, todo.Content, todo.Stage, todo.Completed, todo.ProjectID).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create todo", http.StatusInternalServerError)
		return
	}

	var createdTodo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt time.Time

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = $1`, id).Scan(&createdTodo.ID, &createdTodo.Title, &createdTodo.Content, &createdTodo.Completed, &createdAt, &parentID, &projectID, &createdTodo.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch created todo", http.StatusInternalServerError)
		return
	}

	createdTodo.UserID = userID
	createdTodo.CreatedAt = createdAt

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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	var exists bool
	var parentProjectID sql.NullInt64
	err = database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1 AND user_id = $2), project_id FROM todos WHERE id = $3", parentID, userID, parentID).Scan(&exists, &parentProjectID)
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

	var id int64
	query := `INSERT INTO todos (user_id, title, content, parent_id, stage, completed, project_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = database.DB.QueryRow(query, userID, subtask.Title, subtask.Content, parentID, subtask.Stage, subtask.Completed, pID).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create subtask", http.StatusInternalServerError)
		return
	}

	var createdSubtask models.Todo
	var projectID sql.NullInt64
	var createdAt time.Time

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = $1`, id).Scan(&createdSubtask.ID, &createdSubtask.Title, &createdSubtask.Content, &createdSubtask.Completed, &createdAt, &createdSubtask.ParentID, &projectID, &createdSubtask.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch created subtask", http.StatusInternalServerError)
		return
	}

	createdSubtask.UserID = userID
	createdSubtask.CreatedAt = createdAt
	
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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `UPDATE todos SET title = $1, content = $2, completed = $3, stage = $4, project_id = $5 WHERE id = $6 AND user_id = $7`
	_, err = database.DB.Exec(query, updatedTodo.Title, updatedTodo.Content, updatedTodo.Completed, updatedTodo.Stage, updatedTodo.ProjectID, id, userID)
	if err != nil {
		http.Error(w, "Failed to update todo", http.StatusInternalServerError)
		return
	}

	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt time.Time

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = $1 AND user_id = $2`, id, userID).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch updated todo", http.StatusInternalServerError)
		return
	}

	todo.UserID = userID
	todo.CreatedAt = createdAt

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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	_, err = database.DB.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		http.Error(w, "Failed to delete todo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ToggleTodo toggles completion status
func ToggleTodo(w http.ResponseWriter, r *http.Request) {
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

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	var currentCompleted bool
	err = database.DB.QueryRow("SELECT completed FROM todos WHERE id = $1 AND user_id = $2", id, userID).Scan(&currentCompleted)
	if err != nil {
		http.Error(w, "Todo not found", http.StatusNotFound)
		return
	}

	newCompleted := !currentCompleted
	_, err = database.DB.Exec("UPDATE todos SET completed = $1 WHERE id = $2 AND user_id = $3", newCompleted, id, userID)
	if err != nil {
		http.Error(w, "Failed to toggle todo", http.StatusInternalServerError)
		return
	}

	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt time.Time

	err = database.DB.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = $1`, id).Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		http.Error(w, "Failed to fetch updated todo", http.StatusInternalServerError)
		return
	}

	todo.UserID = userID
	todo.CreatedAt = createdAt

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
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()
	
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
		FROM todos WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		http.Error(w, "Failed to fetch todos", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var todo models.Todo
		var parentID sql.NullInt64
		var projectID sql.NullInt64
		var createdAt time.Time

		err := rows.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
		if err != nil {
			http.Error(w, "Failed to scan todo", http.StatusInternalServerError)
			return
		}

		todo.UserID = userID
		todo.CreatedAt = createdAt

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
