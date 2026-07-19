import React from 'react';

const Dashboard = ({ projects, todos, wikis, fetchProjects, onSelectProject, onDeleteProject }) => {

    const handleDeleteProject = (e, project) => {
        e.stopPropagation(); // prevent triggering onSelectProject
        if (window.confirm(`Are you sure you want to delete project "${project.name}"? This action cannot be undone.`)) {
            onDeleteProject(project.id);
        }
    };

    return (
        <div className="flex-1 overflow-y-auto p-8 bg-bg-base/30">
            <h1 className="text-3xl font-bold text-text-base mb-8">Dashboard</h1>
            
            <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
                {/* Projects Section */}
                <div className="bg-bg-panel border border-border rounded-xl p-6 flex flex-col h-[400px]">
                    <h2 className="text-xl font-semibold text-text-base mb-4 flex items-center justify-between">
                        Projects
                        <span className="bg-primary/20 text-primary text-xs px-2 py-1 rounded-full">{projects.length}</span>
                    </h2>
                    <div className="flex-1 overflow-y-auto space-y-2 pr-2">
                        {projects.length === 0 ? (
                            <p className="text-text-muted text-sm italic">No projects found.</p>
                        ) : (
                            projects.map(p => (
                                <div 
                                    key={p.id} 
                                    className="p-3 bg-bg-subtle border border-border rounded-lg hover:border-primary/50 hover:shadow-md transition-all cursor-pointer flex justify-between items-start group"
                                    onClick={() => onSelectProject(p)}
                                >
                                    <div>
                                        <h3 className="text-text-base font-medium">{p.name}</h3>
                                        {p.description && <p className="text-text-muted text-xs mt-1 line-clamp-1">{p.description}</p>}
                                    </div>
                                    <button 
                                        onClick={(e) => handleDeleteProject(e, p)}
                                        className="text-red-500/50 hover:text-red-500 opacity-0 group-hover:opacity-100 transition-opacity p-1"
                                        title="Delete Project"
                                    >
                                        <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                                        </svg>
                                    </button>
                                </div>
                            ))
                        )}
                    </div>
                </div>

                {/* Open Tasks Section */}
                <div className="bg-bg-panel border border-border rounded-xl p-6 flex flex-col h-[400px]">
                    <h2 className="text-xl font-semibold text-text-base mb-4 flex items-center justify-between">
                        Open Tasks
                        <span className="bg-accent/20 text-accent text-xs px-2 py-1 rounded-full">{todos.length}</span>
                    </h2>
                    <div className="flex-1 overflow-y-auto space-y-2 pr-2">
                        {todos.length === 0 ? (
                            <p className="text-text-muted text-sm italic">You're all caught up!</p>
                        ) : (
                            todos.map(t => (
                                <div key={t.id} className="p-3 bg-bg-subtle border border-border rounded-lg">
                                    <h3 className="text-text-base font-medium text-sm">{t.title}</h3>
                                    <div className="flex items-center gap-2 mt-2">
                                        <span className="text-[10px] uppercase tracking-wider font-bold bg-bg-hover text-text-muted px-2 py-0.5 rounded">
                                            {t.stage}
                                        </span>
                                    </div>
                                </div>
                            ))
                        )}
                    </div>
                </div>

                {/* Latest Wikis Section */}
                <div className="bg-bg-panel border border-border rounded-xl p-6 flex flex-col h-[400px]">
                    <h2 className="text-xl font-semibold text-text-base mb-4 flex items-center justify-between">
                        Latest Wiki Pages
                        <span className="bg-blue-500/20 text-blue-500 text-xs px-2 py-1 rounded-full">{wikis.length}</span>
                    </h2>
                    <div className="flex-1 overflow-y-auto space-y-2 pr-2">
                        {wikis.length === 0 ? (
                            <p className="text-text-muted text-sm italic">No wikis found.</p>
                        ) : (
                            wikis.map(w => (
                                <div key={w.id} className="p-3 bg-bg-subtle border border-border rounded-lg">
                                    <h3 className="text-text-base font-medium text-sm">{w.title}</h3>
                                    <span className="text-xs text-text-muted mt-1 block">{w.category || 'General'}</span>
                                </div>
                            ))
                        )}
                    </div>
                </div>
            </div>
        </div>
    );
};

export default Dashboard;
