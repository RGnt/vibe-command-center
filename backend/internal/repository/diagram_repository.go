package repository

import (
	"database/sql"
	"todo-backend/internal/models"
)

// DiagramRepository ...
type DiagramRepository interface {
	GetAllByUserID(userID int) ([]models.Diagram, error)
	GetByIDAndUserID(id string, userID int) (models.Diagram, error)
	Create(diagram models.Diagram) (models.Diagram, error)
	Update(id string, diagram models.Diagram) (models.Diagram, error)
	Delete(id string, userID int) error
}

type diagramRepository struct {
	db *sql.DB
}

// NewDiagramRepository ...
func NewDiagramRepository(db *sql.DB) DiagramRepository {
	return &diagramRepository{db: db}
}

// GetAllByUserID ...
func (r *diagramRepository) GetAllByUserID(userID int) ([]models.Diagram, error) {
	rows, err := r.db.Query(`SELECT id, user_id, name, diagram_type, code, explanation, created_at, updated_at FROM diagrams WHERE user_id = $1 ORDER BY updated_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var diagrams []models.Diagram
	for rows.Next() {
		var d models.Diagram
		if err := rows.Scan(&d.ID, &d.UserID, &d.Name, &d.DiagramType, &d.Code, &d.Explanation, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		diagrams = append(diagrams, d)
	}
	return diagrams, nil
}

// GetByIDAndUserID ...
func (r *diagramRepository) GetByIDAndUserID(id string, userID int) (models.Diagram, error) {
	var d models.Diagram
	err := r.db.QueryRow(`SELECT id, user_id, name, diagram_type, code, explanation, created_at, updated_at FROM diagrams WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&d.ID, &d.UserID, &d.Name, &d.DiagramType, &d.Code, &d.Explanation, &d.CreatedAt, &d.UpdatedAt)
	return d, err
}

// Create ...
func (r *diagramRepository) Create(diagram models.Diagram) (models.Diagram, error) {
	err := r.db.QueryRow(
		`INSERT INTO diagrams (user_id, name, diagram_type, code, explanation) VALUES ($1, $2, $3, $4, $5) RETURNING id, created_at, updated_at`,
		diagram.UserID, diagram.Name, diagram.DiagramType, diagram.Code, diagram.Explanation).
		Scan(&diagram.ID, &diagram.CreatedAt, &diagram.UpdatedAt)
	return diagram, err
}

// Update ...
func (r *diagramRepository) Update(id string, diagram models.Diagram) (models.Diagram, error) {
	err := r.db.QueryRow(
		`UPDATE diagrams SET name = $1, diagram_type = $2, code = $3, explanation = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 AND user_id = $6 RETURNING updated_at`,
		diagram.Name, diagram.DiagramType, diagram.Code, diagram.Explanation, id, diagram.UserID).
		Scan(&diagram.UpdatedAt)
	return diagram, err
}

// Delete ...
func (r *diagramRepository) Delete(id string, userID int) error {
	_, err := r.db.Exec(`DELETE FROM diagrams WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
