package service

import (
	"fmt"
	"time"
	"todo-backend/internal/models"
)

// ExportProject ...
func (s *projectService) ExportProject(id, userID int) (models.ProjectExportPayload, error) {
	var payload models.ProjectExportPayload

	// 1. Fetch Project
	project, err := s.projectRepo.GetByIDAndUserID(id, userID)
	if err != nil {
		return payload, err
	}
	payload.Project = project

	// 2. Fetch Workflow
	if project.WorkflowID != nil {
		workflow, err := s.workflowRepo.GetByIDAndUserID(*project.WorkflowID, userID)
		if err == nil {
			payload.Workflow = &workflow
		}
	}

	// 3. Fetch Todos (Flat fetch, then build tree)
	flatTodos, err := s.todoRepo.GetAllByProjectIDFlat(userID, id)
	if err != nil {
		return payload, err
	}

	// Build the nested structure
	payload.Todos = buildTodoTree(flatTodos, nil)

	// 4. Fetch Wikis
	wikis, err := s.wikiRepo.GetAllByProjectID(userID, id)
	if err != nil {
		return payload, err
	}
	payload.Wikis = wikis

	return payload, nil
}

func buildTodoTree(flatTodos []models.Todo, parentID *int) []models.Todo {
	var result []models.Todo
	for _, todo := range flatTodos {
		if (parentID == nil && todo.ParentID == nil) || (parentID != nil && todo.ParentID != nil && *todo.ParentID == *parentID) {
			todo.Subtasks = buildTodoTree(flatTodos, &todo.ID)
			result = append(result, todo)
		}
	}
	return result
}

// ImportProject ...
func (s *projectService) ImportProject(userID int, payload models.ProjectExportPayload) (models.Project, error) {
	// Start by creating workflow if it exists
	var workflowID *int
	if payload.Workflow != nil {
		wf := *payload.Workflow
		wf.UserID = userID
		// We clear the ID so it's created fresh
		wf.ID = 0
		for i := range wf.Stages {
			wf.Stages[i].ID = 0
		}
		newWf, err := s.workflowRepo.Create(wf)
		if err != nil {
			return models.Project{}, fmt.Errorf("failed to import workflow: %v", err)
		}
		wid := newWf.ID
		workflowID = &wid
	}

	// Create Project
	proj := payload.Project
	proj.ID = 0
	proj.UserID = userID
	proj.WorkflowID = workflowID
	// Optional: mark it as imported
	proj.Name = "[Imported] " + proj.Name

	newProj, err := s.projectRepo.Create(proj)
	if err != nil {
		return models.Project{}, fmt.Errorf("failed to create project: %v", err)
	}

	// Import Todos recursively
	for _, todo := range payload.Todos {
		if err := s.importTodoNode(userID, newProj.ID, nil, todo); err != nil {
			return newProj, fmt.Errorf("failed to import todo: %v", err)
		}
	}

	// Import Wikis
	for _, wiki := range payload.Wikis {
		wiki.ID = 0
		wiki.UserID = userID
		wiki.ProjectID = &newProj.ID
		wiki.Slug = fmt.Sprintf("%s-%d", wiki.Slug, time.Now().UnixNano()) // Avoid slug collisions
		
		_, err := s.wikiRepo.Create(wiki)
		if err != nil {
			// Log error but don't fail entire import
			fmt.Printf("failed to import wiki %s: %v\n", wiki.Title, err)
		}
	}

	return newProj, nil
}

func (s *projectService) importTodoNode(userID, projectID int, parentID *int, todo models.Todo) error {
	todo.ID = 0
	todo.UserID = userID
	todo.ProjectID = &projectID
	todo.ParentID = parentID
	
	var newTodo models.Todo
	var err error
	if parentID == nil {
		newTodo, err = s.todoRepo.Create(todo)
	} else {
		newTodo, err = s.todoRepo.CreateSubtask(todo)
	}
	
	if err != nil {
		return err
	}

	for _, sub := range todo.Subtasks {
		if err := s.importTodoNode(userID, projectID, &newTodo.ID, sub); err != nil {
			return err
		}
	}

	return nil
}
