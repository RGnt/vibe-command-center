# Frontend Architecture & Components

The frontend is a single-page application built with React, TypeScript, and Vite. It heavily utilizes modern state and routing libraries for a robust, developer-friendly experience.

## Core Technologies
- **Routing**: `@tanstack/react-router` for type-safe file-based routing.
- **State Management**: `zustand` for lightweight global state (e.g. `mermaidStore`).
- **Data Fetching**: `@tanstack/react-query` for server state caching and synchronization.
- **Styling**: Tailwind CSS v4 for utility-first styling with custom CSS variables (e.g., `--bg-base`, `--text-primary`).
- **Editors**: `@monaco-editor/react` configured with custom themes (Tokyo Night) and languages (Mermaid, Markdown).

## High-Level Structure
- `src/main.tsx`: Bootstraps the React DOM, wraps the app in QueryClientProvider, and initializes the Router.
- `src/routes/`: Contains all file-based routes (`__root.tsx`, `login.tsx`, `index.tsx`, `wiki.tsx`, `diagrams.tsx`, `agent.tsx`, etc.).
- `src/routeTree.gen.ts`: Automatically generated route tree managed by TanStack router.
- `src/contexts/AuthContext.tsx`: Manages JWT tokens, exposing user state, login, and registration methods globally.

## Core Feature Routes & Components

### Authentication
- `routes/login.tsx`: Form interface allowing user toggle between Registration and Sign In. Re-routes to the dashboard upon successful auth via `AuthContext`.

### Layout & Global UI
- `routes/__root.tsx`: The root wrapper that conditionally renders `AuthContext` checks and provides global layout.
- `components/layout/Sidebar.tsx`: Navigation rail switching between Projects, Workflows, Wiki, Icons, Diagrams, and Agents.
- `components/layout/Header.tsx`: Contextual breadcrumbs and user profile actions based on the current view.

### Kanban & Projects
- `routes/index.tsx`: The primary dashboard container displaying projects and workflows.
- `components/kanban/KanbanBoard.tsx`: Core drag-and-drop UI component mapping Todos into Workflow stages.

### Wiki Integration
- `routes/wiki.tsx`: Renders a split view of markdown pages.
- `components/wiki/WikiEditor.tsx`: Embeds Monaco Editor for authoring raw markdown.

### Diagrams & Whiteboards
- `routes/diagrams.tsx`: Renders the Mermaid diagram editor and preview.
- `components/mermaid-editor/EditorPreview.tsx`: Wraps `monaco-editor` customized for mermaid syntax.
- `components/mermaid-editor/MermaidRenderer.tsx`: Dynamically renders mermaid diagrams based on user code using standard libraries.

### Agents
- `routes/agent.tsx`: Interface to interact with the Python AI agent harness.
- `components/agent/AgentConfig.tsx`: Form to select wiki prompts, tasks, and display resulting code diffs returned by the agent using `@pierre/diffs/react`.

## Theming & Styling
- Pure utility classes are preferred using Tailwind CSS.
- **Tokyo Night Design System**: The UI utilizes deep, sophisticated dark blues, cyan, and magenta accents (matching the Tokyo Night IDE theme).
- Components utilize custom CSS variables mapped in `index.css` (e.g. `bg-bg-base`, `text-text-muted`) to ensure visual consistency and allow easy theme switching.
