import React from 'react';

const KanbanCard = ({ todo, onDragStart, onTodoUpdate, onTodoDelete, onClick }) => {
    const handleDelete = (e) => {
        e.stopPropagation();
        if (window.confirm('Are you sure you want to delete this todo?')) {
            onTodoDelete(todo.id);
        }
    };

    const handleToggleComplete = (e) => {
        e.stopPropagation();
        onTodoUpdate(todo.id, { ...todo, completed: !todo.completed });
    };

    return (
        <div
            className={`bg-bg-base border rounded-none rounded-tr-lg shadow-sm cursor-pointer hover:-translate-y-1 transition-all duration-200 overflow-hidden group/card ${
                todo.completed 
                    ? 'border-success/40 opacity-75 bg-success/5' 
                    : 'border-border/60 hover:border-primary hover:shadow-md hover:shadow-primary/10'
            }`}
            draggable
            onDragStart={(e) => onDragStart(e, todo)}
            onClick={onClick}
        >
            <div className="p-4">
                <div className="flex justify-between items-start mb-2 gap-2">
                    <span 
                        className={`text-[15px] leading-tight font-semibold transition-colors flex-1 ${
                            todo.completed 
                                ? 'text-success line-through decoration-success/50' 
                                : 'text-text-base group-hover/card:text-primary'
                        }`} 
                        onClick={(e) => {
                            e.stopPropagation();
                            handleToggleComplete(e);
                        }}
                    >
                        {todo.completed && <span className="inline-block mr-2 text-success">✓</span>}
                        {todo.title}
                    </span>
                    <button
                        onClick={(e) => {
                            e.stopPropagation();
                            handleDelete(e);
                        }}
                        className="p-1 -mt-1 -mr-1 text-text-muted opacity-0 group-hover/card:opacity-100 hover:text-danger hover:bg-danger/10 rounded transition-all duration-200 shrink-0"
                        title="Delete"
                    >
                        ✕
                    </button>
                </div>
                
                {todo.content && (
                    <div className={`mt-1.5 text-sm line-clamp-2 ${todo.completed ? 'text-success/70' : 'text-text-muted'}`}>
                        {todo.content}
                    </div>
                )}
                
                <div className="mt-3 pt-3 border-t border-border/40 flex justify-between items-center text-xs font-medium">
                    <span className={`flex items-center gap-1 ${todo.completed ? 'text-success/70' : 'text-text-muted'}`}>
                        <svg className="w-3.5 h-3.5 opacity-70" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        {new Date(todo.created_at).toLocaleDateString(undefined, { month: 'short', day: 'numeric' })}
                    </span>
                    {todo.subtasks?.length > 0 && (
                        <span className={`flex items-center gap-1.5 px-2 py-0.5 rounded-full border ${
                            todo.subtasks.every(st => st.completed) 
                                ? 'bg-success/10 border-success/30 text-success' 
                                : 'bg-accent/10 border-accent/30 text-accent'
                        }`}>
                            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 6h16M4 10h16M4 14h16M4 18h16" />
                            </svg>
                            {todo.subtasks.filter(st => st.completed).length}/{todo.subtasks.length}
                        </span>
                    )}
                </div>
            </div>
        </div>
    );
};

export default KanbanCard;