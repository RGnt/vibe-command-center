import React from 'react';
import { useAuth } from '../../contexts/AuthContext';
import { Link, useRouterState } from '@tanstack/react-router';
import { useSuspenseQuery } from '@tanstack/react-query';
import { projectsQueryOptions } from '../../utils/queries';

const Sidebar = ({ onNewProject, onImportProject, onToggleTheme, isDark }) => {
    const { logout } = useAuth();
    const routerState = useRouterState();
    const currentPath = routerState.location.pathname;

    const { data: projects = [] } = useSuspenseQuery(projectsQueryOptions);

    return (
        <div className="w-64 bg-bg-panel/90 backdrop-blur-xl border-r border-border h-full flex flex-col transition-all duration-300">
            <div className="p-6 border-b border-border/50">
                <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent flex items-center gap-2">
                    <svg className="w-6 h-6 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    Vibe Command Center
                </h1>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                <div>
                    <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">Global</div>
                    <Link
                        to="/"
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            currentPath === '/'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${currentPath === '/' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                        </svg>
                        <span>Dashboard</span>
                    </Link>
                    <Link
                        to="/wiki"
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            currentPath === '/wiki'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${currentPath === '/wiki' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                        </svg>
                        <span>General Wiki</span>
                    </Link>
                </div>

                <div>
                    <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">Tools</div>
                    <Link
                        to="/mermaid-editor"
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            currentPath === '/mermaid-editor'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${currentPath === '/mermaid-editor' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                        <span>Mermaid Editor</span>
                    </Link>
                    <Link
                        to="/agent"
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            currentPath === '/agent'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${currentPath === '/agent' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9.75 17L9 20l-1 1h8l-1-1-.75-3M3 13h18M5 17h14a2 2 0 002-2V5a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                        </svg>
                        <span>Agent Harness</span>
                    </Link>
                    <Link
                        to="/library"
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            currentPath === '/library'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${currentPath === '/library' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                        </svg>
                        <span>Local Library</span>
                    </Link>
                    <Link 
                        to="/graph" 
                        className={`flex items-center gap-3 px-3 py-2 rounded-md transition-colors text-sm font-medium
                            ${currentPath === '/graph' 
                                ? 'bg-surface-active text-text-primary' 
                                : 'text-text-secondary hover:bg-surface-hover hover:text-text-primary'}`}
                    >
                        <svg className={`w-4 h-4 ${currentPath === '/graph' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13.828 10.172a4 4 0 00-5.656 0l-4 4a4 4 0 105.656 5.656l1.102-1.101m-.758-4.899a4 4 0 005.656 0l4-4a4 4 0 00-5.656-5.656l-1.1 1.1" />
                        </svg>
                        <span>Knowledge Graph</span>
                    </Link>
                </div>

                <div>
                    <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">Projects</div>
                
                {projects.map(project => (
                    <div key={project.id} className="mb-1">
                        <Link
                            to={`/projects/${project.id}`}
                            className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                                currentPath === `/projects/${project.id}` || currentPath === `/projects/${project.id}/`
                                    ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                    : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                            }`}
                        >
                            <svg className={`w-4 h-4 ${currentPath === `/projects/${project.id}` ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                            </svg>
                            <span className="truncate">{project.name}</span>
                        </Link>
                        {currentPath.startsWith(`/projects/${project.id}`) && (
                            <Link
                                to={`/projects/${project.id}/wiki`}
                                className={`w-full text-left px-3 py-1.5 pl-10 rounded-lg flex items-center gap-2 transition-all duration-200 border-l-4 ${
                                    currentPath === `/projects/${project.id}/wiki`
                                        ? 'bg-accent/5 border-accent text-accent font-medium' 
                                        : 'border-transparent text-text-muted hover:bg-bg-hover hover:text-text-base'
                                }`}
                            >
                                <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                                </svg>
                                <span className="truncate text-sm">Project Wiki</span>
                            </Link>
                        )}
                    </div>
                ))}

                <button
                    onClick={onNewProject}
                    className="w-full text-left px-3 py-2.5 mt-2 rounded-lg flex items-center gap-3 text-text-muted hover:bg-bg-hover hover:text-text-base transition-all duration-200 border border-dashed border-border/60"
                >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                    <span>New Project</span>
                </button>

                <label className="w-full text-left px-3 py-2.5 mt-2 rounded-lg flex items-center gap-3 text-text-muted hover:bg-bg-hover hover:text-text-base transition-all duration-200 border border-dashed border-border/60 cursor-pointer">
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                    </svg>
                    <span>Import Project</span>
                    <input 
                        type="file" 
                        accept=".json" 
                        onChange={onImportProject} 
                        className="hidden" 
                    />
                </label>
                </div>
            </div>

            <div className="p-4 border-t border-border/50">
                <button
                    onClick={onToggleTheme}
                    className="w-full flex items-center justify-between px-3 py-2 text-text-base hover:bg-bg-hover rounded-lg transition-colors"
                >
                    <div className="flex items-center gap-3">
                        {isDark ? (
                            <svg className="w-5 h-5 text-yellow-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                            </svg>
                        ) : (
                            <svg className="w-5 h-5 text-blue-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                            </svg>
                        )}
                        <span>{isDark ? 'Light Mode' : 'Dark Mode'}</span>
                    </div>
                </button>
                <button 
                    onClick={logout}
                    className="w-full mt-2 flex items-center gap-3 px-3 py-2 text-red-500 hover:bg-red-500/10 rounded-lg transition-colors"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 16l4-4m0 0l-4-4m4 4H7m6 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h4a3 3 0 013 3v1" />
                    </svg>
                    <span>Logout</span>
                </button>
            </div>
        </div>
    );
};

export default Sidebar;
