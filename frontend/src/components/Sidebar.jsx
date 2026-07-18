import React from 'react';

const Sidebar = ({ projects, activeProject, onSelectProject, onNewProject, onToggleTheme, isDark }) => {
    return (
        <div className="w-64 bg-bg-panel/90 backdrop-blur-xl border-r border-border h-full flex flex-col transition-all duration-300">
            <div className="p-6 border-b border-border/50">
                <h1 className="text-2xl font-bold bg-gradient-to-r from-primary to-accent bg-clip-text text-transparent flex items-center gap-2">
                    <svg className="w-6 h-6 text-primary" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                    KanbanX
                </h1>
            </div>

            <div className="flex-1 overflow-y-auto p-4 space-y-4">
                <div>
                    <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">Global</div>
                    <button
                        onClick={() => onSelectProject('global-wiki')}
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            activeProject === 'global-wiki'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${activeProject === 'global-wiki' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                        </svg>
                        <span>General Wiki</span>
                    </button>
                </div>

                <div>
                    <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">Tools</div>
                    <button
                        onClick={() => onSelectProject('mermaid-editor')}
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            activeProject === 'mermaid-editor'
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${activeProject === 'mermaid-editor' ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                        </svg>
                        <span>Mermaid Editor</span>
                    </button>
                </div>

                <div>
                    <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">Projects</div>
                
                {projects.map(project => (
                    <button
                        key={project.id}
                        onClick={() => onSelectProject(project)}
                        className={`w-full text-left px-3 py-2.5 rounded-lg flex items-center gap-3 transition-all duration-200 border-l-4 ${
                            activeProject?.id === project.id 
                                ? 'bg-accent/10 border-accent text-accent font-semibold shadow-inner' 
                                : 'border-transparent text-text-base hover:bg-bg-hover hover:text-text-base'
                        }`}
                    >
                        <svg className={`w-4 h-4 ${activeProject?.id === project.id ? 'text-accent' : 'text-text-muted'}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                        </svg>
                        <span className="truncate">{project.name}</span>
                    </button>
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
                </div>
            </div>

            <div className="p-4 border-t border-border/50">
                <button
                    onClick={onToggleTheme}
                    className="w-full px-4 py-2 bg-bg-base hover:bg-bg-hover text-text-base rounded-lg border border-border transition-colors duration-200 flex items-center justify-center gap-2"
                >
                    {isDark ? (
                        <>
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707M16 12a4 4 0 11-8 0 4 4 0 018 0z" />
                            </svg>
                            Light Mode
                        </>
                    ) : (
                        <>
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M20.354 15.354A9 9 0 018.646 3.646 9.003 9.003 0 0012 21a9.003 9.003 0 008.354-5.646z" />
                            </svg>
                            Dark Mode
                        </>
                    )}
                </button>
            </div>
        </div>
    );
};

export default Sidebar;
