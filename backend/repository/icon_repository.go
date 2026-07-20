package repository

import (
	"database/sql"
	"todo-backend/models"
)

// IconRepository ...
type IconRepository interface {
	GetAllByUserID(userID int) ([]models.Icon, error)
	GetByIDAndUserID(id string, userID int) (models.Icon, error)
	Create(icon models.Icon) (models.Icon, error)
	Update(id string, icon models.Icon) (models.Icon, error)
	Delete(id string, userID int) error
}

type iconRepository struct {
	db *sql.DB
}

// NewIconRepository ...
func NewIconRepository(db *sql.DB) IconRepository {
	return &iconRepository{db: db}
}

// GetAllByUserID ...
func (r *iconRepository) GetAllByUserID(userID int) ([]models.Icon, error) {
	rows, err := r.db.Query(`SELECT id, user_id, name, url, folder, created_at FROM icons WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var icons []models.Icon
	for rows.Next() {
		var icon models.Icon
		if err := rows.Scan(&icon.ID, &icon.UserID, &icon.Name, &icon.URL, &icon.Folder, &icon.CreatedAt); err != nil {
			continue
		}
		icons = append(icons, icon)
	}
	return icons, nil
}

// GetByIDAndUserID ...
func (r *iconRepository) GetByIDAndUserID(id string, userID int) (models.Icon, error) {
	var icon models.Icon
	err := r.db.QueryRow(`SELECT id, user_id, name, url, folder, created_at FROM icons WHERE id = $1 AND user_id = $2`, id, userID).
		Scan(&icon.ID, &icon.UserID, &icon.Name, &icon.URL, &icon.Folder, &icon.CreatedAt)
	return icon, err
}

// Create ...
func (r *iconRepository) Create(icon models.Icon) (models.Icon, error) {
	err := r.db.QueryRow(
		`INSERT INTO icons (user_id, name, url, folder) VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		icon.UserID, icon.Name, icon.URL, icon.Folder).
		Scan(&icon.ID, &icon.CreatedAt)
	return icon, err
}

// Update ...
func (r *iconRepository) Update(id string, icon models.Icon) (models.Icon, error) {
	_, err := r.db.Exec(`UPDATE icons SET name = $1, folder = $2 WHERE id = $3 AND user_id = $4`, icon.Name, icon.Folder, id, icon.UserID)
	if err != nil {
		return icon, err
	}
	return r.GetByIDAndUserID(id, icon.UserID)
}

// Delete ...
func (r *iconRepository) Delete(id string, userID int) error {
	_, err := r.db.Exec(`DELETE FROM icons WHERE id = $1 AND user_id = $2`, id, userID)
	return err
}
