# Backend Models Specification

The system uses standard Go structs to model the database and core business entities. All models are defined in `backend/models/models.go`.

## Core Entities

### User
Represents an authenticated system user.
- **Fields**: `ID`, `Email`, `CreatedAt`
- **Purpose**: Tracks identity and timestamp for when the user joined.

### UserSettings
Represents application-level settings tailored to a user.
- **Fields**: `UserID`, `Theme` (default 'system'), `DefaultProjectID`
- **Purpose**: Stores preferences such as UI theme and the primary project space they land on.

### Project
A container for workflows, todos, wikis, and other resources.
- **Fields**: `ID`, `UserID`, `Name`, `Description`, `WorkflowID`, `CreatedAt`
- **Purpose**: Groups tasks and documentation under specific scopes.

### Workflow & WorkflowStage
Defines a custom kanban board layout.
- **Workflow**: `ID`, `UserID`, `Name`, `CreatedAt`
- **WorkflowStage**: `ID`, `WorkflowID`, `Name`, `Order`
- **Purpose**: Allows users to customize the stages of task progression (e.g., "To Do", "In Progress", "Done").

### Todo
A task to be completed.
- **Fields**: `ID`, `UserID`, `ProjectID`, `Title`, `Description`, `StageID`, `Order`, `IsCompleted`, `CreatedAt`, `UpdatedAt`
- **Purpose**: Tracks work items within projects, associated with specific workflow stages.

### Subtask
Smaller chunks of work belonging to a Todo.
- **Fields**: `ID`, `TodoID`, `Title`, `IsCompleted`
- **Purpose**: Granular task completion tracking.

### WikiPage
A markdown-based documentation page.
- **Fields**: `ID`, `UserID`, `ProjectID`, `Category`, `Title`, `Slug`, `Content`, `CreatedAt`, `UpdatedAt`
- **Purpose**: Provides rich-text documentation capabilities tied to projects.

### Diagram
A vector/graph representation created within the app.
- **Fields**: `ID`, `UserID`, `ProjectID`, `Name`, `Content`, `CreatedAt`, `UpdatedAt`
- **Purpose**: Stores raw diagram representations (e.g., JSON or XML schemas) for in-app rendering.

### Icon
A custom icon or graphical asset uploaded by the user.
- **Fields**: `ID`, `UserID`, `Folder`, `Name`, `Content`, `CreatedAt`
- **Purpose**: Stores SVG content or image URLs grouped by user-defined folders.
