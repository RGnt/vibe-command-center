import React, { useState, useEffect } from 'react';
import KanbanColumn from './KanbanColumn';
import TaskModal from './TaskModal';

const KanbanBoard = ({ todos, workflows, activeProjectWorkflowId, onTodoUpdate, onTodoDelete, onAddSubtask, onTodoCreate }) => {
    const [stages, setStages] = useState([]);
    const [activeWorkflow, setActiveWorkflow] = useState(null);
    const [showWorkflowModal, setShowWorkflowModal] = useState(false);
    const [newWorkflowName, setNewWorkflowName] = useState('');
    const [newWorkflowStages, setNewWorkflowStages] = useState(['To Do', 'In Progress', 'Review', 'Done']);
    const [draggedItem, setDraggedItem] = useState(null);
    const [draggedOver, setDraggedOver] = useState(null);

    // Task Modal State
    const [isTaskModalOpen, setIsTaskModalOpen] = useState(false);
    const [taskModalMode, setTaskModalMode] = useState('create'); // 'create' or 'edit'
    const [taskModalData, setTaskModalData] = useState(null);
    const [taskModalStage, setTaskModalStage] = useState('');

    // Initialize with workflow
    useEffect(() => {
        if (workflows && workflows.length > 0) {
            let workflowToUse = workflows.find(w => w.id === activeProjectWorkflowId) || workflows[0];
            setActiveWorkflow(workflowToUse);
            setStages(workflowToUse.stages.map(stage => stage.name));
        }
    }, [workflows, activeProjectWorkflowId]);

    // Get todos by stage
    const getTodosByStage = (stageName) => {
        return todos.filter(todo => todo.stage === stageName);
    };

    // Handle drag start
    const handleDragStart = (e, todo) => {
        setDraggedItem(todo);
        e.dataTransfer.effectAllowed = 'move';
    };

    // Handle drag over
    const handleDragOver = (e, stageName) => {
        e.preventDefault();
        setDraggedOver(stageName);
    };

    // Handle drag leave
    const handleDragLeave = () => {
        setDraggedOver(null);
    };

    // Handle drop
    const handleDrop = (e, stageName) => {
        e.preventDefault();
        if (draggedItem && draggedItem.stage !== stageName) {
            // Update the todo's stage
            onTodoUpdate(draggedItem.id, { ...draggedItem, stage: stageName });
        }
        setDraggedItem(null);
        setDraggedOver(null);
    };

    // Create new workflow
    const handleCreateWorkflow = async () => {
        if (!newWorkflowName.trim()) return;

        const workflow = {
            name: newWorkflowName,
            stages: newWorkflowStages.map((name, index) => ({
                id: index + 1,
                name: name,
                order: index + 1
            }))
        };

        try {
            const response = await fetch('http://localhost:8080/api/workflows', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(workflow),
            });

            if (!response.ok) {
                throw new Error('Failed to create workflow');
            }

            const newWorkflow = await response.json();
            setActiveWorkflow(newWorkflow);
            setStages(newWorkflow.stages.map(stage => stage.name));
            setShowWorkflowModal(false);
            setNewWorkflowName('');
            setNewWorkflowStages(['To Do', 'In Progress', 'Review', 'Done']);
        } catch (err) {
            console.error('Error creating workflow:', err);
        }
    };

    // Add a new stage to workflow
    const addStage = () => {
        setNewWorkflowStages([...newWorkflowStages, `Stage ${newWorkflowStages.length + 1}`]);
    };

    // Update a stage name
    const updateStage = (index, name) => {
        const updated = [...newWorkflowStages];
        updated[index] = name;
        setNewWorkflowStages(updated);
    };

    // Remove a stage
    const removeStage = (index) => {
        if (newWorkflowStages.length <= 1) return;
        const updated = [...newWorkflowStages];
        updated.splice(index, 1);
        setNewWorkflowStages(updated);
    };

    // Task Modal Handlers
    const openCreateTaskModal = (stageName) => {
        setTaskModalMode('create');
        setTaskModalStage(stageName);
        setTaskModalData(null);
        setIsTaskModalOpen(true);
    };

    const openEditTaskModal = (todo) => {
        setTaskModalMode('edit');
        setTaskModalStage(todo.stage);
        setTaskModalData(todo);
        setIsTaskModalOpen(true);
    };

    const handleTaskModalSave = (taskData) => {
        if (taskModalMode === 'create') {
            onTodoCreate(taskData);
        } else {
            onTodoUpdate(taskData.id, taskData);
        }
    };

    return (
        <div className="flex flex-col h-full bg-transparent overflow-hidden">
            <div className="flex justify-between items-center pb-4 mb-2">
                <h2 className="text-xl font-semibold text-text-base opacity-0 hidden sm:block">Board</h2>
                <div className="flex gap-4 items-center w-full sm:w-auto justify-between sm:justify-end">
                    <select
                        value={activeWorkflow?.id || ''}
                        onChange={(e) => {
                            const workflow = workflows.find(w => w.id === parseInt(e.target.value));
                            setActiveWorkflow(workflow);
                            setStages(workflow.stages.map(stage => stage.name));
                        }}
                        className="p-2 rounded-lg bg-bg-panel/80 backdrop-blur-md text-text-base border border-border focus:outline-none focus:ring-2 focus:ring-primary shadow-sm"
                    >
                        {workflows.map(workflow => (
                            <option key={workflow.id} value={workflow.id}>
                                {workflow.name}
                            </option>
                        ))}
                    </select>
                    <button
                        onClick={() => setShowWorkflowModal(true)}
                        className="px-4 py-2 bg-primary hover:bg-primary-hover text-white rounded-lg font-medium transition-colors duration-200 shadow-sm"
                    >
                        + New Workflow
                    </button>
                </div>
            </div>

            <div className="flex-1 overflow-x-auto overflow-y-hidden">
                <div className="flex gap-4 h-full min-w-max pb-8 pt-2 px-2">
                    {stages.map((stageName, index) => (
                        <KanbanColumn
                            key={stageName}
                            stageName={stageName}
                            todos={getTodosByStage(stageName)}
                            onDragStart={handleDragStart}
                            onDragOver={handleDragOver}
                            onDragLeave={handleDragLeave}
                            onDrop={handleDrop}
                            draggedOver={draggedOver}
                            onTodoUpdate={onTodoUpdate}
                            onTodoDelete={onTodoDelete}
                            onAddSubtask={onAddSubtask}
                            onTodoCreate={() => openCreateTaskModal(stageName)}
                            onTodoEdit={openEditTaskModal}
                            isLast={index === stages.length - 1}
                        />
                    ))}
                </div>
            </div>

            {/* Task Creation/Editing Modal */}
            <TaskModal
                isOpen={isTaskModalOpen}
                onClose={() => setIsTaskModalOpen(false)}
                onSave={handleTaskModalSave}
                onAddSubtask={onAddSubtask}
                onToggleSubtask={onTodoUpdate}
                onDeleteTodo={onTodoDelete}
                initialData={taskModalData}
                isEditMode={taskModalMode === 'edit'}
                stageName={taskModalStage}
            />

            {/* Workflow Creation Modal */}
            {showWorkflowModal && (
                <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
                    <div className="bg-bg-panel rounded-xl border border-border shadow-2xl w-full max-w-lg overflow-hidden flex flex-col max-h-[90vh]">
                        <div className="p-6 border-b border-border">
                            <h3 className="text-2xl font-bold text-text-base">Create New Workflow</h3>
                        </div>
                        
                        <div className="p-6 overflow-y-auto flex-1">
                            <div className="mb-6">
                                <label className="block text-sm font-medium text-text-muted mb-2">Workflow Name</label>
                                <input
                                    type="text"
                                    value={newWorkflowName}
                                    onChange={(e) => setNewWorkflowName(e.target.value)}
                                    placeholder="e.g. Software Development"
                                    className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary placeholder-text-muted"
                                />
                            </div>

                            <div>
                                <div className="flex justify-between items-center mb-4">
                                    <h4 className="text-lg font-semibold text-text-base">Stages</h4>
                                </div>
                                
                                <div className="space-y-3">
                                    {newWorkflowStages.map((stage, index) => (
                                        <div key={index} className="flex gap-3">
                                            <input
                                                type="text"
                                                value={stage}
                                                onChange={(e) => updateStage(index, e.target.value)}
                                                className="flex-1 p-2 rounded bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary"
                                            />
                                            <button
                                                onClick={() => removeStage(index)}
                                                disabled={newWorkflowStages.length <= 1}
                                                className="p-2 w-10 h-10 flex items-center justify-center rounded bg-danger/10 hover:bg-danger/20 text-danger border border-transparent disabled:opacity-50 disabled:cursor-not-allowed transition-colors"
                                            >
                                                ×
                                            </button>
                                        </div>
                                    ))}
                                </div>
                                
                                <button 
                                    onClick={addStage} 
                                    className="mt-4 w-full py-2 border border-dashed border-border rounded-lg text-text-muted hover:text-primary hover:border-primary transition-colors"
                                >
                                    + Add Stage
                                </button>
                            </div>
                        </div>

                        <div className="p-4 bg-bg-base border-t border-border flex justify-end gap-3">
                            <button
                                onClick={() => setShowWorkflowModal(false)}
                                className="px-5 py-2 rounded font-medium text-text-muted hover:text-text-base hover:bg-bg-panel transition-colors"
                            >
                                Cancel
                            </button>
                            <button
                                onClick={handleCreateWorkflow}
                                className="px-5 py-2 rounded font-medium bg-primary hover:bg-primary-hover text-white transition-colors"
                            >
                                Create Workflow
                            </button>
                        </div>
                    </div>
                </div>
            )}
        </div>
    );
};

export default KanbanBoard;