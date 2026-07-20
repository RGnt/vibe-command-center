package service

import (
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// IconService ...
type IconService interface {
	GetIcons(userID int) ([]models.Icon, error)
	GetIcon(id string, userID int) (models.Icon, error)
	CreateIcon(userID int, icon models.Icon) (models.Icon, error)
	UpdateIcon(id string, userID int, icon models.Icon) (models.Icon, error)
	DeleteIcon(id string, userID int) error
}

type iconService struct {
	iconRepo repository.IconRepository
}

// NewIconService ...
func NewIconService(iconRepo repository.IconRepository) IconService {
	return &iconService{
		iconRepo: iconRepo,
	}
}

// GetIcons ...
func (s *iconService) GetIcons(userID int) ([]models.Icon, error) {
	return s.iconRepo.GetAllByUserID(userID)
}

// GetIcon ...
func (s *iconService) GetIcon(id string, userID int) (models.Icon, error) {
	return s.iconRepo.GetByIDAndUserID(id, userID)
}

// CreateIcon ...
func (s *iconService) CreateIcon(userID int, icon models.Icon) (models.Icon, error) {
	icon.UserID = userID
	return s.iconRepo.Create(icon)
}

// UpdateIcon ...
func (s *iconService) UpdateIcon(id string, userID int, icon models.Icon) (models.Icon, error) {
	icon.UserID = userID
	return s.iconRepo.Update(id, icon)
}

// DeleteIcon ...
func (s *iconService) DeleteIcon(id string, userID int) error {
	return s.iconRepo.Delete(id, userID)
}
