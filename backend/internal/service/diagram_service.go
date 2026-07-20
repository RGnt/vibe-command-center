package service

import (
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// DiagramService ...
type DiagramService interface {
	GetDiagrams(userID int) ([]models.Diagram, error)
	GetDiagram(id string, userID int) (models.Diagram, error)
	CreateDiagram(userID int, diagram models.Diagram) (models.Diagram, error)
	UpdateDiagram(id string, userID int, diagram models.Diagram) (models.Diagram, error)
	DeleteDiagram(id string, userID int) error
}

type diagramService struct {
	diagramRepo repository.DiagramRepository
}

// NewDiagramService ...
func NewDiagramService(diagramRepo repository.DiagramRepository) DiagramService {
	return &diagramService{
		diagramRepo: diagramRepo,
	}
}

// GetDiagrams ...
func (s *diagramService) GetDiagrams(userID int) ([]models.Diagram, error) {
	return s.diagramRepo.GetAllByUserID(userID)
}

// GetDiagram ...
func (s *diagramService) GetDiagram(id string, userID int) (models.Diagram, error) {
	return s.diagramRepo.GetByIDAndUserID(id, userID)
}

// CreateDiagram ...
func (s *diagramService) CreateDiagram(userID int, diagram models.Diagram) (models.Diagram, error) {
	diagram.UserID = userID
	return s.diagramRepo.Create(diagram)
}

// UpdateDiagram ...
func (s *diagramService) UpdateDiagram(id string, userID int, diagram models.Diagram) (models.Diagram, error) {
	diagram.UserID = userID
	return s.diagramRepo.Update(id, diagram)
}

// DeleteDiagram ...
func (s *diagramService) DeleteDiagram(id string, userID int) error {
	return s.diagramRepo.Delete(id, userID)
}
