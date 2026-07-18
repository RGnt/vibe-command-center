package service

import (
	"todo-backend/models"
	"todo-backend/repository"
)

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

func NewWorkflowService(workflowRepo repository.WorkflowRepository) WorkflowService {
	return &workflowService{
		workflowRepo: workflowRepo,
	}
}

func (s *workflowService) GetWorkflows(userID int) ([]models.Workflow, error) {
	return s.workflowRepo.GetAllByUserID(userID)
}

func (s *workflowService) GetWorkflow(id, userID int) (models.Workflow, error) {
	return s.workflowRepo.GetByIDAndUserID(id, userID)
}

func (s *workflowService) CreateWorkflow(userID int, workflow models.Workflow) (models.Workflow, error) {
	workflow.UserID = userID
	return s.workflowRepo.Create(workflow)
}

func (s *workflowService) UpdateWorkflow(id, userID int, workflow models.Workflow) (models.Workflow, error) {
	workflow.ID = id
	workflow.UserID = userID
	err := s.workflowRepo.Update(workflow)
	if err != nil {
		return workflow, err
	}
	return s.workflowRepo.GetByIDAndUserID(id, userID)
}

func (s *workflowService) DeleteWorkflow(id, userID int) error {
	return s.workflowRepo.Delete(id, userID)
}
