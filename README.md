# KanbanX

KanbanX is a comprehensive project management and workflow organization tool designed to streamline your development and personal tasks. It features a modern, responsive UI built with React, and a robust backend built with Go and PostgreSQL.

## Features

- **Authentication**: Secure JWT-based user authentication.
- **Project Management**: Create and manage multiple projects and workspaces.
- **Kanban Board**: Highly customizable kanban boards with drag-and-drop functionality for managing tasks and subtasks.
- **Custom Workflows**: Define custom workflow stages for your projects to match your team's specific needs.
- **Integrated Wiki**: Built-in Markdown-based wiki for project documentation, architecture notes, and general guides.
- **Diagrams**: Integrated diagram creation and visualization directly within your workspace.
- **File Uploads & Icons**: Manage custom icons and upload attachments to enrich your workspace.

## Technology Stack

### Frontend
- React (Vite)
- Tailwind CSS
- Context API for state management
- Drag and Drop interfaces

### Backend
- Go (Golang)
- `go-chi` router
- JWT for authentication
- PostgreSQL with `pgvector` and `AGE` extensions

## Prerequisites

To run this application, you must have Docker and Docker Compose installed on your system.

## Running the Application

### Production / Development Mode
You can spin up the application using Docker Compose. The `docker-compose.yml` provides a full environment including the frontend, backend, and PostgreSQL database.

```bash
docker compose up -d --build
```
Once the containers are running:
- **Frontend** is available at: `http://localhost:8080` (or `http://localhost` depending on your Nginx configuration)
- **Backend API** is internally routed via the frontend proxy at `/api/`

### Running the End-to-End Test Suite

KanbanX includes a comprehensive End-to-End (E2E) testing suite built with Playwright. The E2E tests spin up an isolated test environment using `docker-compose.test.yml`.

1. Navigate to the `e2e` directory:
   ```bash
   cd e2e
   ```
2. Install Playwright dependencies (if running for the first time):
   ```bash
   npm install
   npx playwright install
   ```
3. Run the test suite:
   ```bash
   npx playwright test
   ```
   *Note: This command will automatically build the test containers, run the Go backend tests, and then execute the Playwright UI tests against the isolated environment.*

## Project Structure

- `/frontend`: React application, UI components, and Vite configuration.
- `/backend`: Go application, REST API handlers, business logic, and database repositories.
- `/e2e`: Playwright test suite and test infrastructure configuration.
- `/documentation`: Feature specifications, system architecture, and API documentation.

## User Interface

KanbanX utilizes a sleek "Glassmorphism" design with a dark mode color palette tailored for developers and power users. Micro-animations and responsive layouts ensure a smooth experience across different devices.

## API Documentation
The application features a RESTful API. For detailed API endpoints and request/response schemas, refer to the documentation in the `/documentation` directory or browse the handler definitions in `/backend/handlers`.
