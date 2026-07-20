import React, { useState } from 'react';
import MarkdownRenderer from '../shared/MarkdownRenderer';
import WikiHistory from './WikiHistory';

const WikiViewer = ({ activePage, onEdit, onDelete }) => {
    const [showHistory, setShowHistory] = useState(false);
    if (!activePage) {
        return (
            <div className="flex-1 flex flex-col items-center justify-center text-text-muted p-8 text-center">
                <svg className="w-16 h-16 mb-4 text-border" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={1} d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                </svg>
                <p>Select a page from the sidebar or create a new one.</p>
            </div>
        );
    }

    return (
        <div className="flex-1 overflow-y-auto relative">
            {showHistory && (
                <WikiHistory 
                    activePage={activePage} 
                    onClose={() => setShowHistory(false)} 
                />
            )}
            <div className="max-w-4xl mx-auto px-8 py-12">
                <div className="flex justify-between items-start mb-8 pb-4 border-b border-border/50">
                    <div>
                        <h1 className="text-4xl font-bold text-text-base mb-2">{activePage.title}</h1>
                        <div className="flex gap-3 text-sm text-text-muted">
                            <span className="bg-bg-panel px-2 py-0.5 rounded border border-border">
                                {activePage.category || 'General'}
                            </span>
                            <span>Last updated: {new Date(activePage.updated_at).toLocaleDateString()}</span>
                        </div>
                    </div>
                    <div className="flex gap-2">
                        <button 
                            onClick={() => setShowHistory(true)}
                            className="px-3 py-1.5 rounded-lg text-sm font-medium text-text-muted hover:bg-bg-hover hover:text-primary transition-colors flex items-center gap-2 border border-border"
                            title="View History"
                        >
                            <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                            </svg>
                            History
                        </button>
                        <button 
                            onClick={onEdit}
                            className="p-2 rounded-lg text-text-muted hover:bg-bg-hover hover:text-primary transition-colors"
                            title="Edit Page"
                        >
                            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                            </svg>
                        </button>
                        <button 
                            onClick={onDelete}
                            className="p-2 rounded-lg text-text-muted hover:bg-red-500/10 hover:text-red-500 transition-colors"
                            title="Delete Page"
                        >
                            <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                            </svg>
                        </button>
                    </div>
                </div>

                <div className="prose prose-invert max-w-none">
                    <MarkdownRenderer content={activePage.content || ''} />
                </div>
            </div>
        </div>
    );
};

export default WikiViewer;
