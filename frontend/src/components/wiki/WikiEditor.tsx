import React from 'react';
import Editor from '@monaco-editor/react';

const WikiEditor = ({ 
    activePage, 
    editTitle, 
    setEditTitle, 
    editCategory, 
    setEditCategory, 
    editSlug, 
    setEditSlug, 
    editContent, 
    setEditContent, 
    onCancel, 
    onSave 
}) => {
    return (
        <div className="flex-1 flex flex-col p-6 overflow-hidden">
            <div className="flex justify-between items-center mb-4">
                <h2 className="text-xl font-bold text-text-base">Editing Page</h2>
                <div className="flex gap-2">
                    <button 
                        onClick={onCancel}
                        className="px-4 py-1.5 rounded-lg text-sm font-medium text-text-muted hover:bg-bg-hover transition-colors"
                    >
                        Cancel
                    </button>
                    <button 
                        onClick={onSave}
                        className="px-4 py-1.5 rounded-lg text-sm font-medium bg-primary text-primary-content hover:bg-primary/90 transition-colors"
                    >
                        Save Page
                    </button>
                </div>
            </div>

            <div className="space-y-4 mb-4 shrink-0">
                <input
                    type="text"
                    value={editTitle}
                    onChange={(e) => setEditTitle(e.target.value)}
                    placeholder="Page Title"
                    className="w-full bg-bg-panel border border-border rounded-lg px-4 py-2 text-text-base focus:border-primary focus:outline-none"
                />
                <div className="flex gap-4">
                    <input
                        type="text"
                        value={editCategory}
                        onChange={(e) => setEditCategory(e.target.value)}
                        placeholder="Category (e.g. Engineering)"
                        className="flex-1 bg-bg-panel border border-border rounded-lg px-4 py-2 text-text-base focus:border-primary focus:outline-none"
                    />
                    <input
                        type="text"
                        value={editSlug}
                        onChange={(e) => setEditSlug(e.target.value)}
                        placeholder="URL Slug (e.g. my-page)"
                        className="flex-1 bg-bg-panel border border-border rounded-lg px-4 py-2 text-text-base focus:border-primary focus:outline-none font-mono text-sm"
                        disabled={!!activePage}
                    />
                </div>
            </div>

            <div className="flex-1 rounded-lg overflow-hidden border border-border">
                <Editor
                    height="100%"
                    defaultLanguage="markdown"
                    theme="vs-dark"
                    value={editContent}
                    onChange={(value) => setEditContent(value || '')}
                    options={{
                        minimap: { enabled: false },
                        wordWrap: 'on',
                        padding: { top: 16, bottom: 16 },
                        fontSize: 14,
                        lineNumbers: 'off',
                        folding: false,
                        scrollBeyondLastLine: false,
                    }}
                />
            </div>
        </div>
    );
};

export default WikiEditor;
