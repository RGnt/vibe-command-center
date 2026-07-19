# Backend Specification

This document describes the complete architecture of the KanbanX Go backend, organized by layer. The backend follows a clean architecture with strict separation between HTTP handlers, business logic services, and data repositories.

---

## Architecture Overview

The request lifecycle flows through these layers:

`
HTTP Request
  -> main.go (Chi Router + Middleware)
    -> handlers/ (HTTP layer: parse request, call service)
      -> service/ (business logic)
        -> repository/ (data access, raw SQL)
          -> PostgreSQL
`

**Stack:**
- Language: Go 1.21+
- HTTP Router: go-chi/chi v5
- Database Driver: lib/pq (PostgreSQL)
- Authentication: JWT via golang-jwt/jwt v5
- Password Hashing: golang.org/x/crypto/bcrypt

---

## Middleware

### middleware/auth.go

**AuthMiddlewareProvider** wraps protected route groups. It reads a 	oken HttpOnly cookie, validates the JWT signature using the application jwtKey, and injects the resolved user_id into the request context.

| Item | Detail |
|------|--------|
| Cookie Name | token |
| JWT Algorithm | HS256 |
| Context Key | UserContextKey(user_id) |
| On Missing Cookie | 401 Unauthorized |
| On Invalid Token | 401 Unauthorized |

**Helper:** GetUserID(ctx context.Context) (int, bool) retrieves the user ID injected by the middleware. All protected handlers call this first.

---

## Data Models (models/models.go)

All Go structs used as both database row targets and JSON API responses.

| Model | Key Fields | Notes |
|-------|-----------|-------|
| User | ID, Email, CreatedAt | PasswordHash is tagged json: - and never exposed |
| UserSettings | UserID, Theme, DefaultProjectID | Theme defaults to system |
| Project | ID, UserID, Name, Description, WorkflowID, CreatedAt | WorkflowID is nullable (*int) |
| WikiPage | ID, UserID, ProjectID, Category, Title, Slug, Content, CreatedAt, UpdatedAt | ProjectID nullable; Slug unique per user |
| Todo | ID, UserID, Title, Content, Completed, Stage, ParentID, ProjectID, Subtasks, CreatedAt | ParentID nullable (top-level vs subtask) |
| WorkflowStage | ID, Name, Order | Embedded in Workflow.Stages |
| Workflow | ID, UserID, Name, Stages []WorkflowStage, CreatedAt | Stages sorted by Order |
| Diagram | ID, UserID, Name, DiagramType, Code, Explanation, CreatedAt, UpdatedAt | Code is raw Mermaid syntax |
| Icon | ID, UserID, Name, URL, Folder, CreatedAt | URL is relative path or base64 data URI |

### Export Models (models/export.go)

ProjectExportPayload contains Project, optional Workflow pointer, nested Todos tree, and WikiPages for import/export functionality.

---

## HTTP Handlers (handlers/)

All handlers follow this pattern:
1. Extract userID from context via middleware.GetUserID - return 401 if missing.
2. Parse and validate request body or URL params - return 400 on failure.
3. Call appropriate service method - return 4xx/5xx on error.
4. Write JSON response with correct status code.

### Auth Handlers (auth_handlers.go)

| Handler | Method | Route | Auth | Success Response |
|---------|--------|-------|------|-----------------|
| Register | POST | /api/auth/register | Public | 201 + sets token cookie; returns {token, user} |
| Login | POST | /api/auth/login | Public | 200 + sets token cookie; returns {token, user} |
| Logout | POST | /api/auth/logout | Public | 200 + clears token cookie (MaxAge: -1) |
| GetMe | GET | /api/auth/me | Protected | 200 + User object |

Cookie properties: HttpOnly: true, SameSite: Lax, Path: /, MaxAge: 604800 (1 week). Secure is currently false; set to true when deploying with HTTPS.

Error conditions:
- Register: 409 Conflict if email already exists; 500 on other errors.
- Login: 401 Unauthorized if credentials are invalid (service.ErrInvalidCredentials).

### User Handlers (user_handlers.go)

| Handler | Method | Route | Success Response |
|---------|--------|-------|-----------------|
| GetUserSettings | GET | /api/user/settings | 200 + UserSettings |
| UpdateUserSettings | PUT | /api/user/settings | 200 + updated UserSettings |

### Todo Handlers (todo_handlers.go)

| Handler | Method | Route | Notes |
|---------|--------|-------|-------|
| GetTodos | GET | /api/todos | Query param project_id to filter |
| GetTodo | GET | /api/todos/{id} | 404 if not found |
| CreateTodo | POST | /api/todos | 201 on success |
| UpdateTodo | PUT | /api/todos/{id} | Full replacement |
| DeleteTodo | DELETE | /api/todos/{id} | 204 No Content |
| ToggleTodo | PATCH | /api/todos/{id}/toggle | Flips completed flag; 404 if not found |
| CreateSubtask | POST | /api/todos/{parent_id}/subtasks | 404 if parent not found |
| GetTodosByStage | GET | /api/todos/stage | Returns map[string][]Todo; loads stages from user workflow if available |

### Workflow Handlers (workflow_handlers.go)

| Handler | Method | Route |
|---------|--------|-------|
| GetWorkflows | GET | /api/workflows |
| CreateWorkflow | POST | /api/workflows |
| UpdateWorkflow | PUT | /api/workflows/{id} |
| DeleteWorkflow | DELETE | /api/workflows/{id} - 204 |

### Project Handlers (project_handlers.go)

| Handler | Method | Route | Notes |
|---------|--------|-------|-------|
| GetProjects | GET | /api/projects | |
| GetProject | GET | /api/projects/{id} | |
| CreateProject | POST | /api/projects | 201 on success |
| UpdateProject | PUT | /api/projects/{id} | |
| DeleteProject | DELETE | /api/projects/{id} | 204 |
| ExportProject | GET | /api/projects/{id}/export | Returns ProjectExportPayload as JSON |
| ImportProject | POST | /api/projects/import | Accepts ProjectExportPayload, prefixes name with [Imported] |

### Wiki Handlers (wiki_handlers.go)

| Handler | Method | Route | Notes |
|---------|--------|-------|-------|
| GetWikis | GET | /api/wikis | Returns []WikiPage (empty array, never null) |
| GetWiki | GET | /api/wikis/{slug} | Looked up by slug string, not numeric ID |
| CreateWiki | POST | /api/wikis | 201 on success |
| UpdateWiki | PUT | /api/wikis/{id} | Uses numeric ID |
| DeleteWiki | DELETE | /api/wikis/{id} | 204 |

### Diagram Handlers (diagram_handlers.go)

| Handler | Method | Route |
|---------|--------|-------|
| GetDiagrams | GET | /api/diagrams |
| GetDiagram | GET | /api/diagrams/{id} |
| CreateDiagram | POST | /api/diagrams |
| UpdateDiagram | PUT | /api/diagrams/{id} |
| DeleteDiagram | DELETE | /api/diagrams/{id} |

### Icon Handlers (icon_handlers.go)

| Handler | Method | Route |
|---------|--------|-------|
| GetIcons | GET | /api/icons |
| UpdateIcon | PUT | /api/icons/{id} |
| DeleteIcon | DELETE | /api/icons/{id} |

### Upload Handlers (upload_handlers.go)

| Handler | Method | Route | Notes |
|---------|--------|-------|-------|
| UploadFile | POST | /api/upload | Multipart form; saves file, calls IconService.CreateIcon internally |

---

## Services (service/)

Services are the business logic layer. Each is defined as a Go interface + private struct implementation, enabling testability via mocking.

### AuthService (auth_service.go)

Interface methods: Register(email, password string), Login(email, password string), GetUserByID(id int)

- Register: Hashes password with bcrypt.DefaultCost, creates user via repository, generates 1-week HS256 JWT.
- Login: Fetches user by email, compares bcrypt hash. Returns ErrInvalidCredentials on mismatch.
- JWT Claims: { user_id: int, exp: time.Now + 7 days }, signed with jwtKey from environment.

### UserService (user_service.go)

Interface methods: GetUserSettings(userID int), UpdateUserSettings(userID int, settings models.UserSettings)

### ProjectService (project_service.go + project_export.go)

Interface methods: CreateProject, GetProjects, GetProjectByID, UpdateProject, DeleteProject, ExportProject, ImportProject

ExportProject logic:
1. Fetches project, workflow, flat todos, and wikis.
2. Builds nested Todo tree from flat list using recursive buildTodoTree().

ImportProject logic:
1. Creates workflow (with fresh IDs), creates project ([Imported] prefix), imports todos recursively preserving parent-child relationships, imports wikis with timestamp-disambiguated slugs.

### WorkflowService (workflow_service.go)

Interface methods: CreateWorkflow, GetWorkflows, UpdateWorkflow, DeleteWorkflow

### TodoService (todo_service.go)

Interface methods: GetTodos, GetTodo, CreateTodo, CreateSubtask, UpdateTodo, DeleteTodo, ToggleTodo, GetTodosByStage

### WikiService (wiki_service.go)

Interface methods: GetWikis, GetWiki(slug), CreateWiki, UpdateWiki, DeleteWiki

Slug uniqueness is enforced at the DB level (unique constraint).

### DiagramService (diagram_service.go)

Interface methods: GetDiagrams, GetDiagram, CreateDiagram, UpdateDiagram, DeleteDiagram

### IconService (icon_service.go)

Interface methods: GetIcons, CreateIcon, UpdateIcon, DeleteIcon

---

## Repositories (repository/)

Repositories contain raw SQL queries against PostgreSQL. All queries are user-scoped (always include WHERE user_id = N) to prevent cross-user data access.

| Repository File | Entity | Key Interface Methods |
|-----------------|--------|-----------------------|
| user_repository.go | User | CreateUser, GetUserByEmail, GetUserByID |
| user_settings_repository.go | UserSettings | GetByUserID, Upsert |
| project_repository.go | Project | Create, GetAllByUserID, GetByIDAndUserID, Update, Delete |
| workflow_repository.go | Workflow + WorkflowStage | Create, GetAllByUserID, GetByIDAndUserID, Update, Delete |
| todo_repository.go | Todo | Create, GetAllByUserID, GetByID, CreateSubtask, Update, Delete, Toggle, GetByStage |
| todo_repository_flat.go | Todo | GetAllByProjectIDFlat (returns flat list for export) |
| wiki_repository.go | WikiPage | Create, GetAllByUserID, GetBySlugAndUserID, Update, Delete |
| wiki_repository_flat.go | WikiPage | GetAllByProjectID (for export) |
| diagram_repository.go | Diagram | Create, GetAll, GetByID, Update, Delete |
| icon_repository.go | Icon | Create, GetAll, Update, Delete |

Error handling: Repositories return sql.ErrNoRows on not-found, and ErrUserNotFound (custom sentinel) from the user repository.

---

## Security Considerations

| Concern | Implementation |
|---------|---------------|
| Authentication | JWT stored in HttpOnly cookie; not accessible to JavaScript |
| Password Storage | bcrypt with default cost factor |
| User Isolation | All SQL queries include user_id predicate |
| CORS | go-chi/cors allows http:// and https:// with credentials |
| JWT Secret | Loaded from JWT_SECRET env var; falls back to hardcoded dev key |
| SQL Injection | Parameterized queries throughout all repository files |

> WARNING: The fallback jwtKey must be overridden in production via the JWT_SECRET environment variable. The cookie Secure flag is currently false; set it to true when serving over HTTPS.
