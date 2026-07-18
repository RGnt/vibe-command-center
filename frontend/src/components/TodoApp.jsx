import React, { useState, useEffect, useCallback } from 'react';
import KanbanBoard from './KanbanBoard';
import Sidebar from './Sidebar';
import ProjectModal from './ProjectModal';
import WikiView from './WikiView';
import MermaidEditor from './mermaid-editor/MermaidEditor';
import { useTheme } from '../contexts/ThemeContext';

const TodoApp = () => {
    const [projects, setProjects] = useState([]);
    const [activeProject, setActiveProject] = useState(null);
    const [activeTab, setActiveTab] = useState('board'); // 'board' or 'wiki'
    const [todos, setTodos] = useState([]);
    const [workflows, setWorkflows] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    const [showProjectModal, setShowProjectModal] = useState(false);
    
    const { theme, toggleTheme } = useTheme();

    // Fetch projects from backend
    const fetchProjects = async () => {
        try {
            const response = await fetch('/api/projects');
            if (!response.ok) throw new Error('Failed to fetch projects');
            const data = await response.json();
            setProjects(data || []);
            
            // Set first project as active if none selected
            if (data && data.length > 0 && !activeProject) {
                setActiveProject(data[0]);
            }
        } catch (err) {
            setError(err.message);
        }
    };

    // Fetch workflows
    const fetchWorkflows = async () => {
        try {
            const response = await fetch('/api/workflows');
            if (!response.ok) throw new Error('Failed to fetch workflows');
            const data = await response.json();
            setWorkflows(data || []);
        } catch (err) {
            setError(err.message);
        }
    };

    // Fetch todos for active project
    const fetchTodos = useCallback(async () => {
        if (!activeProject || activeProject === 'global-wiki') return;
        
        setLoading(true);
        try {
            const response = await fetch(`/api/todos?project_id=${activeProject.id}`);
            if (!response.ok) throw new Error('Failed to fetch todos');
            const data = await response.json();
            setTodos(data || []);
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    }, [activeProject]);

    // Initial load
    useEffect(() => {
        const loadInitialData = async () => {
            await fetchWorkflows();
            await fetchProjects();
            setLoading(false);
        };
        loadInitialData();
        
        // Listen for internal wiki navigation events (if they happen to come from somewhere else)
        const handleWikiNav = () => {
            setActiveTab('wiki');
        };
        
        window.addEventListener('wikiNavigate', handleWikiNav);
        return () => window.removeEventListener('wikiNavigate', handleWikiNav);
    }, []);

    // Load todos when active project changes
    useEffect(() => {
        if (activeProject && activeProject !== 'global-wiki') {
            fetchTodos();
            if (activeTab === 'wiki' && !projects.find(p => p.id === activeProject.id)) {
                // Keep wiki tab if project supports it
            }
        } else if (activeProject === 'global-wiki') {
            setActiveTab('wiki');
        }
    }, [activeProject, fetchTodos, activeTab, projects]);

    const handleCreateProject = async (projectData) => {
        try {
            const response = await fetch('/api/projects', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(projectData),
            });
            if (!response.ok) throw new Error('Failed to create project');
            const newProject = await response.json();
            setProjects([newProject, ...projects]);
            setActiveProject(newProject);
            setActiveTab('board');
        } catch (err) {
            setError(err.message);
        }
    };

    // Add a new todo
    const addTodo = async ({ title, stage = 'To Do' }) => {
        if (!title.trim() || !activeProject || activeProject === 'global-wiki') return;

        try {
            const response = await fetch('/api/todos', {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title, stage, project_id: activeProject.id }),
            });

            if (!response.ok) throw new Error('Failed to add todo');
            const newTodo = await response.json();
            setTodos([...todos, newTodo]);
        } catch (err) {
            setError(err.message);
        }
    };

    // Add a subtask
    const addSubtask = async (parentId, title) => {
        if (!title.trim() || activeProject === 'global-wiki') return;

        try {
            const response = await fetch(`/api/todos/${parentId}/subtasks`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ title, project_id: activeProject?.id }),
            });

            if (!response.ok) throw new Error('Failed to add subtask');
            const newSubtask = await response.json();

            setTodos(todos.map(todo => {
                if (todo.id === parentId) {
                    return {
                        ...todo,
                        subtasks: [...(todo.subtasks || []), newSubtask]
                    };
                }
                return todo;
            }));
        } catch (err) {
            setError(err.message);
        }
    };

    // Update a todo
    const updateTodo = async (id, updatedTodo) => {
        try {
            const response = await fetch(`/api/todos/${id}`, {
                method: 'PUT',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ ...updatedTodo, project_id: activeProject?.id }),
            });

            if (!response.ok) throw new Error('Failed to update todo');
            const updated = await response.json();
            setTodos(todos.map(todo => todo.id === id ? updated : todo));
        } catch (err) {
            setError(err.message);
        }
    };

    // Delete a todo
    const deleteTodo = async (id) => {
        try {
            const response = await fetch(`/api/todos/${id}`, {
                method: 'DELETE',
            });

            if (!response.ok) throw new Error('Failed to delete todo');
            setTodos(todos.filter(todo => todo.id !== id));
        } catch (err) {
            setError(err.message);
        }
    };

    if (loading && projects.length === 0) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-bg-base">
                <div className="text-xl font-semibold text-primary animate-pulse">Loading Workspace...</div>
            </div>
        );
    }

    if (error) {
        return (
            <div className="flex items-center justify-center min-h-screen bg-bg-base">
                <div className="text-xl font-semibold text-danger">Error: {error}</div>
            </div>
        );
    }

    return (
        <div className="w-full h-screen overflow-hidden flex bg-gradient-to-br from-bg-base via-bg-base to-bg-panel text-text-base">
            <Sidebar 
                projects={projects}
                activeProject={activeProject}
                onSelectProject={(project) => {
                    setActiveProject(project);
                    if (project === 'global-wiki') {
                        setActiveTab('wiki');
                    } else if (activeTab === 'wiki' && project !== 'global-wiki') {
                        setActiveTab('board'); // Default to board when switching projects
                    }
                }}
                onNewProject={() => setShowProjectModal(true)}
                onToggleTheme={toggleTheme}
                isDark={theme === 'dark'}
            />
            
            <main className="flex-1 flex flex-col h-full overflow-hidden">
                {activeProject === 'mermaid-editor' ? (
                    <div className="flex-1 overflow-hidden">
                        <MermaidEditor />
                    </div>
                ) : activeProject === 'global-wiki' ? (
                    <div className="flex-1 overflow-hidden p-6">
                        <WikiView project={null} isGlobal={true} />
                    </div>
                ) : activeProject ? (
                    <>
                        <header className="px-8 py-4 flex flex-col shrink-0 border-b border-border/40">
                            <h2 className="text-3xl font-bold text-text-base">{activeProject.name}</h2>
                            {activeProject.description && (
                                <p className="text-text-muted mt-1 text-sm">{activeProject.description}</p>
                            )}
                            
                            <div className="flex gap-4 mt-4">
                                <button 
                                    onClick={() => setActiveTab('board')}
                                    className={`px-4 py-1.5 rounded-lg text-sm font-semibold transition-all duration-200 ${
                                        activeTab === 'board' 
                                            ? 'bg-primary text-white shadow-md shadow-primary/20' 
                                            : 'text-text-muted hover:text-text-base hover:bg-bg-hover'
                                    }`}
                                >
                                    Board
                                </button>
                                <button 
                                    onClick={() => setActiveTab('wiki')}
                                    className={`px-4 py-1.5 rounded-lg text-sm font-semibold transition-all duration-200 ${
                                        activeTab === 'wiki' 
                                            ? 'bg-accent text-white shadow-md shadow-accent/20' 
                                            : 'text-text-muted hover:text-text-base hover:bg-bg-hover'
                                    }`}
                                >
                                    Wiki
                                </button>
                            </div>
                        </header>
                        <div className="flex-1 overflow-hidden p-6">
                            {activeTab === 'board' ? (
                                <KanbanBoard
                                    todos={todos}
                                    workflows={workflows}
                                    activeProjectWorkflowId={activeProject.workflow_id}
                                    onTodoUpdate={updateTodo}
                                    onTodoDelete={deleteTodo}
                                    onAddSubtask={addSubtask}
                                    onTodoCreate={addTodo}
                                />
                            ) : (
                                <WikiView project={activeProject} isGlobal={false} />
                            )}
                        </div>
                    </>
                ) : (
                    <div className="flex-1 flex items-center justify-center flex-col">
                        <div className="w-16 h-16 bg-bg-panel rounded-2xl flex items-center justify-center mb-4 shadow-sm border border-border">
                            <svg className="w-8 h-8 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                            </svg>
                        </div>
                        <h2 className="text-xl font-semibold text-text-base mb-2">No Project Selected</h2>
                        <p className="text-text-muted mb-6">Select a project from the sidebar or create a new one.</p>
                        <button
                            onClick={() => setShowProjectModal(true)}
                            className="px-6 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors shadow-sm"
                        >
                            Create Project
                        </button>
                    </div>
                )}
            </main>

            <ProjectModal 
                isOpen={showProjectModal}
                onClose={() => setShowProjectModal(false)}
                onSave={handleCreateProject}
                workflows={workflows}
            />
        </div>
    );
};

export default TodoApp;