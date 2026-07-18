# General Application Architecture

KanbanX is structured as a modernized decoupled client-server architecture.

## Overview

The application utilizes three main containers running in Docker:
1. **Frontend**: An Nginx container serving a built Vite/React static application.
2. **Backend**: A compiled Go HTTP server utilizing `go-chi` for routing.
3. **Database**: A PostgreSQL database container enriched with `pgvector` and `AGE` graph extensions.

## Communication Pattern
The Client (React) issues RESTful HTTP calls. 
To bypass CORS complexity and unify the service layer, Nginx on the Frontend container acts as a reverse proxy for the backend API.
- Any request hitting `http://{frontend}/api/*` is seamlessly proxied to `http://backend:8080/*`.

## Backend Architecture Pattern
The Go backend strictly adheres to a three-tier layered architecture:
1. **Handlers (Controllers)**: Sits at the edge. Parses JSON requests, manages HTTP status codes, reads context (e.g., JWT credentials), and calls services.
2. **Services (Business Logic)**: Defines the core algorithms. Encrypts passwords, creates initial default data, handles logical validations, and coordinates repository calls.
3. **Repositories (Data Access)**: Isolates SQL. Executes transactions and builds safe queries via standard `database/sql`. Provides interfaces to make mocking simple.

## Security & State
State is ephemeral. 
- The backend relies entirely on **stateless JWT tokens** passed via the `Authorization: Bearer <token>` header. 
- Each request parses this token, and middleware (`AuthMiddleware`) rejects invalid or expired access directly, extracting the `user_id` into the request context.
- No session affinity is required.

## Testing Architecture
Tests are layered:
1. **Integration/Backend Unit Tests**: Written in Go standard `testing` framework utilizing the test database, ensuring route-level inputs hit database persistence correctly.
2. **E2E UI Tests**: Playwright scripts that build ephemeral, completely sandboxed instances of the three containers, simulating raw user clickpaths.
