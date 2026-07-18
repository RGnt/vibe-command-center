package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

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

func NewDiagramService(diagramRepo repository.DiagramRepository) DiagramService {
	return &diagramService{
		diagramRepo: diagramRepo,
	}
}

func (s *diagramService) GetDiagrams(userID int) ([]models.Diagram, error) {
	return s.diagramRepo.GetAllByUserID(userID)
}

func (s *diagramService) GetDiagram(id string, userID int) (models.Diagram, error) {
	return s.diagramRepo.GetByIDAndUserID(id, userID)
}

func (s *diagramService) CreateDiagram(userID int, diagram models.Diagram) (models.Diagram, error) {
	diagram.UserID = userID
	return s.diagramRepo.Create(diagram)
}

func (s *diagramService) UpdateDiagram(id string, userID int, diagram models.Diagram) (models.Diagram, error) {
	diagram.UserID = userID
	return s.diagramRepo.Update(id, diagram)
}

func (s *diagramService) DeleteDiagram(id string, userID int) error {
	return s.diagramRepo.Delete(id, userID)
}
