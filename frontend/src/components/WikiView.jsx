import React, { useState, useEffect } from 'react';
import MarkdownRenderer from './MarkdownRenderer';

const WikiView = ({ project, isGlobal = false }) => {
    const [pages, setPages] = useState([]);
    const [activePage, setActivePage] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);
    
    // Edit mode state
    const [isEditing, setIsEditing] = useState(false);
    const [editTitle, setEditTitle] = useState('');
    const [editCategory, setEditCategory] = useState('');
    const [editSlug, setEditSlug] = useState('');
    const [editContent, setEditContent] = useState('');

    useEffect(() => {
        fetchPages();
        
        // Listen for internal wiki navigation events
        const handleWikiNav = (e) => {
            const slug = e.detail.slug;
            setPages(currentPages => {
                const targetPage = currentPages.find(p => p.slug === slug);
                if (targetPage) {
                    setActivePage(targetPage);
                    setIsEditing(false);
                }
                return currentPages;
            });
        };
        
        window.addEventListener('wikiNavigate', handleWikiNav);
        return () => window.removeEventListener('wikiNavigate', handleWikiNav);
    }, [project, isGlobal]);

    const fetchPages = async () => {
        setLoading(true);
        try {
            const url = isGlobal 
                ? '/api/wikis' 
                : `/api/wikis?project_id=${project.id}`;
            const response = await fetch(url);
            if (!response.ok) throw new Error('Failed to fetch wiki pages');
            
            const data = await response.json();
            setPages(data || []);
            
            // Set default active page if none selected
            if (data && data.length > 0 && !activePage) {
                setActivePage(data[0]);
            }
        } catch (err) {
            setError(err.message);
        } finally {
            setLoading(false);
        }
    };

    const handleCreateNew = () => {
        setIsEditing(true);
        setActivePage(null);
        setEditTitle('New Page');
        setEditCategory('General');
        setEditSlug(`new-page-${Date.now()}`);
        setEditContent('# New Page\n\nWrite something here...');
    };

    const handleEdit = () => {
        if (!activePage) return;
        setIsEditing(true);
        setEditTitle(activePage.title);
        setEditCategory(activePage.category);
        setEditSlug(activePage.slug);
        setEditContent(activePage.content);
    };

    const handleSave = async () => {
        const payload = {
            title: editTitle,
            category: editCategory,
            slug: editSlug,
            content: editContent,
            project_id: isGlobal ? null : project.id
        };

        try {
            const url = activePage 
                ? `/api/wikis/${activePage.id}`
                : '/api/wikis';
            const method = activePage ? 'PUT' : 'POST';

            const response = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });

            if (!response.ok) throw new Error('Failed to save wiki page');
            const savedPage = await response.json();
            
            await fetchPages();
            setActivePage(savedPage);
            setIsEditing(false);
        } catch (err) {
            alert(err.message);
        }
    };

    const handleDelete = async () => {
        if (!activePage) return;
        if (!window.confirm('Are you sure you want to delete this wiki page?')) return;

        try {
            const response = await fetch(`/api/wikis/${activePage.id}`, {
                method: 'DELETE'
            });

            if (!response.ok) throw new Error('Failed to delete wiki page');
            
            setActivePage(null);
            await fetchPages();
        } catch (err) {
            alert(err.message);
        }
    };

    // Group pages by category
    const categories = pages.reduce((acc, page) => {
        const cat = page.category || 'General';
        if (!acc[cat]) acc[cat] = [];
        acc[cat].push(page);
        return acc;
    }, {});

    if (loading && pages.length === 0) {
        return <div className="flex items-center justify-center h-full">Loading Wiki...</div>;
    }

    return (
        <div className="flex h-full w-full bg-bg-base overflow-hidden rounded-xl border border-border shadow-sm">
            {/* Wiki Sidebar */}
            <div className="w-64 bg-bg-panel/50 border-r border-border flex flex-col shrink-0">
                <div className="p-4 border-b border-border/50 flex justify-between items-center">
                    <h3 className="font-semibold text-text-base">
                        {isGlobal ? 'General Wiki' : 'Project Wiki'}
                    </h3>
                    <button 
                        onClick={handleCreateNew}
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
                                        onClick={() => {
                                            setActivePage(page);
                                            setIsEditing(false);
                                        }}
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

            {/* Wiki Content Area */}
            <div className="flex-1 flex flex-col bg-bg-base/30 relative">
                {isEditing ? (
                    <div className="flex-1 flex flex-col p-6 overflow-hidden">
                        <div className="flex justify-between items-center mb-4">
                            <h2 className="text-xl font-bold text-text-base">Editing Page</h2>
                            <div className="flex gap-2">
                                <button 
                                    onClick={() => {
                                        setIsEditing(false);
                                        if (!activePage) fetchPages();
                                    }}
                                    className="px-4 py-1.5 rounded-lg text-sm font-medium text-text-muted hover:bg-bg-hover transition-colors"
                                >
                                    Cancel
                                </button>
                                <button 
                                    onClick={handleSave}
                                    className="px-4 py-1.5 rounded-lg text-sm font-medium bg-primary text-white hover:bg-primary-hover transition-colors shadow-sm"
                                >
                                    Save
                                </button>
                            </div>
                        </div>
                        <div className="grid grid-cols-2 gap-4 mb-4">
                            <div>
                                <label className="block text-xs font-medium text-text-muted mb-1">Title</label>
                                <input 
                                    type="text" 
                                    value={editTitle}
                                    onChange={e => setEditTitle(e.target.value)}
                                    className="w-full p-2 rounded-lg bg-bg-panel border border-border text-sm focus:ring-1 focus:ring-primary focus:outline-none"
                                />
                            </div>
                            <div>
                                <label className="block text-xs font-medium text-text-muted mb-1">Category</label>
                                <input 
                                    type="text" 
                                    value={editCategory}
                                    onChange={e => setEditCategory(e.target.value)}
                                    className="w-full p-2 rounded-lg bg-bg-panel border border-border text-sm focus:ring-1 focus:ring-primary focus:outline-none"
                                />
                            </div>
                            <div className="col-span-2">
                                <label className="block text-xs font-medium text-text-muted mb-1">Slug (for internal linking)</label>
                                <input 
                                    type="text" 
                                    value={editSlug}
                                    onChange={e => setEditSlug(e.target.value)}
                                    className="w-full p-2 rounded-lg bg-bg-panel border border-border text-sm focus:ring-1 focus:ring-primary focus:outline-none font-mono"
                                />
                            </div>
                        </div>
                        <div className="flex-1 flex flex-col min-h-[300px]">
                            <label className="block text-xs font-medium text-text-muted mb-1">Content (Markdown)</label>
                            <textarea 
                                value={editContent}
                                onChange={e => setEditContent(e.target.value)}
                                className="flex-1 w-full p-4 rounded-lg bg-bg-panel border border-border text-sm focus:ring-1 focus:ring-primary focus:outline-none font-mono resize-none"
                            />
                        </div>
                    </div>
                ) : activePage ? (
                    <div className="flex-1 overflow-y-auto p-8">
                        <div className="flex justify-between items-start mb-6 pb-4 border-b border-border/50">
                            <div>
                                <h1 className="text-3xl font-bold text-text-base mb-2">{activePage.title}</h1>
                                <div className="text-xs text-text-muted">
                                    Last updated: {new Date(activePage.updated_at).toLocaleString()}
                                </div>
                            </div>
                            <div className="flex gap-2">
                                <button 
                                    onClick={handleEdit}
                                    className="px-3 py-1.5 text-sm font-medium rounded-lg text-primary bg-primary/10 hover:bg-primary/20 transition-colors"
                                >
                                    Edit
                                </button>
                                <button 
                                    onClick={handleDelete}
                                    className="px-3 py-1.5 text-sm font-medium rounded-lg text-danger bg-danger/10 hover:bg-danger/20 transition-colors"
                                >
                                    Delete
                                </button>
                            </div>
                        </div>
                        <div className="bg-bg-panel/30 p-6 rounded-xl border border-border/30">
                            <MarkdownRenderer content={activePage.content} />
                        </div>
                    </div>
                ) : (
                    <div className="flex-1 flex items-center justify-center text-text-muted">
                        Select a page or create a new one.
                    </div>
                )}
            </div>
        </div>
    );
};

export default WikiView;
