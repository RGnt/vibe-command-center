package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todo-backend/database"
	"todo-backend/repository"
	"todo-backend/service"
	"todo-backend/testutils"

	"github.com/go-chi/chi/v5"
)

var (
	authHandler     *AuthHandler
	userHandler     *UserHandler
	projectHandler  *ProjectHandler
	workflowHandler *WorkflowHandler
	todoHandler     *TodoHandler
	wikiHandler     *WikiHandler
	iconHandler     *IconHandler
	diagramHandler  *DiagramHandler
	uploadHandler   *UploadHandler
)

func TestMain(m *testing.M) {
	testutils.SetupTestDB()
	
	// Setup dependencies for all tests here
	userRepo := repository.NewPostgresUserRepository(database.DB)
	authService := service.NewAuthService(userRepo, testutils.TestJWTSecret)
	authHandler = NewAuthHandler(authService)

	userSettingsRepo := repository.NewUserSettingsRepository(database.DB)
	userService := service.NewUserService(userSettingsRepo)
	userHandler = NewUserHandler(userService)

	workflowRepo := repository.NewWorkflowRepository(database.DB)
	workflowService := service.NewWorkflowService(workflowRepo)
	workflowHandler = NewWorkflowHandler(workflowService)

	todoRepo := repository.NewTodoRepository(database.DB)
	todoService := service.NewTodoService(todoRepo)
	todoHandler = NewTodoHandler(todoService)

	wikiRepo := repository.NewWikiRepository(database.DB)
	wikiService := service.NewWikiService(wikiRepo)
	wikiHandler = NewWikiHandler(wikiService)

	projectRepo := repository.NewProjectRepository(database.DB)
	projectService := service.NewProjectService(projectRepo, todoRepo, wikiRepo, workflowRepo)
	projectHandler = NewProjectHandler(projectService)

	iconRepo := repository.NewIconRepository(database.DB)
	iconService := service.NewIconService(iconRepo)
	iconHandler = NewIconHandler(iconService)

	diagramRepo := repository.NewDiagramRepository(database.DB)
	diagramService := service.NewDiagramService(diagramRepo)
	diagramHandler = NewDiagramHandler(diagramService)

	uploadHandler = NewUploadHandler(iconService)

	code := m.Run()
	testutils.ClearDB()
	os.Exit(code)
}

func setupAuthRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/auth/register", authHandler.Register)
	r.Post("/api/auth/login", authHandler.Login)
	return r
}

func TestRegister(t *testing.T) {
	testutils.ClearDB()
	r := setupAuthRouter()

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["user"] == nil {
		t.Errorf("Expected user object in response")
	}
	if resp["token"] == nil {
		t.Errorf("Expected token in response")
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	testutils.ClearDB()
	r := setupAuthRouter()

	payload := map[string]string{
		"email":    "test@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)

	// First request
	req1, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	rr1 := httptest.NewRecorder()
	r.ServeHTTP(rr1, req1)

	// Second request with same email
	req2, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	req2.Header.Set("Content-Type", "application/json")
	rr2 := httptest.NewRecorder()
	r.ServeHTTP(rr2, req2)

	if rr2.Code != http.StatusConflict {
		t.Errorf("Expected status code %d, got %d", http.StatusConflict, rr2.Code)
	}
}

func TestLogin(t *testing.T) {
	testutils.ClearDB()
	r := setupAuthRouter()

	// Register user first
	payload := map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	reqReg.Header.Set("Content-Type", "application/json")
	rrReg := httptest.NewRecorder()
	r.ServeHTTP(rrReg, reqReg)

	// Test Login
	reqLogin, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(body))
	reqLogin.Header.Set("Content-Type", "application/json")
	rrLogin := httptest.NewRecorder()
	r.ServeHTTP(rrLogin, reqLogin)

	if rrLogin.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rrLogin.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rrLogin.Body).Decode(&resp)

	if resp["token"] == nil {
		t.Errorf("Expected token in login response")
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	testutils.ClearDB()
	r := setupAuthRouter()

	// Register user first
	payload := map[string]string{
		"email":    "login@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(payload)
	reqReg, _ := http.NewRequest("POST", "/api/auth/register", bytes.NewBuffer(body))
	reqReg.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(httptest.NewRecorder(), reqReg)

	// Test Login with wrong password
	badPayload := map[string]string{
		"email":    "login@example.com",
		"password": "wrongpassword",
	}
	badBody, _ := json.Marshal(badPayload)
	reqLogin, _ := http.NewRequest("POST", "/api/auth/login", bytes.NewBuffer(badBody))
	reqLogin.Header.Set("Content-Type", "application/json")
	rrLogin := httptest.NewRecorder()
	r.ServeHTTP(rrLogin, reqLogin)

	if rrLogin.Code != http.StatusUnauthorized {
		t.Errorf("Expected status code %d, got %d", http.StatusUnauthorized, rrLogin.Code)
	}
}

func TestGetMe(t *testing.T) {
	testutils.ClearDB()

	// Instead of hitting the actual router, we can test GetMe using our AuthContext wrapper
	r := chi.NewRouter()
	r.Get("/api/auth/me", testutils.AuthContext(1, authHandler.GetMe))

	// Need a user in DB with ID 1
	database.DB.Exec("INSERT INTO users (id, email, password_hash) VALUES (1, 'me@example.com', 'hash')")

	req, _ := http.NewRequest("GET", "/api/auth/me", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["email"] != "me@example.com" {
		t.Errorf("Expected email me@example.com, got %v", resp["email"])
	}
}
