package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

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

func NewIconService(iconRepo repository.IconRepository) IconService {
	return &iconService{
		iconRepo: iconRepo,
	}
}

func (s *iconService) GetIcons(userID int) ([]models.Icon, error) {
	return s.iconRepo.GetAllByUserID(userID)
}

func (s *iconService) GetIcon(id string, userID int) (models.Icon, error) {
	return s.iconRepo.GetByIDAndUserID(id, userID)
}

func (s *iconService) CreateIcon(userID int, icon models.Icon) (models.Icon, error) {
	icon.UserID = userID
	return s.iconRepo.Create(icon)
}

func (s *iconService) UpdateIcon(id string, userID int, icon models.Icon) (models.Icon, error) {
	icon.UserID = userID
	return s.iconRepo.Update(id, icon)
}

func (s *iconService) DeleteIcon(id string, userID int) error {
	return s.iconRepo.Delete(id, userID)
}
