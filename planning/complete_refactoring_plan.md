# Complete Backend Refactoring Plan

## Goals
1. **Separation of Concerns:** Separate HTTP routing/handling, business logic, and database access into distinct layers.
2. **Dependency Injection:** Use struct-based dependencies rather than global package-level functions and variables.
3. **Security & Best Practices:** Abstract away sensitive configurations (like JWT keys) and improve testability by utilizing interfaces.

## Target Architecture
The backend will be refactored into three primary layers:
- **Handlers (Presentation Layer):** Responsible for parsing HTTP requests, delegating to the service layer, and formatting HTTP responses.
- **Services (Business Logic Layer):** Responsible for executing business rules, orchestrating complex operations, and calling repositories.
- **Repositories (Data Access Layer):** Responsible for executing raw SQL queries and interacting with the database.

## Phased Approach

### Phase 1: Foundation and Authentication (Completed)
1. **Create Packages:** Create `/repository` and `/service` directories.
2. **User/Auth Repository:** Create `UserRepository` for DB operations related to users (create user, fetch by email, set up default projects/workflows/wikis).
3. **Auth Service:** Create `AuthService` handling password hashing, password verification, and JWT generation.
4. **Auth Handlers:** Convert `Register`, `Login`, and `GetMe` to methods on an `AuthHandler` struct injected with `AuthService`.
5. **Middleware:** Refactor `AuthMiddleware` to accept dependencies (like the JWT secret) rather than using a global variable.
6. **Main Wiring:** Update `main.go` to instantiate the DB, Repository, Service, and Handler, then mount the new handler methods to the Chi router.

### Phase 2: User Settings and Profiles
1. **Repository:** Create `UserSettingsRepository`.
2. **Service:** Create `UserService` to handle fetching and updating user settings.
3. **Handlers:** Refactor `user_handlers.go` into a `UserHandler` struct injected with `UserService`. Update `main.go` routes.

### Phase 3: Core Workflow and Projects
1. **Repositories:** Create `ProjectRepository` and `WorkflowRepository`.
2. **Services:** Create `ProjectService` and `WorkflowService`. Move business logic (like ensuring a project belongs to the user, managing workflow stages) from handlers to these services.
3. **Handlers:** Refactor `project_handlers.go` and `workflow_handlers.go` into `ProjectHandler` and `WorkflowHandler` structs.

### Phase 4: Todos and Subtasks
1. **Repository:** Create `TodoRepository`. It should handle fetching todos, creating subtasks, toggling completion status, and fetching by stage.
2. **Service:** Create `TodoService`. This service will manage the complex logic of todo creation, stage transitions, and parent-child relationship validation.
3. **Handlers:** Refactor `todo_handlers.go` into a `TodoHandler` struct. This will be a significant refactor due to the complexity of the todo operations. Update related tests.

### Phase 5: Auxiliary Features (Diagrams, Icons, Wikis, Uploads)
1. **Repositories:** Create `DiagramRepository`, `IconRepository`, `WikiRepository`.
2. **Services:** Create `DiagramService`, `IconService`, `WikiService`, `UploadService`. Move file handling logic into `UploadService`.
3. **Handlers:** Refactor `diagram_handlers.go`, `icon_handlers.go`, `wiki_handlers.go`, and `upload_handlers.go` into their respective struct-based handlers. Update `main.go` with the final set of injected dependencies.

### Phase 6: Final Review and Cleanup
1. **Configuration Management:** Ensure all sensitive configurations (secrets, DB connection strings) are loaded via environment variables rather than hardcoded.
2. **Testing:** Verify all existing unit tests pass, and write new unit tests for the newly created services and repositories using mocked interfaces.
3. **Dead Code Elimination:** Remove any remaining global functions or unused structures from the original architecture.
