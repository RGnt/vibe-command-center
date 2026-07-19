package main

import (
	"log"
	"net/http"

	"os"

	"todo-backend/database"
	"todo-backend/handlers"
	mymiddleware "todo-backend/middleware"
	"todo-backend/repository"
	"todo-backend/service"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	// Initialize database
	database.InitDB()
	defer database.DB.Close()

	// Initialize dependencies
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "my_super_secret_key_change_me_in_prod" // Default for development
	}

	userRepo := repository.NewPostgresUserRepository(database.DB)
	authService := service.NewAuthService(userRepo, []byte(jwtSecret))
	authHandler := handlers.NewAuthHandler(authService)
	authMiddleware := mymiddleware.NewAuthMiddlewareProvider([]byte(jwtSecret))

	userSettingsRepo := repository.NewUserSettingsRepository(database.DB)
	userService := service.NewUserService(userSettingsRepo)
	userHandler := handlers.NewUserHandler(userService)

	workflowRepo := repository.NewWorkflowRepository(database.DB)
	workflowService := service.NewWorkflowService(workflowRepo)
	workflowHandler := handlers.NewWorkflowHandler(workflowService)

	todoRepo := repository.NewTodoRepository(database.DB)
	todoService := service.NewTodoService(todoRepo)
	todoHandler := handlers.NewTodoHandler(todoService)

	wikiRepo := repository.NewWikiRepository(database.DB)
	wikiService := service.NewWikiService(wikiRepo)
	wikiHandler := handlers.NewWikiHandler(wikiService)

	projectRepo := repository.NewProjectRepository(database.DB)
	projectService := service.NewProjectService(projectRepo, todoRepo, wikiRepo, workflowRepo)
	projectHandler := handlers.NewProjectHandler(projectService)

	diagramRepo := repository.NewDiagramRepository(database.DB)
	diagramService := service.NewDiagramService(diagramRepo)
	diagramHandler := handlers.NewDiagramHandler(diagramService)

	iconRepo := repository.NewIconRepository(database.DB)
	iconService := service.NewIconService(iconRepo)
	iconHandler := handlers.NewIconHandler(iconService)

	uploadHandler := handlers.NewUploadHandler(iconService)

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
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/logout", authHandler.Logout)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Middleware)

		r.Get("/api/auth/me", authHandler.GetMe)
		
		r.Get("/api/user/settings", userHandler.GetUserSettings)
		r.Put("/api/user/settings", userHandler.UpdateUserSettings)

		// Routes
		r.Get("/api/todos", todoHandler.GetTodos)
		r.Get("/api/todos/{id}", todoHandler.GetTodo)
		r.Post("/api/todos", todoHandler.CreateTodo)
		r.Post("/api/todos/{parent_id}/subtasks", todoHandler.CreateSubtask)
		r.Put("/api/todos/{id}", todoHandler.UpdateTodo)
		r.Delete("/api/todos/{id}", todoHandler.DeleteTodo)
		r.Patch("/api/todos/{id}/toggle", todoHandler.ToggleTodo)
		r.Get("/api/todos/stage", todoHandler.GetTodosByStage)

		// Workflow routes
		r.Get("/api/workflows", workflowHandler.GetWorkflows)
		r.Post("/api/workflows", workflowHandler.CreateWorkflow)
		r.Put("/api/workflows/{id}", workflowHandler.UpdateWorkflow)
		r.Delete("/api/workflows/{id}", workflowHandler.DeleteWorkflow)

		// Project routes
		r.Get("/api/projects", projectHandler.GetProjects)
		r.Get("/api/projects/{id}", projectHandler.GetProject)
		r.Get("/api/projects/{id}/export", projectHandler.ExportProject)
		r.Post("/api/projects/import", projectHandler.ImportProject)
		r.Post("/api/projects", projectHandler.CreateProject)
		r.Put("/api/projects/{id}", projectHandler.UpdateProject)
		r.Delete("/api/projects/{id}", projectHandler.DeleteProject)

		// Upload route
		r.Post("/api/upload", uploadHandler.UploadFile)

		// Icon routes
		r.Get("/api/icons", iconHandler.GetIcons)
		r.Put("/api/icons/{id}", iconHandler.UpdateIcon)
		r.Delete("/api/icons/{id}", iconHandler.DeleteIcon)

		// Diagram routes
		r.Get("/api/diagrams", diagramHandler.GetDiagrams)
		r.Get("/api/diagrams/{id}", diagramHandler.GetDiagram)
		r.Post("/api/diagrams", diagramHandler.CreateDiagram)
		r.Put("/api/diagrams/{id}", diagramHandler.UpdateDiagram)
		r.Delete("/api/diagrams/{id}", diagramHandler.DeleteDiagram)

		// Wiki routes
		r.Get("/api/wikis", wikiHandler.GetWikis)
		r.Get("/api/wikis/{slug}", wikiHandler.GetWiki)
		r.Post("/api/wikis", wikiHandler.CreateWiki)
		r.Put("/api/wikis/{id}", wikiHandler.UpdateWiki)
		r.Delete("/api/wikis/{id}", wikiHandler.DeleteWiki)
	})

	log.Println("Server starting on :8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
