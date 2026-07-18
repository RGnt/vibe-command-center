package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

type ProjectService interface {
	GetProjects(userID int) ([]models.Project, error)
	GetProject(id, userID int) (models.Project, error)
	CreateProject(userID int, project models.Project) (models.Project, error)
	UpdateProject(id, userID int, project models.Project) (models.Project, error)
	DeleteProject(id, userID int) error
}

type projectService struct {
	projectRepo repository.ProjectRepository
}

func NewProjectService(projectRepo repository.ProjectRepository) ProjectService {
	return &projectService{
		projectRepo: projectRepo,
	}
}

func (s *projectService) GetProjects(userID int) ([]models.Project, error) {
	return s.projectRepo.GetAllByUserID(userID)
}

func (s *projectService) GetProject(id, userID int) (models.Project, error) {
	return s.projectRepo.GetByIDAndUserID(id, userID)
}

func (s *projectService) CreateProject(userID int, project models.Project) (models.Project, error) {
	project.UserID = userID
	return s.projectRepo.Create(project)
}

func (s *projectService) UpdateProject(id, userID int, project models.Project) (models.Project, error) {
	project.ID = id
	project.UserID = userID
	err := s.projectRepo.Update(project)
	if err != nil {
		return project, err
	}
	return project, nil
}

func (s *projectService) DeleteProject(id, userID int) error {
	return s.projectRepo.Delete(id, userID)
}
