package models

import "time"

// Project represents a project containing workflows, tasks, and wikis
type Project struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkflowID  *int      `json:"workflow_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// WikiPage represents a single wiki document
type WikiPage struct {
	ID        int       `json:"id"`
	ProjectID *int      `json:"project_id,omitempty"`
	Category  string    `json:"category"`
	Title     string    `json:"title"`
	Slug      string    `json:"slug"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Todo represents a single todo item
type Todo struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Completed bool      `json:"completed"`
	CreatedAt time.Time `json:"created_at"`
	ParentID  *int      `json:"parent_id,omitempty"`
	ProjectID *int      `json:"project_id,omitempty"`
	Subtasks  []Todo    `json:"subtasks,omitempty"`
	Stage     string    `json:"stage"`
}

// WorkflowStage represents a stage in the workflow
type WorkflowStage struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Order int    `json:"order"`
}

// Workflow represents the customizable workflow configuration
type Workflow struct {
	ID        int             `json:"id"`
	Name      string          `json:"name"`
	Stages    []WorkflowStage `json:"stages"`
	CreatedAt time.Time       `json:"created_at"`
}
