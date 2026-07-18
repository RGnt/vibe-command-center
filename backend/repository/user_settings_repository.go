package repository

import (
	"database/sql"
	"todo-backend/models"
)

type UserSettingsRepository interface {
	GetByUserID(userID int) (models.UserSettings, error)
	Upsert(settings models.UserSettings) error
}

type userSettingsRepository struct {
	db *sql.DB
}

func NewUserSettingsRepository(db *sql.DB) UserSettingsRepository {
	return &userSettingsRepository{db: db}
}

func (r *userSettingsRepository) GetByUserID(userID int) (models.UserSettings, error) {
	var settings models.UserSettings
	err := r.db.QueryRow(`
		SELECT user_id, theme, default_project_id 
		FROM user_settings WHERE user_id = $1
	`, userID).Scan(&settings.UserID, &settings.Theme, &settings.DefaultProjectID)

	return settings, err
}

func (r *userSettingsRepository) Upsert(settings models.UserSettings) error {
	_, err := r.db.Exec(`
		INSERT INTO user_settings (user_id, theme, default_project_id) 
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id) DO UPDATE SET 
			theme = EXCLUDED.theme,
			default_project_id = EXCLUDED.default_project_id
	`, settings.UserID, settings.Theme, settings.DefaultProjectID)

	return err
}
