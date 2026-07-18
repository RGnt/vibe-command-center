package database

import (
	"database/sql"
	"log"
	"os"
	"sync"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection
var DB *sql.DB

// Mutex is used to synchronize database access
var Mutex sync.RWMutex

// InitDB initializes the database connection and creates tables
func InitDB() {
	var err error
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://todo_user:todo_password@localhost:5432/todo_db?sslmode=disable"
	}
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}

	// Create extensions for AGE and pgvector
	_, err = DB.Exec(`CREATE EXTENSION IF NOT EXISTS age;`)
	if err != nil {
		log.Println("Note: Failed to create age extension (might need superuser or not installed):", err)
	}
	_, err = DB.Exec(`CREATE EXTENSION IF NOT EXISTS vector;`)
	if err != nil {
		log.Println("Note: Failed to create vector extension:", err)
	}

	// Create tables
	createTables()
}

// InitTestDB initializes a database connection for testing
func InitTestDB() {
	// For testing, just connect to the same DB or a test DB if configured
	InitDB()
}

// createTables creates the necessary database tables
func createTables() {
	// Create workflows table first (since projects references it)
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS workflows (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create workflows table:", err)
	}

	// Create projects table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		workflow_id INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE SET NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create projects table:", err)
	}

	// Create todos table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		title TEXT NOT NULL,
		content TEXT DEFAULT '',
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		parent_id INTEGER,
		project_id INTEGER,
		stage TEXT DEFAULT 'To Do',
		FOREIGN KEY (parent_id) REFERENCES todos(id) ON DELETE CASCADE,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Failed to create todos table:", err)
	}

	// Create workflow_stages table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS workflow_stages (
		id SERIAL PRIMARY KEY,
		workflow_id INTEGER,
		name TEXT NOT NULL,
		"order" INTEGER NOT NULL,
		FOREIGN KEY (workflow_id) REFERENCES workflows(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Failed to create workflow_stages table:", err)
	}

	// Create wiki_pages table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS wiki_pages (
		id SERIAL PRIMARY KEY,
		project_id INTEGER,
		category TEXT DEFAULT 'General',
		title TEXT NOT NULL,
		slug TEXT NOT NULL UNIQUE,
		content TEXT DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (project_id) REFERENCES projects(id) ON DELETE CASCADE
	)`)
	if err != nil {
		log.Fatal("Failed to create wiki_pages table:", err)
	}

	// Insert default data if tables are empty
	insertDefaultData()
}

// insertDefaultData inserts default workflow, projects, and stages if empty
func insertDefaultData() {
	var workflowID int64
	var projectID int64

	// Check if workflows exist
	var workflowCount int
	err := DB.QueryRow("SELECT COUNT(*) FROM workflows").Scan(&workflowCount)
	if err != nil {
		log.Println("Error checking workflows:", err)
		return
	}

	if workflowCount == 0 {
		// Insert default workflow
		workflowQuery := `INSERT INTO workflows (name) VALUES ($1) RETURNING id`
		err := DB.QueryRow(workflowQuery, "Default Workflow").Scan(&workflowID)
		if err != nil {
			log.Println("Error inserting default workflow:", err)
			return
		}

		// Insert default stages
		stages := []string{"To Do", "In Progress", "Review", "Done"}
		for i, stageName := range stages {
			stageQuery := `INSERT INTO workflow_stages (workflow_id, name, "order") VALUES ($1, $2, $3)`
			_, err := DB.Exec(stageQuery, workflowID, stageName, i+1)
			if err != nil {
				log.Println("Error inserting stage:", err)
			}
		}
	} else {
		DB.QueryRow("SELECT id FROM workflows ORDER BY id LIMIT 1").Scan(&workflowID)
	}

	// Check if projects exist
	var projectCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM projects").Scan(&projectCount)
	if err != nil {
		log.Println("Error checking projects:", err)
		return
	}

	if projectCount == 0 {
		projectQuery := `INSERT INTO projects (name, description, workflow_id) VALUES ($1, $2, $3) RETURNING id`
		err := DB.QueryRow(projectQuery, "General Project", "Default project for tasks", workflowID).Scan(&projectID)
		if err != nil {
			log.Println("Error inserting default project:", err)
			return
		}
	} else {
		DB.QueryRow("SELECT id FROM projects ORDER BY id LIMIT 1").Scan(&projectID)
	}

	// Update any todos that don't have a project
	DB.Exec("UPDATE todos SET project_id = $1 WHERE project_id IS NULL", projectID)

	// Check if todos exist
	var todoCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM todos").Scan(&todoCount)
	if err != nil {
		log.Println("Error checking todos:", err)
		return
	}

	if todoCount == 0 {
		// Insert sample todos
		todos := []struct {
			title     string
			stage     string
			completed bool
		}{
			{"Research Go framework options", "To Do", false},
			{"Design database schema", "In Progress", false},
			{"Implement authentication", "Done", true},
			{"Write documentation", "Review", false},
		}

		for _, todo := range todos {
			query := `INSERT INTO todos (title, stage, completed, project_id) VALUES ($1, $2, $3, $4)`
			_, err := DB.Exec(query, todo.title, todo.stage, todo.completed, projectID)
			if err != nil {
				log.Println("Error inserting sample todo:", err)
			}
		}
	}

	// Check if wiki pages exist
	var wikiCount int
	err = DB.QueryRow("SELECT COUNT(*) FROM wiki_pages").Scan(&wikiCount)
	if err == nil && wikiCount == 0 {
		wikiQuery := `INSERT INTO wiki_pages (project_id, category, title, slug, content) VALUES ($1, $2, $3, $4, $5)`
		DB.Exec(wikiQuery, projectID, "General", "Welcome to the Wiki", "welcome", "# Welcome\n\nThis is the default index page for this project's wiki.\n\nYou can use markdown, tables, and math!\n\nLink to another page: [Architecture](#wiki:architecture)")
		DB.Exec(wikiQuery, projectID, "Engineering", "Architecture Overview", "architecture", "# Architecture\n\nWe use a simple Go backend with a React frontend.")
		
		// General global wiki (no project)
		wikiQueryNullProj := `INSERT INTO wiki_pages (category, title, slug, content) VALUES ($1, $2, $3, $4)`
		DB.Exec(wikiQueryNullProj, "Guides", "Global Guide", "global-guide", "# Global Guide\n\nThis wiki page does not belong to any specific project.")
	}
}
