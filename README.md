# Agentic IDE Todo & Knowledge App

Welcome to the Agentic IDE Todo App! This is a feature-rich, modern application built with a React frontend, a Go backend, and a PostgreSQL database (equipped with pgvector and AGE extensions). It offers Kanban boards, a Mermaid diagram editor, a full Markdown Wiki, and a highly polished Tokyo Night aesthetic.

## Features

- **Kanban Board**: Visualize and manage tasks with stages ("To Do", "In Progress", "Done"). Support for nested subtasks.
- **Wiki System**: Create and edit Markdown-based wiki pages. Supports LaTeX equations and seamless internal linking.
- **Mermaid Editor**: Build architecture diagrams and flowcharts visually.
  - Choose from various diagram types (graph TD, sequence diagram, class diagram, etc.).
  - Drag and drop shapes or custom icons onto the canvas.
  - Upload custom icons (SVG, PNG, JPEG, WEBP) and organize them into collapsible folders.
  - Save diagrams to the database and embed them into the Wiki or Kanban tasks.
- **Tokyo Night Theme**: An elegant, dynamic visual design that emphasizes aesthetics and usability.
- **Dockerized Environment**: The entire stack, including a comprehensive E2E test suite using Playwright, runs inside Docker.

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/) and [Docker Compose](https://docs.docker.com/compose/)

### Running the Application

To start the main application stack, simply run:

```bash
docker compose up -d --build
```

This will spin up:
- The **Frontend** on `http://localhost:80`
- The **Backend API** (internally routed)
- The **PostgreSQL Database**
- **Adminer** on `http://localhost:8081` for direct database management

The application will be accessible at `http://localhost:80`.

### Running End-to-End (E2E) Tests

We have a complete E2E test suite configured using Playwright that tests the application in isolation.

1. Navigate to the `e2e` directory:
   ```bash
   cd e2e
   ```
2. Run the Playwright test suite (this will automatically build the test containers, run the tests, and tear down the environment):
   ```bash
   npx playwright test
   ```

*(Note: The E2E test environment runs on a separate port `8089` to avoid conflicting with the production container.)*

## Architecture

- **Frontend**: React, Vite, Tailwind CSS v4, React Flow (for canvas interactions), and Zustand (for state management).
- **Backend**: Go 1.23, utilizing a standard `net/http` router, and connecting to Postgres via the `pgx` driver.
- **Database**: PostgreSQL 16 with `pgvector` and `Apache AGE` (Graph Database extension) built-in. Data schemas include `todos`, `workflows`, `wikis`, `diagrams`, and `icons`.

## Documentation

For more detailed information about the API and system behaviors, check the `/documentation` directory in this repository:
- `/documentation/api/openapi.yaml` - Complete OpenAPI 3.0 specifications.
- `/documentation/specs/` - Behavior-Driven Development (BDD) `.feature` files describing all use cases for Todos, Workflows, Wikis, and the Diagram Editor.
