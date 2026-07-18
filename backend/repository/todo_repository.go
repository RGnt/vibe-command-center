package repository

import (
	"database/sql"
	"time"
	"todo-backend/models"
)

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
}

type todoRepository struct {
	db *sql.DB
}

func NewTodoRepository(db *sql.DB) TodoRepository {
	return &todoRepository{db: db}
}

func (r *todoRepository) GetAllTopLevel(userID int, projectID *int) ([]models.Todo, error) {
	var rows *sql.Rows
	var err error

	if projectID != nil {
		rows, err = r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
			FROM todos WHERE parent_id IS NULL AND project_id = $1 AND user_id = $2 ORDER BY created_at DESC`, *projectID, userID)
	} else {
		rows, err = r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
			FROM todos WHERE parent_id IS NULL AND user_id = $1 ORDER BY created_at DESC`, userID)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *todoRepository) GetByIDAndUserID(id, userID int) (models.Todo, error) {
	row := r.db.QueryRow(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE id = $1 AND user_id = $2`, id, userID)
	return scanTodoRow(row, userID)
}

func (r *todoRepository) GetSubtasks(parentID, userID int) ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE parent_id = $1 AND user_id = $2 ORDER BY created_at DESC`, parentID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

func (r *todoRepository) Create(todo models.Todo) (models.Todo, error) {
	var id int64
	query := `INSERT INTO todos (user_id, title, content, stage, completed, project_id) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := r.db.QueryRow(query, todo.UserID, todo.Title, todo.Content, todo.Stage, todo.Completed, todo.ProjectID).Scan(&id)
	if err != nil {
		return todo, err
	}
	
	return r.GetByIDAndUserID(int(id), todo.UserID)
}

func (r *todoRepository) CreateSubtask(subtask models.Todo) (models.Todo, error) {
	var id int64
	query := `INSERT INTO todos (user_id, title, content, parent_id, stage, completed, project_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := r.db.QueryRow(query, subtask.UserID, subtask.Title, subtask.Content, subtask.ParentID, subtask.Stage, subtask.Completed, subtask.ProjectID).Scan(&id)
	if err != nil {
		return subtask, err
	}

	return r.GetByIDAndUserID(int(id), subtask.UserID)
}

func (r *todoRepository) Update(todo models.Todo) error {
	query := `UPDATE todos SET title = $1, content = $2, completed = $3, stage = $4, project_id = $5 WHERE id = $6 AND user_id = $7`
	_, err := r.db.Exec(query, todo.Title, todo.Content, todo.Completed, todo.Stage, todo.ProjectID, todo.ID, todo.UserID)
	return err
}

func (r *todoRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = $1 AND user_id = $2", id, userID)
	return err
}

func (r *todoRepository) ToggleCompleted(id, userID int) (bool, error) {
	var currentCompleted bool
	err := r.db.QueryRow("SELECT completed FROM todos WHERE id = $1 AND user_id = $2", id, userID).Scan(&currentCompleted)
	if err != nil {
		return false, err
	}

	newCompleted := !currentCompleted
	_, err = r.db.Exec("UPDATE todos SET completed = $1 WHERE id = $2 AND user_id = $3", newCompleted, id, userID)
	return newCompleted, err
}

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

func (r *todoRepository) GetAllByUserID(userID int) ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

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

	err := scanner.Scan(&todo.ID, &todo.Title, &todo.Content, &todo.Completed, &createdAt, &parentID, &projectID, &todo.Stage)
	if err != nil {
		return todo, err
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

	return todo, nil
}

func scanTodoRow(row *sql.Row, userID int) (models.Todo, error) {
	return scanTodo(row, userID)
}
