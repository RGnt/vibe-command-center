package main

import (
	"log"
	"net/http"

	"todo-backend/database"
	"todo-backend/handlers"
	mymiddleware "todo-backend/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Initialize database
	database.InitDB()
	defer database.DB.Close()

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// CORS middleware
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	}))

	// Public routes
	r.Post("/api/auth/register", handlers.Register)
	r.Post("/api/auth/login", handlers.Login)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(mymiddleware.AuthMiddleware)

		r.Get("/api/auth/me", handlers.GetMe)
		
		r.Get("/api/user/settings", handlers.GetUserSettings)
		r.Put("/api/user/settings", handlers.UpdateUserSettings)

		// Routes
		r.Get("/api/todos", handlers.GetTodos)
		r.Get("/api/todos/{id}", handlers.GetTodo)
		r.Post("/api/todos", handlers.CreateTodo)
		r.Post("/api/todos/{parent_id}/subtasks", handlers.CreateSubtask)
		r.Put("/api/todos/{id}", handlers.UpdateTodo)
		r.Delete("/api/todos/{id}", handlers.DeleteTodo)
		r.Patch("/api/todos/{id}/toggle", handlers.ToggleTodo)
		r.Get("/api/todos/stage", handlers.GetTodosByStage)

		// Workflow routes
		r.Get("/api/workflows", handlers.GetWorkflows)
		r.Post("/api/workflows", handlers.CreateWorkflow)
		r.Put("/api/workflows/{id}", handlers.UpdateWorkflow)
		r.Delete("/api/workflows/{id}", handlers.DeleteWorkflow)

		// Project routes
		r.Get("/api/projects", handlers.GetProjects)
		r.Get("/api/projects/{id}", handlers.GetProject)
		r.Post("/api/projects", handlers.CreateProject)
		r.Put("/api/projects/{id}", handlers.UpdateProject)
		r.Delete("/api/projects/{id}", handlers.DeleteProject)

		// Upload route
		r.Post("/api/upload", handlers.UploadFile)

		// Icon routes
		r.Get("/api/icons", handlers.GetIcons)
		r.Put("/api/icons/{id}", handlers.UpdateIcon)
		r.Delete("/api/icons/{id}", handlers.DeleteIcon)

		// Diagram routes
		r.Get("/api/diagrams", handlers.GetDiagrams)
		r.Get("/api/diagrams/{id}", handlers.GetDiagram)
		r.Post("/api/diagrams", handlers.CreateDiagram)
		r.Put("/api/diagrams/{id}", handlers.UpdateDiagram)
		r.Delete("/api/diagrams/{id}", handlers.DeleteDiagram)

		// Wiki routes
		r.Get("/api/wikis", handlers.GetWikis)
		r.Get("/api/wikis/{slug}", handlers.GetWiki)
		r.Post("/api/wikis", handlers.CreateWiki)
		r.Put("/api/wikis/{id}", handlers.UpdateWiki)
		r.Delete("/api/wikis/{id}", handlers.DeleteWiki)
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
