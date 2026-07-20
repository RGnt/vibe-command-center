package repository

import (
	"database/sql"
	"time"
	"todo-backend/models"
)

// ProjectRepository ...
type ProjectRepository interface {
	GetAllByUserID(userID int) ([]models.Project, error)
	GetByIDAndUserID(id, userID int) (models.Project, error)
	Create(project models.Project) (models.Project, error)
	Update(project models.Project) error
	Delete(id, userID int) error
}

type projectRepository struct {
	db *sql.DB
}

// NewProjectRepository ...
func NewProjectRepository(db *sql.DB) ProjectRepository {
	return &projectRepository{db: db}
}

// GetAllByUserID ...
func (r *projectRepository) GetAllByUserID(userID int) ([]models.Project, error) {
	rows, err := r.db.Query(`SELECT id, name, description, workflow_id, created_at FROM projects WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var projects []models.Project
	for rows.Next() {
		var project models.Project
		var workflowID sql.NullInt64
		var createdAt time.Time

		err := rows.Scan(&project.ID, &project.Name, &project.Description, &workflowID, &createdAt)
		if err != nil {
			return nil, err
		}
		project.UserID = userID
		project.CreatedAt = createdAt

		if workflowID.Valid {
			wid := int(workflowID.Int64)
			project.WorkflowID = &wid
		}

		projects = append(projects, project)
	}

	return projects, nil
}

// GetByIDAndUserID ...
func (r *projectRepository) GetByIDAndUserID(id, userID int) (models.Project, error) {
	var project models.Project
	var workflowID sql.NullInt64
	var createdAt time.Time

	err := r.db.QueryRow(`SELECT id, name, description, workflow_id, created_at FROM projects WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&project.ID, &project.Name, &project.Description, &workflowID, &createdAt)
	if err != nil {
		return project, err
	}
	
	project.UserID = userID
	project.CreatedAt = createdAt

	if workflowID.Valid {
		wid := int(workflowID.Int64)
		project.WorkflowID = &wid
	}

	return project, nil
}

// Create ...
func (r *projectRepository) Create(project models.Project) (models.Project, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return project, err
	}
	defer func() { _ = tx.Rollback() }()

	var id int64
	query := `INSERT INTO projects (user_id, name, description, workflow_id) VALUES ($1, $2, $3, $4) RETURNING id`
	err = tx.QueryRow(query, project.UserID, project.Name, project.Description, project.WorkflowID).Scan(&id)
	if err != nil {
		return project, err
	}

	project.ID = int(id)
	project.CreatedAt = time.Now()

	// Create a default wiki page for the new project
	_, err = tx.Exec(`INSERT INTO wiki_pages (user_id, project_id, category, title, slug, content) VALUES ($1, $2, 'General', 'Index', 'index', '# Index\n\nWelcome to your new project wiki.')`, project.UserID, project.ID)
	if err != nil {
		return project, err
	}

	if err := tx.Commit(); err != nil {
		return project, err
	}

	return project, nil
}

// Update ...
func (r *projectRepository) Update(project models.Project) error {
	query := `UPDATE projects SET name = $1, description = $2, workflow_id = $3 WHERE id = $4 AND user_id = $5`
	_, err := r.db.Exec(query, project.Name, project.Description, project.WorkflowID, project.ID, project.UserID)
	return err
}

// Delete ...
func (r *projectRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM projects WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
