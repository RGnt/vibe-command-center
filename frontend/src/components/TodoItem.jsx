import React, { useState } from 'react';

const TodoItem = ({ todo, onUpdate, onDelete, onToggle, onAddSubtask }) => {
    const [isEditing, setIsEditing] = useState(false);
    const [editTitle, setEditTitle] = useState(todo.title);
    const [showSubtasks, setShowSubtasks] = useState(true);
    const [newSubtask, setNewSubtask] = useState('');

    const handleUpdate = (e) => {
        e.preventDefault();
        if (editTitle.trim()) {
            onUpdate(todo.id, { ...todo, title: editTitle });
            setIsEditing(false);
        }
    };

    const handleDelete = () => {
        if (window.confirm('Are you sure you want to delete this todo and all its subtasks?')) {
            onDelete(todo.id);
        }
    };

    const handleAddSubtask = (e) => {
        e.preventDefault();
        if (newSubtask.trim()) {
            onAddSubtask(todo.id, newSubtask);
            setNewSubtask('');
        }
    };

    const toggleSubtasks = () => {
        setShowSubtasks(!showSubtasks);
    };

    return (
        <div className={`bg-bg-panel border border-border rounded-lg shadow-sm transition-colors overflow-hidden ${todo.completed ? 'opacity-70' : ''}`}>
            <div className="p-4 flex items-start justify-between gap-4">
                <div className="flex-1 min-w-0">
                    {isEditing ? (
                        <form onSubmit={handleUpdate} className="flex gap-2">
                            <input
                                type="text"
                                value={editTitle}
                                onChange={(e) => setEditTitle(e.target.value)}
                                className="flex-1 p-2 text-sm bg-bg-base border border-primary rounded text-text-base focus:outline-none focus:ring-1 focus:ring-primary"
                                autoFocus
                            />
                            <button type="submit" className="px-3 py-1.5 text-sm font-medium bg-primary text-white rounded hover:bg-primary-hover">Save</button>
                        </form>
                    ) : (
                        <div>
                            <span 
                                onClick={() => onToggle(todo.id)} 
                                className={`text-base font-medium cursor-pointer transition-colors block ${todo.completed ? 'line-through text-text-muted' : 'text-text-base hover:text-primary'}`}
                            >
                                {todo.title}
                            </span>
                            <span className="text-xs text-text-muted block mt-1">
                                {new Date(todo.created_at).toLocaleDateString()}
                            </span>
                        </div>
                    )}
                </div>
                
                <div className="flex shrink-0 items-center gap-2">
                    {!isEditing && (
                        <button
                            onClick={() => setIsEditing(true)}
                            className="px-3 py-1.5 text-sm font-medium text-text-base bg-bg-base border border-border rounded hover:bg-bg-hover hover:border-primary transition-colors"
                        >
                            Edit
                        </button>
                    )}
                    <button
                        onClick={handleDelete}
                        className="px-3 py-1.5 text-sm font-medium text-danger bg-danger/10 border border-transparent rounded hover:bg-danger/20 transition-colors"
                    >
                        Delete
                    </button>
                </div>
            </div>

            {/* Subtasks section */}
            {(todo.subtasks?.length > 0 || !isEditing) && (
                <div className="border-t border-border bg-bg-base/50">
                    <button
                        onClick={toggleSubtasks}
                        className="w-full text-left px-4 py-2 text-sm font-medium text-text-muted hover:text-text-base hover:bg-bg-hover transition-colors flex items-center gap-2"
                    >
                        <span className="w-4 inline-block text-center">{showSubtasks ? '▼' : '▶'}</span> 
                        Subtasks ({todo.subtasks?.length || 0})
                    </button>

                    {showSubtasks && (
                        <div className="p-4 pt-0 space-y-2">
                            {todo.subtasks?.map(subtask => (
                                <div key={subtask.id} className="flex items-center justify-between group py-1">
                                    <div className="flex flex-col">
                                        <span
                                            onClick={() => onToggle(subtask.id)}
                                            className={`text-sm cursor-pointer ${subtask.completed ? 'line-through text-text-muted' : 'text-text-base hover:text-primary'}`}
                                        >
                                            • {subtask.title}
                                        </span>
                                        <span className="text-[10px] text-text-muted ml-3 mt-0.5">
                                            {new Date(subtask.created_at).toLocaleDateString()}
                                        </span>
                                    </div>
                                </div>
                            ))}

                            <form onSubmit={handleAddSubtask} className="flex mt-3 gap-2">
                                <input
                                    type="text"
                                    value={newSubtask}
                                    onChange={(e) => setNewSubtask(e.target.value)}
                                    placeholder="Add a subtask..."
                                    className="flex-1 p-2 text-sm bg-bg-base border border-border rounded text-text-base focus:outline-none focus:border-primary placeholder-text-muted"
                                />
                                <button type="submit" className="px-4 py-2 text-sm font-medium bg-bg-panel hover:bg-primary text-text-base hover:text-white rounded border border-border hover:border-transparent transition-colors">Add Subtask</button>
                            </form>
                        </div>
                    )}
                </div>
            )}
        </div>
    );
};

export default TodoItem;