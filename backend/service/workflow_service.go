package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

// WorkflowService ...
type WorkflowService interface {
	GetWorkflows(userID int) ([]models.Workflow, error)
	GetWorkflow(id, userID int) (models.Workflow, error)
	CreateWorkflow(userID int, workflow models.Workflow) (models.Workflow, error)
	UpdateWorkflow(id, userID int, workflow models.Workflow) (models.Workflow, error)
	DeleteWorkflow(id, userID int) error
}

type workflowService struct {
	workflowRepo repository.WorkflowRepository
}

// NewWorkflowService ...
func NewWorkflowService(workflowRepo repository.WorkflowRepository) WorkflowService {
	return &workflowService{
		workflowRepo: workflowRepo,
	}
}

// GetWorkflows ...
func (s *workflowService) GetWorkflows(userID int) ([]models.Workflow, error) {
	return s.workflowRepo.GetAllByUserID(userID)
}

// GetWorkflow ...
func (s *workflowService) GetWorkflow(id, userID int) (models.Workflow, error) {
	return s.workflowRepo.GetByIDAndUserID(id, userID)
}

// CreateWorkflow ...
func (s *workflowService) CreateWorkflow(userID int, workflow models.Workflow) (models.Workflow, error) {
	workflow.UserID = userID
	return s.workflowRepo.Create(workflow)
}

// UpdateWorkflow ...
func (s *workflowService) UpdateWorkflow(id, userID int, workflow models.Workflow) (models.Workflow, error) {
	workflow.ID = id
	workflow.UserID = userID
	err := s.workflowRepo.Update(workflow)
	if err != nil {
		return workflow, err
	}
	return s.workflowRepo.GetByIDAndUserID(id, userID)
}

// DeleteWorkflow ...
func (s *workflowService) DeleteWorkflow(id, userID int) error {
	return s.workflowRepo.Delete(id, userID)
}
