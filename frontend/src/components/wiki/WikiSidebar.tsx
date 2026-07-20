import React from 'react';

const WikiSidebar = ({ 
    isGlobal, 
    pages, 
    activePage, 
    isEditing, 
    onPageSelect, 
    onCreateNew 
}) => {
    // Group top-level pages by category
    const topLevelPages = pages.filter(p => !p.parent_id);
    const categories = topLevelPages.reduce((acc, page) => {
        const cat = page.category || 'General';
        if (!acc[cat]) acc[cat] = [];
        acc[cat].push(page);
        return acc;
    }, {});

    const renderPageNode = (page, depth = 0) => {
        const children = pages.filter(p => p.parent_id === page.id);
        const isSelected = activePage?.id === page.id && !isEditing;
        
        return (
            <div key={page.id} className="w-full">
                <button
                    onClick={() => onPageSelect(page)}
                    className={`w-full text-left py-1.5 rounded-lg text-sm transition-colors flex items-center ${
                        isSelected
                            ? 'bg-accent/10 text-accent font-medium border-l-2 border-accent'
                            : 'text-text-base hover:bg-bg-hover border-l-2 border-transparent'
                    }`}
                    style={{ paddingLeft: `${(depth * 1.5) + 0.75}rem`, paddingRight: '0.75rem' }}
                >
                    {children.length > 0 && (
                        <svg className="w-3 h-3 mr-1.5 text-text-muted shrink-0" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                        </svg>
                    )}
                    {children.length === 0 && <span className="w-4.5 mr-1.5 inline-block shrink-0"></span>}
                    <span className="truncate">{page.title}</span>
                </button>
                {children.length > 0 && (
                    <div className="flex flex-col mt-0.5">
                        {children.map(child => renderPageNode(child, depth + 1))}
                    </div>
                )}
            </div>
        );
    };

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
                            {categories[cat].map(page => renderPageNode(page, 0))}
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
};

export default WikiSidebar;
