import React from 'react';

const WikiSidebar = ({ 
    isGlobal, 
    pages, 
    activePage, 
    isEditing, 
    onPageSelect, 
    onCreateNew 
}) => {
    // Group pages by category
    const categories = pages.reduce((acc, page) => {
        const cat = page.category || 'General';
        if (!acc[cat]) acc[cat] = [];
        acc[cat].push(page);
        return acc;
    }, {});

    return (
        <div className="w-64 bg-bg-panel/50 border-r border-border flex flex-col shrink-0">
            <div className="p-4 border-b border-border/50 flex justify-between items-center">
                <h3 className="font-semibold text-text-base">
                    {isGlobal ? 'General Wiki' : 'Project Wiki'}
                </h3>
                <button 
                    onClick={onCreateNew}
                    className="p-1 rounded hover:bg-primary/20 text-primary transition-colors"
                    title="New Page"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 4v16m8-8H4" />
                    </svg>
                </button>
            </div>
            <div className="flex-1 overflow-y-auto p-3 space-y-4">
                {Object.keys(categories).map(cat => (
                    <div key={cat}>
                        <div className="text-xs font-bold text-text-muted uppercase tracking-wider mb-2 px-2">{cat}</div>
                        <div className="space-y-1">
                            {categories[cat].map(page => (
                                <button
                                    key={page.id}
                                    onClick={() => onPageSelect(page)}
                                    className={`w-full text-left px-3 py-1.5 rounded-lg text-sm transition-colors ${
                                        activePage?.id === page.id && !isEditing
                                            ? 'bg-accent/10 text-accent font-medium border-l-2 border-accent'
                                            : 'text-text-base hover:bg-bg-hover border-l-2 border-transparent'
                                    }`}
                                >
                                    {page.title}
                                </button>
                            ))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default WikiSidebar;
