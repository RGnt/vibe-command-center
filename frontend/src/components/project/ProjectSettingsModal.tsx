import React, { useState, useEffect } from 'react';

const ProjectSettingsModal = ({ isOpen, onClose, project, onSave }) => {
    const [name, setName] = useState('');
    const [description, setDescription] = useState('');
    const [customFieldSchema, setCustomFieldSchema] = useState([]);

    useEffect(() => {
        if (isOpen && project) {
            setName(project.name || '');
            setDescription(project.description || '');
            
            let parsedSchema = [];
            if (project.custom_field_schema) {
                try {
                    parsedSchema = typeof project.custom_field_schema === 'string' 
                        ? JSON.parse(project.custom_field_schema) 
                        : project.custom_field_schema;
                } catch (e) {
                    console.error('Failed to parse custom field schema', e);
                }
            }
            setCustomFieldSchema(parsedSchema || []);
        }
    }, [isOpen, project]);

    if (!isOpen) return null;

    const handleSave = (e) => {
        e.preventDefault();
        if (name.trim()) {
            onSave({ 
                name, 
                description, 
                custom_field_schema: JSON.stringify(customFieldSchema)
            });
            onClose();
        }
    };

    const addField = () => {
        setCustomFieldSchema([...customFieldSchema, { name: 'New Field', type: 'Text' }]);
    };

    const updateField = (index, field) => {
        const newSchema = [...customFieldSchema];
        newSchema[index] = field;
        setCustomFieldSchema(newSchema);
    };

    const removeField = (index) => {
        const newSchema = [...customFieldSchema];
        newSchema.splice(index, 1);
        setCustomFieldSchema(newSchema);
    };

    return (
        <div className="fixed inset-0 bg-black/60 backdrop-blur-sm flex items-center justify-center z-50 p-4">
            <div className="bg-bg-panel rounded-xl border border-border shadow-2xl w-full max-w-2xl overflow-hidden flex flex-col max-h-[90vh]">
                <div className="p-6 border-b border-border flex justify-between items-center bg-bg-base">
                    <h3 className="text-2xl font-bold text-text-base">Project Settings</h3>
                    <button onClick={onClose} className="text-text-muted hover:text-text-base hover:bg-bg-hover p-2 rounded-lg transition-colors">✕</button>
                </div>
                
                <div className="p-6 overflow-y-auto flex-1 bg-bg-panel space-y-6">
                    <div>
                        <label className="block text-sm font-medium text-text-muted mb-2">Project Name</label>
                        <input
                            type="text"
                            value={name}
                            onChange={(e) => setName(e.target.value)}
                            className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary"
                        />
                    </div>

                    <div>
                        <label className="block text-sm font-medium text-text-muted mb-2">Description</label>
                        <textarea
                            value={description}
                            onChange={(e) => setDescription(e.target.value)}
                            className="w-full p-3 rounded-lg bg-bg-base border border-border text-text-base focus:outline-none focus:ring-2 focus:ring-primary min-h-[100px] resize-y"
                        />
                    </div>

                    <div className="border-t border-border pt-6">
                        <div className="flex justify-between items-center mb-4">
                            <div>
                                <h4 className="text-lg font-semibold text-text-base">Custom Task Fields</h4>
                                <p className="text-sm text-text-muted">Define templates for tasks in this project.</p>
                            </div>
                            <button onClick={addField} className="px-4 py-2 bg-primary/10 text-primary hover:bg-primary/20 rounded-lg font-medium transition-colors text-sm">
                                + Add Field
                            </button>
                        </div>

                        <div className="space-y-4">
                            {customFieldSchema.map((field, index) => (
                                <div key={index} className="flex gap-3 items-start bg-bg-base p-4 rounded-lg border border-border">
                                    <div className="flex-1 space-y-3">
                                        <div className="flex gap-3">
                                            <div className="flex-1">
                                                <label className="block text-xs font-medium text-text-muted mb-1">Field Name</label>
                                                <input
                                                    type="text"
                                                    value={field.name}
                                                    onChange={(e) => updateField(index, { ...field, name: e.target.value })}
                                                    className="w-full p-2 rounded-lg bg-bg-panel border border-border text-text-base text-sm focus:outline-none focus:ring-1 focus:ring-primary"
                                                />
                                            </div>
                                            <div className="w-1/3">
                                                <label className="block text-xs font-medium text-text-muted mb-1">Type</label>
                                                <select
                                                    value={field.type}
                                                    onChange={(e) => updateField(index, { ...field, type: e.target.value })}
                                                    className="w-full p-2 rounded-lg bg-bg-panel border border-border text-text-base text-sm focus:outline-none focus:ring-1 focus:ring-primary"
                                                >
                                                    <option value="Text">Text</option>
                                                    <option value="Number">Number</option>
                                                    <option value="Dropdown">Dropdown</option>
                                                </select>
                                            </div>
                                        </div>
                                        {field.type === 'Dropdown' && (
                                            <div>
                                                <label className="block text-xs font-medium text-text-muted mb-1">Options (comma separated)</label>
                                                <input
                                                    type="text"
                                                    value={(field.options || []).join(', ')}
                                                    onChange={(e) => updateField(index, { ...field, options: e.target.value.split(',').map(s => s.trim()).filter(Boolean) })}
                                                    className="w-full p-2 rounded-lg bg-bg-panel border border-border text-text-base text-sm focus:outline-none focus:ring-1 focus:ring-primary"
                                                    placeholder="Option 1, Option 2, Option 3"
                                                />
                                            </div>
                                        )}
                                    </div>
                                    <button 
                                        onClick={() => removeField(index)}
                                        className="p-2 text-text-muted hover:text-danger hover:bg-danger/10 rounded transition-colors mt-6"
                                    >
                                        ✕
                                    </button>
                                </div>
                            ))}
                            {customFieldSchema.length === 0 && (
                                <div className="text-center p-6 border border-dashed border-border rounded-lg text-text-muted italic text-sm">
                                    No custom fields defined yet. Add fields like "Story Points" or "Environment".
                                </div>
                            )}
                        </div>
                    </div>
                </div>

                <div className="p-4 bg-bg-base border-t border-border flex justify-end gap-3">
                    <button onClick={onClose} className="px-5 py-2 rounded-lg font-medium text-text-muted hover:text-text-base hover:bg-bg-panel transition-colors">
                        Cancel
                    </button>
                    <button onClick={handleSave} disabled={!name.trim()} className="px-6 py-2 rounded-lg font-medium bg-primary hover:bg-primary-hover text-white transition-colors disabled:opacity-50">
                        Save Settings
                    </button>
                </div>
            </div>
        </div>
    );
};

export default ProjectSettingsModal;
