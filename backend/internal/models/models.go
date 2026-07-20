package models

import "time"

// User represents an authenticated user
type User struct {
	ID           int       `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // never serialize password hash
	CreatedAt    time.Time `json:"created_at"`
}

// UserSettings represents user-specific settings
type UserSettings struct {
	UserID           int    `json:"user_id"`
	Theme            string `json:"theme"`
	DefaultProjectID *int   `json:"default_project_id,omitempty"`
}

// Project represents a project containing workflows, tasks, and wikis
type Project struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	WorkflowID  *int      `json:"workflow_id,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
}

// WikiPage represents a single wiki document
type WikiPage struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
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
	UserID    int       `json:"user_id"`
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
	UserID    int             `json:"user_id"`
	Name      string          `json:"name"`
	Stages    []WorkflowStage `json:"stages"`
	CreatedAt time.Time       `json:"created_at"`
}

// Diagram represents a saved mermaid diagram
type Diagram struct {
	ID          int    `json:"id"`
	UserID      int    `json:"user_id"`
	Name        string `json:"name"`
	DiagramType string `json:"diagram_type"`
	Code        string `json:"code"`
	Explanation string `json:"explanation"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// Icon represents an uploaded icon image
type Icon struct {
	ID        int    `json:"id"`
	UserID    int    `json:"user_id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Folder    string `json:"folder"`
	CreatedAt string `json:"created_at"`
}

