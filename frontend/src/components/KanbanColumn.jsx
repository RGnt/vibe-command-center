import React from 'react';
import KanbanCard from './KanbanCard';

const KanbanColumn = ({
    stageName,
    todos,
    onDragStart,
    onDragOver,
    onDragLeave,
    onDrop,
    draggedOver,
    onTodoUpdate,
    onTodoDelete,
    onAddSubtask,
    onTodoCreate,
    onTodoEdit,
    isLast
}) => {
    // Map stages to Tokyo Night colors to make them pop
    const getStageColor = (name) => {
        const n = name.toLowerCase();
        if (n.includes('done') || n.includes('complete')) return 'text-success bg-success/10 border-success/30';
        if (n.includes('review') || n.includes('test')) return 'text-accent bg-accent/10 border-accent/30';
        if (n.includes('progress') || n.includes('doing')) return 'text-primary bg-primary/10 border-primary/30';
        return 'text-warning bg-warning/10 border-warning/30'; // Default / To Do
    };

    const colorClasses = getStageColor(stageName);

    return (
        <div
            className={`flex flex-col w-[340px] shrink-0 rounded-none rounded-tl-xl bg-bg-panel/50 backdrop-blur-md border shadow-sm transition-all duration-300 ${
                draggedOver === stageName 
                    ? `ring-2 ring-primary bg-bg-hover/60 scale-[1.01] border-primary/50` 
                    : 'border-border/50 hover:border-border'
            }`}
            onDragOver={(e) => onDragOver(e, stageName)}
            onDragLeave={onDragLeave}
            onDrop={(e) => onDrop(e, stageName)}
        >
            <div className={`flex justify-between items-center p-4 border-b border-border/50 bg-bg-base/30 rounded-none rounded-tl-xl`}>
                <h3 className="font-semibold text-text-base text-lg tracking-wide flex items-center gap-2">
                    <div className={`w-2 h-2 rounded-full ${colorClasses.split(' ')[1].replace('/10', '')}`}></div>
                    {stageName}
                </h3>
                <span className={`px-2.5 py-0.5 rounded-full text-xs font-bold border ${colorClasses}`}>
                    {todos.length}
                </span>
            </div>

            <div className="p-3 flex-1 overflow-y-auto space-y-3">
                {todos.map(todo => (
                    <KanbanCard
                        key={todo.id}
                        todo={todo}
                        onDragStart={onDragStart}
                        onTodoUpdate={onTodoUpdate}
                        onTodoDelete={onTodoDelete}
                        onAddSubtask={onAddSubtask}
                        onClick={() => onTodoEdit(todo)}
                    />
                ))}

                <button
                    onClick={onTodoCreate}
                    className="w-full mt-3 py-2 text-text-muted hover:text-text-base hover:bg-bg-hover rounded-lg transition-colors flex items-center justify-center gap-2 border border-dashed border-border hover:border-text-muted"
                >
                    <span>+ Add Task</span>
                </button>
            </div>
        </div>
    );
};

export default KanbanColumn;