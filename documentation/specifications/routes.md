# API Routes Specification

All routes are defined via `go-chi` and routed from `backend/main.go` using a central middleware stack, including `AuthMiddleware` for JWT-based endpoint protection.

## Authentication (`/api/auth`)
- `POST /register`: Accepts `email` and `password`. Returns JWT.
- `POST /login`: Accepts `email` and `password`. Returns JWT.
- `GET /me`: Returns details for the authenticated user.

## User Settings (`/api/user`)
- `GET /settings`: Retrieves the theme and default project ID.
- `PUT /settings`: Updates user global state preferences.

## Projects (`/api/projects`)
- `POST /`: Creates a project.
- `GET /`: Lists all projects for the user.
- `GET /{id}`: Gets detailed information on a single project.
- `PUT /{id}`: Updates name/description.
- `DELETE /{id}`: Deletes project.

## Workflows (`/api/workflows`)
- `POST /`: Creates a workflow with customizable stages.
- `GET /`: Lists user workflows.
- `PUT /{id}`: Updates a workflow's details/stages.
- `DELETE /{id}`: Deletes the workflow.

## Todos (`/api/todos`)
- `POST /`: Creates a new kanban task.
- `GET /`: Lists tasks (can be filtered by Project).
- `PUT /{id}`: Updates a task (e.g., changing its stage or order).
- `DELETE /{id}`: Removes a task.
- `PATCH /{id}/toggle`: Quickly marks a task as done/undone.
- `POST /{id}/subtasks`: Creates a new subtask attached to a todo.

## Wikis (`/api/wikis`)
- `POST /`: Creates a new markdown wiki page.
- `GET /`: Lists all wiki pages.
- `GET /{slug}`: Fetches a single wiki page by unique string literal.
- `PUT /{id}`: Edits the page content.
- `DELETE /{id}`: Deletes the page.

## Diagrams (`/api/diagrams`)
- `POST /`: Saves a new diagram.
- `GET /`: Lists all diagrams.
- `GET /{id}`: Gets diagram payload.
- `PUT /{id}`: Edits a diagram.
- `DELETE /{id}`: Deletes a diagram.

## Icons (`/api/icons`)
- `POST /`: Adds a custom icon/SVG.
- `GET /`: Lists all icons.
- `PUT /{id}`: Updates an icon payload.
- `DELETE /{id}`: Removes the icon.

## Uploads (`/api/upload`)
- `POST /icon`: Handles multipart form data for uploading physical icon files. Parses the request and passes it to the Icon service.
