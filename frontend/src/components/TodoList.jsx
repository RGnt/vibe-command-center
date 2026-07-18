import React from 'react';
import TodoItem from './TodoItem';

const TodoList = ({ todos, onUpdate, onDelete, onToggle, onAddSubtask }) => {
    if (todos.length === 0) {
        return <div className="text-center p-8 text-text-muted bg-bg-panel rounded-xl border border-dashed border-border mt-4">No todos yet. Add one above!</div>;
    }

    return (
        <div className="flex flex-col gap-3 mt-4">
            {todos.map(todo => (
                <TodoItem
                    key={todo.id}
                    todo={todo}
                    onUpdate={onUpdate}
                    onDelete={onDelete}
                    onToggle={onToggle}
                    onAddSubtask={onAddSubtask}
                />
            ))}
        </div>
    );
};

export default TodoList;