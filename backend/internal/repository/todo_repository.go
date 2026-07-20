package repository

import (
	"database/sql"
	"encoding/json"
	"time"
	"todo-backend/internal/models"
)

// TodoRepository ...
type TodoRepository interface {
	GetAllTopLevel(userID int, projectID *int) ([]models.Todo, error)
	GetByIDAndUserID(id, userID int) (models.Todo, error)
	GetSubtasks(parentID, userID int) ([]models.Todo, error)
	Create(todo models.Todo) (models.Todo, error)
	CreateSubtask(subtask models.Todo) (models.Todo, error)
	Update(todo models.Todo) error
	Delete(id, userID int) error
	ToggleCompleted(id, userID int) (bool, error)
	CheckExistsAndProjectID(id, userID int) (bool, *int, error)
	GetAllByUserID(userID int) ([]models.Todo, error)
	GetAllByProjectIDFlat(userID int, projectID int) ([]models.Todo, error)
}

type todoRepository struct {
	db *sql.DB
}

// NewTodoRepository ...
func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

// GetAllTopLevel ...
func (r *todoRepository) GetAllTopLevel(userID int, projectID *int) ([]models.Todo, error) {
	var rows *sql.Rows
	var err error

	if projectID != nil {
		rows, err = r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage, task_type, priority, custom_fields 
			FROM todos WHERE parent_id IS NULL AND project_id = $1 AND user_id = $2 ORDER BY created_at DESC`, *projectID, userID)
	} else {
		rows, err = r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage, task_type, priority, custom_fields 
			FROM todos WHERE parent_id IS NULL AND user_id = $1 ORDER BY created_at DESC`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var todos []models.Todo
	for rows.Next() {
		todo, err := scanTodo(rows, userID)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}

// GetByIDAndUserID ...
func (r *todoRepository) GetByIDAndUserID(id, userID int) (models.Todo, error) {
	row := r.db.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage, task_type, priority, custom_fields 
		FROM todos WHERE id = $1 AND user_id = $2`, id, userID)
	return scanTodoRow(row, userID)
}

// GetSubtasks ...
func (r *todoRepository) GetSubtasks(parentID, userID int) ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage, task_type, priority, custom_fields 
		FROM todos WHERE parent_id = $1 AND user_id = $2 ORDER BY created_at DESC`, parentID, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var subtasks []models.Todo
	for rows.Next() {
		todo, err := scanTodo(rows, userID)
		if err != nil {
			return nil, err
		}
		subtasks = append(subtasks, todo)
	}
	return subtasks, nil
}

// Create ...
func (r *todoRepository) Create(todo models.Todo) (models.Todo, error) {
	var id int64
	query := `INSERT INTO todos (user_id, title, content, stage, completed, project_id, task_type, priority, custom_fields) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9) RETURNING id`
	
	customFieldsJSON, _ := json.Marshal(todo.CustomFields)
	if todo.CustomFields == nil {
		customFieldsJSON = []byte("{}")
	}

	err := r.db.QueryRow(query, todo.UserID, todo.Title, todo.Content, todo.Stage, todo.Completed, todo.ProjectID, todo.TaskType, todo.Priority, customFieldsJSON).Scan(&id)
	if err != nil {
		return todo, err
	}
	
	return r.GetByIDAndUserID(int(id), todo.UserID)
}

// CreateSubtask ...
func (r *todoRepository) CreateSubtask(subtask models.Todo) (models.Todo, error) {
	var id int64
	query := `INSERT INTO todos (user_id, title, content, parent_id, stage, completed, project_id, task_type, priority, custom_fields) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`

	customFieldsJSON, _ := json.Marshal(subtask.CustomFields)
	if subtask.CustomFields == nil {
		customFieldsJSON = []byte("{}")
	}

	err := r.db.QueryRow(query, subtask.UserID, subtask.Title, subtask.Content, subtask.ParentID, subtask.Stage, subtask.Completed, subtask.ProjectID, subtask.TaskType, subtask.Priority, customFieldsJSON).Scan(&id)
	if err != nil {
		return subtask, err
	}

	return r.GetByIDAndUserID(int(id), subtask.UserID)
}

// Update ...
func (r *todoRepository) Update(todo models.Todo) error {
	query := `UPDATE todos SET title = $1, content = $2, completed = $3, stage = $4, project_id = $5, task_type = $6, priority = $7, custom_fields = $8 WHERE id = $9 AND user_id = $10`
	
	customFieldsJSON, _ := json.Marshal(todo.CustomFields)
	if todo.CustomFields == nil {
		customFieldsJSON = []byte("{}")
	}

	_, err := r.db.Exec(query, todo.Title, todo.Content, todo.Completed, todo.Stage, todo.ProjectID, todo.TaskType, todo.Priority, customFieldsJSON, todo.ID, todo.UserID)
	return err
}

// Delete ...
func (r *todoRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", id, userID)
	return err
}

// ToggleCompleted atomically flips the completed state of a todo.
func (r *todoRepository) ToggleCompleted(id, userID int) (bool, error) {
	var newCompleted bool
	err := r.db.QueryRow(
		`UPDATE todos SET completed = NOT completed WHERE id = $1 AND user_id = $2 RETURNING completed`,
		id, userID,
	).Scan(&newCompleted)
	return newCompleted, err
}

// CheckExistsAndProjectID ...
func (r *todoRepository) CheckExistsAndProjectID(id, userID int) (bool, *int, error) {
	var exists bool
	var parentProjectID sql.NullInt64
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM todos WHERE id = $1 AND user_id = $2), project_id FROM todos WHERE id = $3", id, userID, id).Scan(&exists, &parentProjectID)
	
	var pID *int
	if parentProjectID.Valid {
		pid := int(parentProjectID.Int64)
		pID = &pid
	}
	
	return exists, pID, err
}

// GetAllByUserID ...
func (r *todoRepository) GetAllByUserID(userID int) ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage, task_type, priority, custom_fields 
		FROM todos WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var todos []models.Todo
	for rows.Next() {
		todo, err := scanTodo(rows, userID)
		if err != nil {
			return nil, err
		}
		todos = append(todos, todo)
	}
	return todos, nil
}


// helper function for scanning
func scanTodo(scanner interface{ Scan(dest ...interface{}) error }, userID int) (models.Todo, error) {
	var todo models.Todo
	var parentID sql.NullInt64
	var projectID sql.NullInt64
	var createdAt time.Time
	var customFieldsJSON []byte

	err := scanner.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage, &todo.TaskType, &todo.Priority, &customFieldsJSON)
	if err != nil {
		return todo, err
	}

	todo.UserID = userID
	todo.CreatedAt = createdAt

	if len(customFieldsJSON) > 0 {
		_ = json.Unmarshal(customFieldsJSON, &todo.CustomFields)
	}
	if todo.CustomFields == nil {
		todo.CustomFields = make(map[string]interface{})
	}

	if parentID.Valid {
		pid := int(parentID.Int64)
		todo.ParentID = &pid
	}
	
	if projectID.Valid {
		pid := int(projectID.Int64)
		todo.ProjectID = &pid
	}

	return todo, nil
}

func scanTodoRow(row *sql.Row, userID int) (models.Todo, error) {
	return scanTodo(row, userID)
}
