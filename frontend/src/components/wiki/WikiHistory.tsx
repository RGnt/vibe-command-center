import React, { useState, useEffect } from 'react';
import MarkdownRenderer from '../shared/MarkdownRenderer';

const WikiHistory = ({ activePage, onClose }) => {
    const [revisions, setRevisions] = useState([]);
    const [loading, setLoading] = useState(true);
    const [selectedRevision, setSelectedRevision] = useState(null);

    useEffect(() => {
        const fetchRevisions = async () => {
            try {
                const res = await fetch(`/api/wikis/${activePage.id}/revisions`);
                if (!res.ok) throw new Error('Failed to fetch revisions');
                const data = await res.json();
                setRevisions(data || []);
            } catch (err) {
                console.error(err);
            } finally {
                setLoading(false);
            }
        };
        fetchRevisions();
    }, [activePage.id]);

    return (
        <div className="absolute inset-0 bg-bg-base/95 backdrop-blur z-50 flex overflow-hidden">
            <div className="w-80 bg-bg-panel border-r border-border flex flex-col">
                <div className="p-4 border-b border-border flex justify-between items-center">
                    <h3 className="font-bold text-text-base">Revision History</h3>
                    <button onClick={onClose} className="p-1 hover:bg-bg-hover rounded text-text-muted">
                        <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                        </svg>
                    </button>
                </div>
                <div className="flex-1 overflow-y-auto p-4 space-y-2">
                    {loading ? (
                        <div className="text-text-muted text-sm text-center">Loading...</div>
                    ) : revisions.length === 0 ? (
                        <div className="text-text-muted text-sm text-center">No previous revisions.</div>
                    ) : (
                        revisions.map((rev) => (
                            <button
                                key={rev.id}
                                onClick={() => setSelectedRevision(rev)}
                                className={`w-full text-left p-3 rounded-lg border ${
                                    selectedRevision?.id === rev.id 
                                        ? 'border-primary bg-primary/10' 
                                        : 'border-border hover:bg-bg-hover'
                                }`}
                            >
                                <div className="text-sm font-medium text-text-base">
                                    {new Date(rev.created_at).toLocaleString()}
                                </div>
                            </button>
                        ))
                    )}
                </div>
            </div>
            
            <div className="flex-1 flex flex-col p-8 overflow-y-auto">
                <div className="max-w-4xl mx-auto w-full">
                    {selectedRevision ? (
                        <>
                            <div className="mb-6 pb-4 border-b border-border">
                                <h2 className="text-2xl font-bold text-text-base mb-2">Historical Snapshot</h2>
                                <p className="text-text-muted">
                                    Viewing version from {new Date(selectedRevision.created_at).toLocaleString()}
                                </p>
                            </div>
                            <div className="prose prose-invert max-w-none opacity-80">
                                <MarkdownRenderer content={selectedRevision.content || ''} />
                            </div>
                        </>
                    ) : (
                        <div className="h-full flex items-center justify-center text-text-muted">
                            Select a revision from the sidebar to view its content.
                        </div>
                    )}
                </div>
            </div>
        </div>
    );
};

export default WikiHistory;
