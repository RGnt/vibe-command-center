# Frontend Component Architecture

The frontend is a single-page application built with React and Vite, using Tailwind CSS for sophisticated, dark-themed Glassmorphism styling.

## High-Level Structure
- `App.jsx`: The root component. Handles conditional rendering based on authentication state.
- `main.jsx`: Bootstraps the React DOM.

## Context Providers
- `AuthContext.jsx`: Intercepts API requests to inject the Bearer token. Exposes user state, login, and registration methods globally.

## Core Feature Components

### Authentication
- `Login.jsx`: Form interface allowing user toggle between Registration and Sign In. Emits to `AuthContext`.

### Layout & Global UI
- `TodoApp.jsx`: The primary dashboard container. Houses the `Sidebar`, `Header`, and current view states.
- `Sidebar`: Navigation rail switching between Projects, Workflows, Wiki, Icons, and Diagrams.
- `Header`: Contextual breadcrumbs and global actions based on the current view.

### Kanban & Projects
- `ProjectsView`: Lists workspaces.
- `WorkflowsView`: Manages columns and logic for customizable boards.
- `KanbanBoard`: Core drag-and-drop UI component mapping Todos into Workflow stages. Provides inline editing, moving, and subtask generation.

### Wiki Integration
- `WikiView`: Renders markdown content dynamically using `react-markdown` or basic pre-formatted parsing. Provides an editor mode to alter raw markdown.

### Diagrams & Whiteboards
- `DiagramsView`: A dedicated pane for managing and saving graph/node diagram structures.

### Resources
- `IconsView`: Grid interface allowing users to upload (via multipart) or paste raw SVG code, structured into custom folders.

## Theming & Styling
- Pure utility classes are preferred using Tailwind CSS.
- **Glassmorphism**: Modals, forms, and cards use `backdrop-blur-xl`, semi-transparent backgrounds (`bg-[#1f2335]/80`), and distinct box shadows.
- Micro-interactions (e.g. `hover:-translate-y-0.5`, `transition-all duration-200`) are applied universally to buttons and interactive elements for a premium feel.
