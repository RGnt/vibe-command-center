import React, { useState, useEffect } from 'react';
import WikiSidebar from './wiki/WikiSidebar';
import WikiEditor from './wiki/WikiEditor';
import WikiViewer from './wiki/WikiViewer';

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
    const [editParentId, setEditParentId] = useState('');

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
        setEditParentId('');
    };

    const handleEdit = () => {
        if (!activePage) return;
        setIsEditing(true);
        setEditTitle(activePage.title);
        setEditCategory(activePage.category);
        setEditSlug(activePage.slug);
        setEditContent(activePage.content);
        setEditParentId(activePage.parent_id || '');
    };

    const handleSave = async () => {
        const payload = {
            title: editTitle,
            category: editCategory,
            slug: editSlug,
            content: editContent,
            project_id: isGlobal ? null : project.id,
            parent_id: editParentId ? parseInt(editParentId, 10) : null
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

    if (loading && pages.length === 0) {
        return <div className="flex items-center justify-center h-full">Loading Wiki...</div>;
    }

    return (
        <div className="flex h-full w-full bg-bg-base overflow-hidden rounded-xl border border-border shadow-sm">
            <WikiSidebar 
                isGlobal={isGlobal}
                pages={pages}
                activePage={activePage}
                isEditing={isEditing}
                onPageSelect={(page) => {
                    setActivePage(page);
                    setIsEditing(false);
                }}
                onCreateNew={handleCreateNew}
            />

            <div className="flex-1 flex flex-col bg-bg-base/30 relative">
                {isEditing ? (
                    <WikiEditor 
                        activePage={activePage}
                        editTitle={editTitle}
                        setEditTitle={setEditTitle}
                        editCategory={editCategory}
                        setEditCategory={setEditCategory}
                        editSlug={editSlug}
                        setEditSlug={setEditSlug}
                        editContent={editContent}
                        setEditContent={setEditContent}
                        editParentId={editParentId}
                        setEditParentId={setEditParentId}
                        pages={pages}
                        onCancel={() => {
                            setIsEditing(false);
                            if (!activePage) fetchPages();
                        }}
                        onSave={handleSave}
                    />
                ) : (
                    <WikiViewer 
                        activePage={activePage}
                        onEdit={handleEdit}
                        onDelete={handleDelete}
                    />
                )}
            </div>
        </div>
    );
};

export default WikiView;
