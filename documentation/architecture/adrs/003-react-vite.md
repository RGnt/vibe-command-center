# ADR 003: React and Vite for Frontend

## Status
Accepted

## Context
We need a reactive frontend to support rich, interactive drag-and-drop kanban boards, diagram editing, and single-page navigation without full page reloads.

## Decision
We chose **React** as our UI library, scaffolded and built using **Vite**.

## Rationale
- **React**: Provides a vast ecosystem (e.g. `react-beautiful-dnd`, `lucide-react`) and components that map directly to application state.
- **Vite**: Replaces Create-React-App or Webpack. Provides near-instant Hot Module Replacement (HMR) during development and heavily optimizes ES modules during build.
- **Styling**: Tailwind CSS is used alongside React. This prevents CSS scoping issues and enables rapid layout iterations for our premium Glassmorphism aesthetic.

## Consequences
- Requires Javascript enabled on the client.
- The build step produces a static artifact that needs a lightweight web server (Nginx) for production.
