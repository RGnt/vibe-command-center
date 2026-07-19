import React, { useState, useEffect } from 'react';

const ProjectModal = ({ isOpen, onClose, onSave, workflows }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [workflowId, setWorkflowId] = useState('');

    useEffect(() => {
        if (isOpen) {
            setName('');
            setDescription('');
            if (workflows && workflows.length > 0) {
                setWorkflowId(workflows[0].id.toString());
            }
        }
    }, [isOpen, workflows]);

    if (!isOpen) return null;

    const handleSave = (e) => {
        e.preventDefault();
        if (name.trim()) {
            onSave({ 
                name, 
                description, 
                workflow_id: workflowId ? parseInt(workflowId) : null 
            });
            onClose();
        }
    };

    return (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-bg-panel rounded-xl border border-border shadow-2xl w-full max-w-lg overflow-hidden flex flex-col">
                <div className="p-6 border-b border-border flex justify-between items-center bg-bg-base">
                    <h3 className="text-2xl font-bold text-text-base">
                        Create New Project
                    </h3>
                    <button 
                        onClick={onClose}
                        className="text-text-muted hover:text-text-base hover:bg-bg-hover p-2 rounded-lg transition-colors"
                    >
                        ✕
                    </button>
                </div>
                
                <div className="p-6 overflow-y-auto flex-1 bg-bg-panel">
                    <div className="space-y-6">
                        <div>
                            <label className="block text-sm font-medium text-text-muted mb-2">Project Name</label>
                            <input
                                type="text"
                                value={name}
                                onChange={(e) => setName(e.target.value)}
                                placeholder="E.g., Website Redesign"
                                className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary placeholder-text-muted"
                                autoFocus
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-text-muted mb-2">Description</label>
                            <textarea
                                value={description}
                                onChange={(e) => setDescription(e.target.value)}
                                placeholder="What is this project about?"
                                className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary placeholder-text-muted min-h-[100px] resize-y"
                            />
                        </div>

                        <div>
                            <label className="block text-sm font-medium text-text-muted mb-2">Default Workflow</label>
                            <select
                                value={workflowId}
                                onChange={(e) => setWorkflowId(e.target.value)}
                                className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary"
                            >
                                <option value="">Select a workflow</option>
                                {workflows?.map(w => (
                                    <option key={w.id} value={w.id}>{w.name}</option>
                                ))}
                            </select>
                        </div>
                    </div>
                </div>

                <div className="p-4 bg-bg-base border-t border-border flex justify-end items-center gap-3">
                    <button
                        onClick={onClose}
                        className="px-5 py-2 rounded-lg font-medium text-text-muted hover:text-text-base hover:bg-bg-panel transition-colors"
                    >
                        Cancel
                    </button>
                    <button
                        onClick={handleSave}
                        disabled={!name.trim()}
                        className="px-6 py-2 rounded-lg font-medium bg-primary hover:bg-primary-hover text-white transition-colors disabled:opacity-50"
                    >
                        Create Project
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ProjectModal;
