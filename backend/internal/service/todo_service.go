package service

import (
	"errors"
	"todo-backend/internal/models"
	"todo-backend/internal/repository"
)

// TodoService ...
type TodoService interface {
	GetTodos(userID int, projectID *int) ([]models.Todo, error)
	GetTodo(id, userID int) (models.Todo, error)
	CreateTodo(userID int, todo models.Todo) (models.Todo, error)
	CreateSubtask(userID, parentID int, subtask models.Todo) (models.Todo, error)
	UpdateTodo(id, userID int, todo models.Todo) (models.Todo, error)
	DeleteTodo(id, userID int) error
	ToggleTodo(id, userID int) (models.Todo, error)
	GetTodosByStage(userID int, stages []string) ([]map[string]interface{}, error)
}

type todoService struct {
	todoRepo repository.TodoRepository
}

// NewTodoService ...
func NewTodoService(todoRepo repository.TodoRepository) TodoService {
	return &todoService{
		todoRepo: todoRepo,
	}
}

// GetTodos ...
func (s *todoService) GetTodos(userID int, projectID *int) ([]models.Todo, error) {
	todos, err := s.todoRepo.GetAllTopLevel(userID, projectID)
	if err != nil {
		return nil, err
	}

	for i := range todos {
		subtasks, err := s.todoRepo.GetSubtasks(todos[i].ID, userID)
		if err == nil {
			todos[i].Subtasks = subtasks
		}
	}

	return todos, nil
}

// GetTodo ...
func (s *todoService) GetTodo(id, userID int) (models.Todo, error) {
	todo, err := s.todoRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return todo, err
	}

	subtasks, err := s.todoRepo.GetSubtasks(id, userID)
	if err == nil {
		todo.Subtasks = subtasks
	}

	return todo, nil
}

// CreateTodo ...
func (s *todoService) CreateTodo(userID int, todo models.Todo) (models.Todo, error) {
	todo.UserID = userID
	return s.todoRepo.Create(todo)
}

// CreateSubtask ...
func (s *todoService) CreateSubtask(userID, parentID int, subtask models.Todo) (models.Todo, error) {
	exists, parentProjectID, err := s.todoRepo.CheckExistsAndProjectID(parentID, userID)
	if err != nil || !exists {
		return subtask, errors.New("parent todo not found")
	}

	subtask.UserID = userID
	subtask.ParentID = &parentID
	
	if parentProjectID != nil {
		subtask.ProjectID = parentProjectID
	}

	return s.todoRepo.CreateSubtask(subtask)
}

// UpdateTodo ...
func (s *todoService) UpdateTodo(id, userID int, todo models.Todo) (models.Todo, error) {
	todo.ID = id
	todo.UserID = userID
	err := s.todoRepo.Update(todo)
	if err != nil {
		return todo, err
	}
	return s.todoRepo.GetByIDAndUserID(id, userID)
}

// DeleteTodo ...
func (s *todoService) DeleteTodo(id, userID int) error {
	return s.todoRepo.Delete(id, userID)
}

// ToggleTodo ...
func (s *todoService) ToggleTodo(id, userID int) (models.Todo, error) {
	_, err := s.todoRepo.ToggleCompleted(id, userID)
	if err != nil {
		return models.Todo{}, err
	}
	return s.todoRepo.GetByIDAndUserID(id, userID)
}

// GetTodosByStage ...
func (s *todoService) GetTodosByStage(userID int, stages []string) ([]map[string]interface{}, error) {
	stageMap := make(map[string][]models.Todo)
	for _, stage := range stages {
		stageMap[stage] = []models.Todo{}
	}

	todos, err := s.todoRepo.GetAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	for _, todo := range todos {
		if _, exists := stageMap[todo.Stage]; exists {
			stageMap[todo.Stage] = append(stageMap[todo.Stage], todo)
		} else {
			stageMap["To Do"] = append(stageMap["To Do"], todo) // default fallback
		}
	}

	var result []map[string]interface{}
	for _, stage := range stages {
		result = append(result, map[string]interface{}{
			"name":  stage,
			"todos": stageMap[stage],
		})
	}

	return result, nil
}
