package database

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/lib/pq" // PostgreSQL driver
)

// DB is the global database connection
var DB *sql.DB

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

	// Initialize Apache AGE graph
	_, err = DB.Exec(`LOAD 'age';`)
	if err != nil {
		log.Println("Note: Failed to load age extension:", err)
	} else {
		// Set search path so we don't have to fully qualify cypher calls, but keep public first for CREATE TABLE
		_, _ = DB.Exec(`SET search_path = public, ag_catalog, "$user";`)
		
		// Create the graph
		// Apache AGE's create_graph throws an error if it already exists, so we just log it
		_, err = DB.Exec(`SELECT create_graph('knowledge_graph');`)
		if err != nil {
			log.Println("Note: knowledge_graph might already exist or failed to create:", err)
		}
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
	// We are going to add user_id to existing tables.
	// We will create the users table first.
	_, err := DB.Exec(`CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email TEXT UNIQUE NOT NULL,
		password_hash TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create users table:", err)
	}

	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS user_settings (
		user_id INTEGER PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
		theme TEXT DEFAULT 'system',
		default_project_id INTEGER
	)`)
	if err != nil {
		log.Fatal("Failed to create user_settings table:", err)
	}

	// Create workflows table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS workflows (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create workflows table:", err)
	}

	// Create projects table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS projects (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		description TEXT DEFAULT '',
		workflow_id INTEGER REFERENCES workflows(id) ON DELETE SET NULL,
		custom_field_schema JSONB DEFAULT '[]',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create projects table:", err)
	}

	// Create todos table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS todos (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		title TEXT NOT NULL,
		content TEXT DEFAULT '',
		completed BOOLEAN DEFAULT FALSE,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		parent_id INTEGER REFERENCES todos(id) ON DELETE CASCADE,
		project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
		stage TEXT DEFAULT 'To Do',
		task_type TEXT DEFAULT 'Task',
		priority TEXT DEFAULT 'Medium',
		custom_fields JSONB DEFAULT '{}'
	)`)
	if err != nil {
		log.Fatal("Failed to create todos table:", err)
	}

	// Create workflow_stages table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS workflow_stages (
		id SERIAL PRIMARY KEY,
		workflow_id INTEGER REFERENCES workflows(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		"order" INTEGER NOT NULL
	)`)
	if err != nil {
		log.Fatal("Failed to create workflow_stages table:", err)
	}

	// Create wiki_pages table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS wiki_pages (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		project_id INTEGER REFERENCES projects(id) ON DELETE CASCADE,
		parent_id INTEGER REFERENCES wiki_pages(id) ON DELETE CASCADE,
		category TEXT DEFAULT 'General',
		title TEXT NOT NULL,
		slug TEXT NOT NULL,
		content TEXT DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		UNIQUE(user_id, slug)
	)`)
	if err != nil {
		log.Fatal("Failed to create wiki_pages table:", err)
	}

	// Migrate existing wiki_pages
	_, err = DB.Exec(`ALTER TABLE wiki_pages ADD COLUMN IF NOT EXISTS parent_id INTEGER REFERENCES wiki_pages(id) ON DELETE CASCADE`)
	if err != nil {
		log.Fatal("Failed to alter wiki_pages table:", err)
	}

	// Create wiki_page_revisions table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS wiki_page_revisions (
		id SERIAL PRIMARY KEY,
		wiki_page_id INTEGER REFERENCES wiki_pages(id) ON DELETE CASCADE,
		content TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create wiki_page_revisions table:", err)
	}

	// Create icons table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS icons (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		url TEXT NOT NULL,
		folder TEXT DEFAULT 'General',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create icons table:", err)
	}

	// Create diagrams table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS diagrams (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		name TEXT NOT NULL,
		diagram_type TEXT DEFAULT 'graph TD',
		code TEXT DEFAULT '',
		explanation TEXT DEFAULT '',
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create diagrams table:", err)
	}

	// Create library_documents table
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS library_documents (
		id SERIAL PRIMARY KEY,
		user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
		original_name TEXT NOT NULL,
		filename TEXT NOT NULL,
		filepath TEXT NOT NULL,
		size BIGINT NOT NULL,
		mime_type TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create library_documents table:", err)
	}
	
	// Ensure user_id column exists if tables were already there from previous version
	tables := []string{"workflows", "projects", "todos", "wiki_pages", "icons", "diagrams", "library_documents"}
	for _, table := range tables {
		_, _ = DB.Exec(`ALTER TABLE ` + table + ` ADD COLUMN IF NOT EXISTS user_id INTEGER REFERENCES users(id) ON DELETE CASCADE`)
	}
	_, _ = DB.Exec(`ALTER TABLE icons ADD COLUMN IF NOT EXISTS folder TEXT DEFAULT 'General'`)

	// Migrate existing tables for Task extensions
	_, _ = DB.Exec(`ALTER TABLE projects ADD COLUMN IF NOT EXISTS custom_field_schema JSONB DEFAULT '[]'`)
	_, _ = DB.Exec(`ALTER TABLE todos ADD COLUMN IF NOT EXISTS task_type TEXT DEFAULT 'Task'`)
	_, _ = DB.Exec(`ALTER TABLE todos ADD COLUMN IF NOT EXISTS priority TEXT DEFAULT 'Medium'`)
	_, _ = DB.Exec(`ALTER TABLE todos ADD COLUMN IF NOT EXISTS custom_fields JSONB DEFAULT '{}'`)

	// Revoked JWT tokens blocklist (AUTH-01)
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS revoked_tokens (
		jti        TEXT PRIMARY KEY,
		expires_at TIMESTAMP NOT NULL,
		revoked_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create revoked_tokens table:", err)
	}

	// Refresh tokens (AUTH-09)
	_, err = DB.Exec(`CREATE TABLE IF NOT EXISTS refresh_tokens (
		id         SERIAL PRIMARY KEY,
		user_id    INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
		token_hash TEXT NOT NULL UNIQUE,
		expires_at TIMESTAMP NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		log.Fatal("Failed to create refresh_tokens table:", err)
	}
}
