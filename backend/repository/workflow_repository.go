package repository

import (
	"database/sql"
	"time"
	"todo-backend/models"
)

type WorkflowRepository interface {
	GetAllByUserID(userID int) ([]models.Workflow, error)
	GetByIDAndUserID(id, userID int) (models.Workflow, error)
	Create(workflow models.Workflow) (models.Workflow, error)
	Update(workflow models.Workflow) error
	Delete(id, userID int) error
}

type workflowRepository struct {
	db *sql.DB
}

func NewWorkflowRepository(db *sql.DB) WorkflowRepository {
	return &workflowRepository{db: db}
}

func (r *workflowRepository) GetAllByUserID(userID int) ([]models.Workflow, error) {
	rows, err := r.db.Query(`SELECT id, name, created_at FROM workflows WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workflows []models.Workflow
	for rows.Next() {
		var workflow models.Workflow
		var createdAt time.Time

		err := rows.Scan(&workflow.ID, &workflow.Name, &createdAt)
		if err != nil {
			return nil, err
		}
		workflow.UserID = userID
		workflow.CreatedAt = createdAt

		stageRows, err := r.db.Query(`SELECT id, name, "order" FROM workflow_stages WHERE workflow_id = $1 ORDER BY "order"`, workflow.ID)
		if err != nil {
			return nil, err
		}
		defer stageRows.Close()

		var stages []models.WorkflowStage
		for stageRows.Next() {
			var stage models.WorkflowStage
			err := stageRows.Scan(&stage.ID, &stage.Name, &stage.Order)
			if err != nil {
				return nil, err
			}
			stages = append(stages, stage)
		}
		workflow.Stages = stages

		workflows = append(workflows, workflow)
	}

	return workflows, nil
}

func (r *workflowRepository) GetByIDAndUserID(id, userID int) (models.Workflow, error) {
	var workflow models.Workflow
	var createdAt time.Time

	err := r.db.QueryRow(`SELECT id, name, created_at FROM workflows WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&workflow.ID, &workflow.Name, &createdAt)
	if err != nil {
		return workflow, err
	}
	
	workflow.UserID = userID
	workflow.CreatedAt = createdAt

	stageRows, err := r.db.Query(`SELECT id, name, "order" FROM workflow_stages WHERE workflow_id = $1 ORDER BY "order"`, id)
	if err != nil {
		return workflow, err
	}
	defer stageRows.Close()

	var stages []models.WorkflowStage
	for stageRows.Next() {
		var stage models.WorkflowStage
		err := stageRows.Scan(&stage.ID, &stage.Name, &stage.Order)
		if err != nil {
			return workflow, err
		}
		stages = append(stages, stage)
	}
	workflow.Stages = stages

	return workflow, nil
}

func (r *workflowRepository) Create(workflow models.Workflow) (models.Workflow, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return workflow, err
	}

	var id int64
	query := `INSERT INTO workflows (user_id, name) VALUES ($1, $2) RETURNING id`
	err = tx.QueryRow(query, workflow.UserID, workflow.Name).Scan(&id)
	if err != nil {
		tx.Rollback()
		return workflow, err
	}

	for i, stage := range workflow.Stages {
		stageQuery := `INSERT INTO workflow_stages (workflow_id, name, "order") VALUES ($1, $2, $3)`
		_, err := tx.Exec(stageQuery, id, stage.Name, i+1)
		if err != nil {
			tx.Rollback()
			return workflow, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return workflow, err
	}

	return r.GetByIDAndUserID(int(id), workflow.UserID)
}

func (r *workflowRepository) Update(workflow models.Workflow) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec("UPDATE workflows SET name = $1 WHERE id = $2 AND user_id = $3", workflow.Name, workflow.ID, workflow.UserID)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec("DELETE FROM workflow_stages WHERE workflow_id = $1", workflow.ID)
	if err != nil {
		tx.Rollback()
		return err
	}

	for i, stage := range workflow.Stages {
		stageQuery := `INSERT INTO workflow_stages (workflow_id, name, "order") VALUES ($1, $2, $3)`
		_, err := tx.Exec(stageQuery, workflow.ID, stage.Name, i+1)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (r *workflowRepository) Delete(id, userID int) error {
	_, err := r.db.Exec("DELETE FROM workflows WHERE id = $1 AND user_id = $2", id, userID)
	return err
}
