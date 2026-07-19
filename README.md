# Vibe Command Center

Vibe Command Center is a comprehensive project management and workflow organization tool designed to streamline your development and personal tasks. It features a modern, responsive UI built with React, a robust backend built with Go and PostgreSQL, and a powerful **AI Agent Harness** for intelligent document ingestion and task automation.

## Features

- **Authentication**: "Secure" JWT-based user authentication.
- **Project Management**: Create and manage multiple projects and workspaces.
- **Kanban Board**: Highly customizable kanban boards with drag-and-drop functionality for managing tasks and subtasks.
- **Custom Workflows**: Define custom workflow stages for your projects to match your team's specific needs.
- **Integrated Wiki**: Built-in Markdown-based wiki for project documentation, architecture notes, and general guides.
- **Diagrams**: Integrated diagram creation and visualization directly within your workspace.
- **Local Library & Ingestion**: Securely upload local documents (PDFs, Office docs, images) without publishing them publicly. Ingest these documents directly into the Global Wiki using the AI Agent Harness.
- **AI Agent Harness**: Local GPU-accelerated container leveraging PyTorch and `docling` to intelligently convert unstructured documents to Markdown, answer queries, and execute complex operations securely in sandboxed environments.

## Technology Stack

### Frontend
- React (Vite) & TypeScript
- Tailwind CSS
- Context API & TanStack Query for state management
- `react-resizable-panels` for modern split-pane layouts

### Backend
- Go (Golang)
- `go-chi` router
- JWT for authentication
- PostgreSQL with `pgvector` and `AGE` extensions

### Agent Harness
- Python 3.11 with FastAPI
- PyTorch (CUDA-enabled)
- `docling` for advanced document conversion
- `openai-agents` SDK for orchestrating AI tasks

## Prerequisites

To run this application, you must have Docker and Docker Compose installed on your system.
For hardware acceleration of document conversion, a CUDA-compatible GPU is highly recommended but not strictly required.

## Running the Application

### Production / Development Mode
You can spin up the application using Docker Compose. The `docker-compose.yml` provides a full environment including the frontend, backend, agent-harness, and PostgreSQL database.

```bash
docker compose up -d --build
```

Once the containers are running:
- **Frontend UI** is available at: `http://localhost:8080` (or `http://localhost` depending on your Nginx configuration)
- **Backend API** is internally routed via the frontend proxy at `/api/`
- **Agent Harness API** is securely routed internally to handle AI and ingestion workloads

### Running the End-to-End Test Suite

Vibe Command Center includes a comprehensive End-to-End (E2E) testing suite built with Playwright. The E2E tests spin up an isolated test environment using `docker-compose.test.yml`.

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
- `/agent-harness`: Python FastAPI service orchestrating the Sandboxed AI Agent execution and GPU document ingestion.
- `/e2e`: Playwright test suite and test infrastructure configuration.
- `/documentation`: Feature specifications, system architecture, and API documentation.

## User Interface

Vibe Command Center utilizes a sleek "Glassmorphism" design with a dark mode color palette tailored for developers and power users. Micro-animations and responsive layouts ensure a smooth experience across different devices.

## API Documentation
The application features a RESTful API. For detailed API endpoints and request/response schemas, refer to the documentation in the `/documentation` directory or browse the handler definitions in `/backend/handlers`.
