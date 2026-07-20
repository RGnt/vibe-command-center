package repository

import (
	"todo-backend/internal/models"
)

// GetAllByProjectIDFlat ...
func (r *todoRepository) GetAllByProjectIDFlat(userID int, projectID int) ([]models.Todo, error) {
	rows, err := r.db.Query(`SELECT id, title, content, completed, created_at, parent_id, project_id, stage 
		FROM todos WHERE user_id = $1 AND project_id = $2 ORDER BY created_at ASC`, userID, projectID)
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
