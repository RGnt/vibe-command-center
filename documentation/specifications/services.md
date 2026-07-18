# Backend Services Specification

The business logic of the KanbanX application is centralized in the `backend/service` package. These services sit between the HTTP handlers and the data repositories.

## Core Services

### AuthService
- **Responsibilities**: Handles user registration, login authentication, and JWT generation.
- **Methods**: `Register`, `Login`, `GetUserByID`.
- **Dependencies**: Uses `UserRepository` to persist and verify identity credentials, and `bcrypt` for password hashing.

### UserService
- **Responsibilities**: Manages user-specific preferences and global state.
- **Methods**: `GetUserSettings`, `UpdateUserSettings`.
- **Dependencies**: Integrates with `UserSettingsRepository`.

### ProjectService
- **Responsibilities**: Provides CRUD operations for Projects. Ensures resources are appropriately bounded to specific users.
- **Methods**: `CreateProject`, `GetProjects`, `GetProjectByID`, `UpdateProject`, `DeleteProject`.

### WorkflowService
- **Responsibilities**: Handles the custom logic for Kanban board columns/stages.
- **Methods**: `CreateWorkflow`, `GetWorkflows`, `UpdateWorkflow`, `DeleteWorkflow`.
- **Specific Logic**: Ensures workflows inherently tie into stages (e.g., managing the order array and ensuring `stage_id` continuity).

### TodoService
- **Responsibilities**: Core kanban task interactions including completion toggles and stage movements.
- **Methods**: `CreateTodo`, `GetTodos`, `UpdateTodo`, `DeleteTodo`, `ToggleTodo`, `CreateSubtask`.
- **Dependencies**: Interacts heavily with `TodoRepository`. 

### WikiService
- **Responsibilities**: Document organization via markdown pages. Enforces distinct `slug` uniqueness within user scopes.
- **Methods**: `CreateWikiPage`, `GetWikiPages`, `GetWikiPageBySlug`, `UpdateWikiPage`, `DeleteWikiPage`.

### DiagramService
- **Responsibilities**: Manages vector/graph diagrams stored as content strings.
- **Methods**: `CreateDiagram`, `GetDiagrams`, `GetDiagram`, `UpdateDiagram`, `DeleteDiagram`.

### IconService
- **Responsibilities**: Manages custom SVG paths or asset metadata across virtual "folders".
- **Methods**: `CreateIcon`, `GetIcons`, `UpdateIcon`, `DeleteIcon`.
