package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

// ProjectService ...
type ProjectService interface {
	GetProjects(userID int) ([]models.Project, error)
	GetProject(id, userID int) (models.Project, error)
	CreateProject(userID int, project models.Project) (models.Project, error)
	UpdateProject(id, userID int, project models.Project) (models.Project, error)
	DeleteProject(id, userID int) error
	ExportProject(id, userID int) (models.ProjectExportPayload, error)
	ImportProject(userID int, payload models.ProjectExportPayload) (models.Project, error)
}

type projectService struct {
	projectRepo  repository.ProjectRepository
	todoRepo     repository.TodoRepository
	wikiRepo     repository.WikiRepository
	workflowRepo repository.WorkflowRepository
}

// NewProjectService ...
func NewProjectService(
	projectRepo repository.ProjectRepository,
	todoRepo repository.TodoRepository,
	wikiRepo repository.WikiRepository,
	workflowRepo repository.WorkflowRepository,
) ProjectService {
	return &projectService{
		projectRepo:  projectRepo,
		todoRepo:     todoRepo,
		wikiRepo:     wikiRepo,
		workflowRepo: workflowRepo,
	}
}

// GetProjects ...
func (s *projectService) GetProjects(userID int) ([]models.Project, error) {
	return s.projectRepo.GetAllByUserID(userID)
}

// GetProject ...
func (s *projectService) GetProject(id, userID int) (models.Project, error) {
	return s.projectRepo.GetByIDAndUserID(id, userID)
}

// CreateProject ...
func (s *projectService) CreateProject(userID int, project models.Project) (models.Project, error) {
	project.UserID = userID
	return s.projectRepo.Create(project)
}

// UpdateProject ...
func (s *projectService) UpdateProject(id, userID int, project models.Project) (models.Project, error) {
	project.ID = id
	project.UserID = userID
	err := s.projectRepo.Update(project)
	if err != nil {
		return project, err
	}
	return project, nil
}

// DeleteProject ...
func (s *projectService) DeleteProject(id, userID int) error {
	return s.projectRepo.Delete(id, userID)
}
