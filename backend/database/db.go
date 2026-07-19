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
		stage TEXT DEFAULT 'To Do'
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
}
