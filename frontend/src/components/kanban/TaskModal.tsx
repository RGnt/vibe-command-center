import React, { useState, useEffect } from 'react';
import MarkdownRenderer from '../shared/MarkdownRenderer';

const TaskModal = ({ isOpen, onClose, onSave, onAddSubtask, onToggleSubtask, onDeleteTodo, initialData, isEditMode, stageName }) => {
    const [title, setTitle] = useState('');
    const [content, setContent] = useState('');
    const [newSubtask, setNewSubtask] = useState('');
    const [activeTab, setActiveTab] = useState('write'); // 'write' or 'preview'

    useEffect(() => {
        if (isOpen) {
            if (isEditMode && initialData) {
                setTitle(initialData.title || '');
                setContent(initialData.content || '');
            } else {
                setTitle('');
                setContent('');
            }
            setNewSubtask('');
            setActiveTab('write');
        }
    }, [isOpen, initialData, isEditMode]);

    if (!isOpen) return null;

    const handleSave = (e) => {
        e.preventDefault();
        if (title.trim()) {
            if (isEditMode) {
                onSave({ ...initialData, title, content });
            } else {
                onSave({ title, content, stage: stageName });
            }
            onClose();
        }
    };

    const handleAddSubtask = (e) => {
        e.preventDefault();
        if (newSubtask.trim() && initialData?.id) {
            onAddSubtask(initialData.id, newSubtask);
            setNewSubtask('');
        }
    };

    const handleDelete = () => {
        if (window.confirm('Are you sure you want to delete this task?')) {
            onDeleteTodo(initialData.id);
            onClose();
        }
    };

    return (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-bg-panel rounded-xl border border-border shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col max-h-[90vh]">
                <div className="p-6 border-b border-border flex justify-between items-center bg-bg-base">
                    <h3 className="text-2xl font-bold text-text-base">
                        {isEditMode ? 'Edit Task' : 'Create New Task'}
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
                            <label className="block text-sm font-medium text-text-muted mb-2">Title</label>
                            <input
                                type="text"
                                value={title}
                                onChange={(e) => setTitle(e.target.value)}
                                placeholder="Task title..."
                                className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary placeholder-text-muted"
                                autoFocus
                            />
                        </div>

                        <div>
                            <div className="flex justify-between items-center mb-2">
                                <label className="block text-sm font-medium text-text-muted">Description</label>
                                <div className="flex bg-bg-base rounded-lg border border-border overflow-hidden">
                                    <button 
                                        onClick={() => setActiveTab('write')}
                                        className={`px-3 py-1 text-xs font-medium transition-colors ${activeTab === 'write' ? 'bg-primary text-white' : 'text-text-muted hover:bg-bg-hover'}`}
                                    >
                                        Write
                                    </button>
                                    <button 
                                        onClick={() => setActiveTab('preview')}
                                        className={`px-3 py-1 text-xs font-medium transition-colors ${activeTab === 'preview' ? 'bg-primary text-white' : 'text-text-muted hover:bg-bg-hover'}`}
                                    >
                                        Preview
                                    </button>
                                </div>
                            </div>
                            
                            {activeTab === 'write' ? (
                                <textarea
                                    value={content}
                                    onChange={(e) => setContent(e.target.value)}
                                    placeholder="Add more details using Markdown, tables, Mermaid diagrams, or math formulas ($E=mc^2$)..."
                                    className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary placeholder-text-muted min-h-[160px] resize-y font-mono text-sm"
                                />
                            ) : (
                                <div className="w-full p-4 rounded-lg bg-bg-base border border-border min-h-[160px] overflow-auto">
                                    {content.trim() ? (
                                        <MarkdownRenderer content={content} />
                                    ) : (
                                        <div className="text-text-muted text-sm italic">Nothing to preview.</div>
                                    )}
                                </div>
                            )}
                        </div>

                        {isEditMode && initialData && (
                            <div className="pt-4 border-t border-border">
                                <h4 className="text-lg font-semibold text-text-base mb-4">Subtasks</h4>
                                
                                <div className="space-y-2 mb-4">
                                    {initialData.subtasks?.map(subtask => (
                                        <div key={subtask.id} className="flex items-center gap-3 bg-bg-base p-3 rounded-lg border border-border">
                                            <input 
                                                type="checkbox" 
                                                checked={subtask.completed}
                                                onChange={() => onToggleSubtask(subtask.id, { ...subtask, completed: !subtask.completed })}
                                                className="w-4 h-4 text-primary bg-bg-panel border-border rounded focus:ring-primary"
                                            />
                                            <span
                                                className={`flex-1 text-sm ${subtask.completed ? 'line-through text-text-muted' : 'text-text-base'}`}
                                            >
                                                {subtask.title}
                                            </span>
                                        </div>
                                    ))}
                                    {(!initialData.subtasks || initialData.subtasks.length === 0) && (
                                        <div className="text-text-muted text-sm italic">No subtasks added yet.</div>
                                    )}
                                </div>

                                <form onSubmit={handleAddSubtask} className="flex gap-2">
                                    <input
                                        type="text"
                                        value={newSubtask}
                                        onChange={(e) => setNewSubtask(e.target.value)}
                                        placeholder="Add a new subtask..."
                                        className="flex-1 p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary placeholder-text-muted"
                                    />
                                    <button 
                                        type="submit" 
                                        disabled={!newSubtask.trim()}
                                        className="px-4 py-2 font-medium bg-bg-base hover:bg-primary text-text-base hover:text-white rounded-lg border border-border hover:border-transparent transition-colors disabled:opacity-50"
                                    >
                                        Add
                                    </button>
                                </form>
                            </div>
                        )}
                    </div>
                </div>

                <div className="p-4 bg-bg-base border-t border-border flex justify-between items-center">
                    <div>
                        {isEditMode && (
                            <button
                                onClick={handleDelete}
                                className="px-4 py-2 rounded-lg font-medium text-danger hover:bg-danger/10 transition-colors"
                            >
                                Delete Task
                            </button>
                        )}
                    </div>
                    <div className="flex gap-3">
                        <button
                            onClick={onClose}
                            className="px-5 py-2 rounded-lg font-medium text-text-muted hover:text-text-base hover:bg-bg-panel transition-colors"
                        >
                            Cancel
                        </button>
                        <button
                            onClick={handleSave}
                            disabled={!title.trim()}
                            className="px-6 py-2 rounded-lg font-medium bg-primary hover:bg-primary-hover text-white transition-colors disabled:opacity-50"
                        >
                            {isEditMode ? 'Save Changes' : 'Create Task'}
                        </button>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default TaskModal;
