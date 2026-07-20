# Golang Project Structure Guidelines

In Go, there is no single "official" project structure mandated by the compiler, but the community has coalesced around a few standard conventions (often inspired by `golang-standards/project-layout`). 

Our goal is to keep things simple, avoid over-engineering, but maintain a clean separation of concerns as the project grows.

## 1. Core Directories

### `/cmd`
This directory contains the main applications for the project. Each application should have its own subdirectory (e.g., `/cmd/server/main.go`, `/cmd/cli/main.go`). The `main.go` file should be small—its primary responsibility is configuration parsing, dependency injection, and starting the application. It should *not* contain core business logic.

### `/internal`
Private application and library code goes here. This is a special directory recognized by the Go compiler; packages inside `/internal` cannot be imported by applications outside of this repository. 
- Use this directory to hide the internal implementation details (e.g., `handlers`, `service`, `repository`, `models`) from the public API.
- This prevents accidental coupling if someone tries to use your backend as a library.

### `/pkg` (Optional)
Library code that is explicitly designed to be used by *other* external projects. If you are not building a shared library, do not use `/pkg`.

## 2. Package Organization

### Domain-Driven vs Layer-Driven
- **Layer-Driven (Current)**: Grouping files by technical layer (e.g., `/handlers`, `/service`, `/repository`). This is acceptable for small-to-medium web APIs.
- **Domain-Driven (Preferred for large projects)**: Grouping files by business domain (e.g., `/internal/auth`, `/internal/wiki`). 
- *Guideline*: Start with a layered approach (or flat), but as a module grows too large, consider breaking it into domain packages.

## 3. Applying the Standard Layout

To align a standard web backend with these principles:
1. Move the entry point (`main.go`) to `cmd/server/main.go`.
2. Move all business logic, routing, and data access layers (`handlers`, `middleware`, `models`, `repository`, `service`, `database`, `testutils`) into `internal/`.
3. Update all internal import paths to reflect the new `internal/` namespace.
