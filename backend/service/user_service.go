package service

import (
	"database/sql"
	"todo-backend/models"
	"todo-backend/repository"
)

type UserService interface {
	GetSettings(userID int) (models.UserSettings, error)
	UpdateSettings(userID int, settings models.UserSettings) (models.UserSettings, error)
}

type userService struct {
	settingsRepo repository.UserSettingsRepository
}

func NewUserService(settingsRepo repository.UserSettingsRepository) UserService {
	return &userService{
		settingsRepo: settingsRepo,
	}
}

func (s *userService) GetSettings(userID int) (models.UserSettings, error) {
	settings, err := s.settingsRepo.GetByUserID(userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// If not found, return default
			return models.UserSettings{
				UserID: userID,
				Theme:  "system",
			}, nil
		}
		return settings, err
	}
	return settings, nil
}

func (s *userService) UpdateSettings(userID int, settings models.UserSettings) (models.UserSettings, error) {
	settings.UserID = userID
	err := s.settingsRepo.Upsert(settings)
	if err != nil {
		return settings, err
	}
	return settings, nil
}
