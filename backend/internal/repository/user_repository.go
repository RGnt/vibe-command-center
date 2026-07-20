package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"todo-backend/internal/models"
)

// ErrUserNotFound ...
var ErrUserNotFound = errors.New("user not found")
// ErrEmailExists ...
var ErrEmailExists = errors.New("email might already exist")

// UserRepository ...
type UserRepository interface {
	CreateUser(email, passwordHash string) (models.User, error)
	GetUserByEmail(email string) (models.User, string, error)
	GetUserByID(id int) (models.User, error)
}

type postgresUserRepository struct {
	db *sql.DB
}

// NewPostgresUserRepository ...
func NewPostgresUserRepository(db *sql.DB) UserRepository {
	return &postgresUserRepository{db: db}
}

// CreateUser ...
func (r *postgresUserRepository) CreateUser(email, passwordHash string) (models.User, error) {
	var user models.User

	// Begin transaction
	tx, err := r.db.Begin()
	if err != nil {
		return user, err
	}
	defer func() { _ = tx.Rollback() }()

	err = tx.QueryRow(`
		INSERT INTO users (email, password_hash)
		VALUES ($1, $2) RETURNING id, email, created_at
	`, email, passwordHash).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		return user, fmt.Errorf("%w: %w", ErrEmailExists, err)
	}

	// Create default user settings
	_, err = tx.Exec(`INSERT INTO user_settings (user_id, theme) VALUES ($1, 'system')`, user.ID)
	if err != nil {
		return user, err
	}

	// Create a default workflow
	var workflowID int
	err = tx.QueryRow(`INSERT INTO workflows (user_id, name) VALUES ($1, 'Default Workflow') RETURNING id`, user.ID).Scan(&workflowID)
	if err != nil {
		return user, err
	}

	stages := []string{"To Do", "In Progress", "Review", "Done"}
	for i, stageName := range stages {
		_, err = tx.Exec(`INSERT INTO workflow_stages (workflow_id, name, "order") VALUES ($1, $2, $3)`, workflowID, stageName, i+1)
		if err != nil {
			return user, err
		}
	}

	// Create a default project
	var projectID int
	err = tx.QueryRow(`INSERT INTO projects (user_id, name, description, workflow_id) VALUES ($1, 'General Project', 'Default project', $2) RETURNING id`, user.ID, workflowID).Scan(&projectID)
	if err != nil {
		return user, err
	}

	// Update user setting default project
	_, err = tx.Exec(`UPDATE user_settings SET default_project_id = $1 WHERE user_id = $2`, projectID, user.ID)
	if err != nil {
		return user, err
	}

	// Add a default wiki page
	_, err = tx.Exec(`INSERT INTO wiki_pages (user_id, project_id, category, title, slug, content) VALUES ($1, $2, 'General', 'Welcome to the Wiki', 'welcome', '# Welcome\n\nThis is your private wiki.')`, user.ID, projectID)
	if err != nil {
		return user, err
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return user, err
	}

	return user, nil
}

// GetUserByEmail ...
func (r *postgresUserRepository) GetUserByEmail(email string) (models.User, string, error) {
	var user models.User
	var hash string
	err := r.db.QueryRow(`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`, email).Scan(&user.ID, &user.Email, &hash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, "", ErrUserNotFound
		}
		return user, "", err
	}
	return user, hash, nil
}

// GetUserByID ...
func (r *postgresUserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User
	err := r.db.QueryRow(`SELECT id, email, created_at FROM users WHERE id = $1`, id).Scan(&user.ID, &user.Email, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, ErrUserNotFound
		}
		return user, err
	}
	return user, nil
}
