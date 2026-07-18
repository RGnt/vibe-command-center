import React, { useEffect, useState } from 'react';
import { Panel } from '@xyflow/react';

const EditorSidebar = () => {
    const [icons, setIcons] = useState([]);
    const [expandedFolders, setExpandedFolders] = useState({});

    const fetchIcons = async () => {
        try {
            const res = await fetch('/api/icons');
            if (res.ok) {
                const data = await res.json();
                setIcons(data || []);
            }
        } catch (err) {
            console.error("Failed to fetch icons", err);
        }
    };

    useEffect(() => {
        fetchIcons();
    }, []);

    const handleRenameIcon = async (id, currentName) => {
        const newName = prompt("Enter new name for icon:", currentName);
        if (newName && newName !== currentName) {
            try {
                await fetch(`/api/icons/${id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ name: newName })
                });
                fetchIcons();
            } catch (err) {
                console.error("Failed to rename icon", err);
            }
        }
    };

    const onDragStart = (event, nodeType, shape, url = null) => {
        event.dataTransfer.setData('application/reactflow', nodeType);
        event.dataTransfer.setData('application/reactflow/shape', shape);
        if (url) {
            event.dataTransfer.setData('application/reactflow/url', url);
        }
        event.dataTransfer.effectAllowed = 'move';
    };

    const groupedIcons = icons.reduce((acc, icon) => {
        const folder = icon.folder || 'General';
        if (!acc[folder]) acc[folder] = [];
        acc[folder].push(icon);
        return acc;
    }, {});

    const toggleFolder = (folder) => {
        setExpandedFolders(prev => ({
            ...prev,
            [folder]: prev[folder] === undefined ? false : !prev[folder]
        }));
    };

    const [moveDialog, setMoveDialog] = useState(null);

    const handleMoveIcon = (id, currentFolder) => {
        setMoveDialog({ id, currentFolder });
    };

    const handleConfirmMove = async (newFolder) => {
        if (newFolder && newFolder !== moveDialog.currentFolder) {
            try {
                await fetch(`/api/icons/${moveDialog.id}`, {
                    method: 'PUT',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({ folder: newFolder })
                });
                fetchIcons();
            } catch (err) {
                console.error("Failed to move icon", err);
            }
        }
        setMoveDialog(null);
    };

    return (
        <>
        <aside className="w-64 bg-bg-base border-r border-border p-4 flex flex-col gap-4 overflow-y-auto">
            <h2 className="text-text-base font-semibold text-lg">Palette</h2>
            <p className="text-text-muted text-sm mb-2">Drag shapes onto the canvas.</p>
            
            <div className="flex flex-col gap-3">
                <div 
                    className="border-2 border-primary bg-primary/10 text-primary p-3 rounded cursor-grab text-center font-medium"
                    onDragStart={(event) => onDragStart(event, 'mermaidNode', 'rectangle')}
                    draggable
                >
                    Process (Rectangle)
                </div>
                <div 
                    className="border-2 border-primary bg-primary/10 text-primary p-3 rounded-full cursor-grab text-center font-medium"
                    onDragStart={(event) => onDragStart(event, 'mermaidNode', 'round')}
                    draggable
                >
                    Round (Action)
                </div>
                <div 
                    className="border-2 border-primary bg-primary/10 text-primary p-3 cursor-grab text-center font-medium"
                    style={{ borderRadius: '50% 50% 50% 50% / 15% 15% 15% 15%' }}
                    onDragStart={(event) => onDragStart(event, 'mermaidNode', 'stadium')}
                    draggable
                >
                    Stadium (Terminal)
                </div>
                <div 
                    className="border-2 border-primary bg-primary/10 text-primary p-3 cursor-grab text-center font-medium"
                    style={{ borderRadius: '15px' }}
                    onDragStart={(event) => onDragStart(event, 'mermaidNode', 'cylinder')}
                    draggable
                >
                    Database (Cylinder)
                </div>
                <div 
                    className="border-2 border-primary bg-primary/10 text-primary p-3 rounded-full cursor-grab text-center font-medium aspect-square flex items-center justify-center w-24 self-center"
                    onDragStart={(event) => onDragStart(event, 'mermaidNode', 'circle')}
                    draggable
                >
                    Circle
                </div>
                <div 
                    className="border-2 border-primary bg-primary/10 text-primary p-3 cursor-grab text-center font-medium transform -skew-x-12 w-[80%] self-center"
                    onDragStart={(event) => onDragStart(event, 'mermaidNode', 'rhombus')}
                    draggable
                >
                    Decision (Rhombus)
                </div>
            </div>

            <div className="mt-8">
                <h3 className="text-text-base font-semibold text-md mb-2">Custom Icons</h3>
                {icons.length === 0 ? (
                    <p className="text-text-muted text-sm">No icons uploaded yet.</p>
                ) : (
                    <div className="flex flex-col gap-2">
                        {Object.keys(groupedIcons).sort().map(folder => {
                            const isExpanded = expandedFolders[folder] !== false; // default true
                            return (
                                <div key={folder} className="flex flex-col border border-border rounded overflow-hidden">
                                    <button 
                                        onClick={() => toggleFolder(folder)}
                                        className="bg-bg-subtle p-2 flex items-center justify-between text-sm font-medium hover:bg-bg-base transition-colors"
                                    >
                                        <span>📁 {folder}</span>
                                        <svg className={`w-4 h-4 transform transition-transform ${isExpanded ? 'rotate-180' : ''}`} fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 9l-7 7-7-7" />
                                        </svg>
                                    </button>
                                    {isExpanded && (
                                        <div className="p-2 grid grid-cols-2 gap-2 bg-bg-base">
                                            {groupedIcons[folder].map(icon => (
                                                <div key={icon.id} className="relative group border border-border rounded p-2 bg-bg-subtle flex flex-col items-center justify-center cursor-grab hover:border-primary transition-colors">
                                                    <img 
                                                        src={icon.url} 
                                                        alt={icon.name}
                                                        className="w-10 h-10 object-contain mb-1"
                                                        onDragStart={(event) => onDragStart(event, 'mermaidNode', 'rectangle', icon.url)}
                                                        draggable
                                                    />
                                                    <span className="text-xs text-text-muted truncate w-full text-center" title={icon.name}>
                                                        {icon.name}
                                                    </span>
                                                    <div className="absolute top-1 right-1 opacity-0 group-hover:opacity-100 flex gap-1 transition-opacity">
                                                        <button 
                                                            onClick={() => handleRenameIcon(icon.id, icon.name)}
                                                            className="p-1 bg-bg-base rounded shadow hover:text-primary"
                                                            title="Rename"
                                                        >
                                                            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15.232 5.232l3.536 3.536m-2.036-5.036a2.5 2.5 0 113.536 3.536L6.5 21.036H3v-3.572L16.732 3.732z" />
                                                            </svg>
                                                        </button>
                                                        <button 
                                                            onClick={() => handleMoveIcon(icon.id, icon.folder)}
                                                            className="p-1 bg-bg-base rounded shadow hover:text-primary"
                                                            title="Move to Folder"
                                                        >
                                                            <svg className="w-3 h-3" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                                                                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M3 7v10a2 2 0 002 2h14a2 2 0 002-2V9a2 2 0 00-2-2h-6l-2-2H5a2 2 0 00-2 2z" />
                                                            </svg>
                                                        </button>
                                                    </div>
                                                </div>
                                            ))}
                                        </div>
                                    )}
                                </div>
                            );
                        })}
                    </div>
                )}
            </div>

            <div className="mt-4">
                <input
                    type="file"
                    accept="image/png, image/jpeg, image/svg+xml, image/webp, image/gif"
                    id="image-upload"
                    className="hidden"
                    multiple
                    onChange={async (e) => {
                        const files = Array.from(e.target.files);
                        if (!files.length) return;

                        e.target.value = '';

                        for (const file of files) {
                            const formData = new FormData();
                            formData.append('file', file);

                            try {
                                const res = await fetch('/api/upload', {
                                    method: 'POST',
                                    body: formData
                                });
                                if (!res.ok) throw new Error(`Upload failed for ${file.name}`);
                                const data = await res.json();
                                
                                // Refresh icons list
                                fetchIcons();
                            } catch (err) {
                                console.error(err);
                                alert(`Failed to upload ${file.name}.`);
                            }
                        }
                    }}
                />
                <label 
                    htmlFor="image-upload"
                    className="flex items-center justify-center gap-2 w-full px-4 py-2 bg-primary/20 hover:bg-primary/30 text-primary border border-primary/50 border-dashed rounded-lg cursor-pointer transition-colors"
                >
                    <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12" />
                    </svg>
                    Upload Image
                </label>
            </div>
        </aside>

        {moveDialog && (
            <div className="fixed inset-0 bg-black/50 flex items-center justify-center z-50">
                <div className="bg-bg-base p-6 rounded-lg shadow-xl border border-border w-80">
                    <h3 className="text-lg font-semibold mb-4 text-text-base">Move Icon</h3>
                    <div className="flex flex-col gap-2 mb-6">
                        <label className="text-sm text-text-muted">Select or type folder name:</label>
                        <input 
                            type="text" 
                            list="folder-options"
                            className="w-full bg-bg-subtle border border-border rounded p-2 text-text-base focus:border-primary focus:outline-none"
                            defaultValue={moveDialog.currentFolder}
                            id="folder-input"
                            autoFocus
                        />
                        <datalist id="folder-options">
                            {Object.keys(groupedIcons).map(folder => (
                                <option key={folder} value={folder} />
                            ))}
                        </datalist>
                    </div>
                    <div className="flex justify-end gap-2">
                        <button 
                            onClick={() => setMoveDialog(null)}
                            className="px-4 py-2 rounded text-text-muted hover:bg-bg-subtle transition-colors"
                        >
                            Cancel
                        </button>
                        <button 
                            onClick={() => handleConfirmMove(document.getElementById('folder-input').value)}
                            className="px-4 py-2 bg-primary text-bg-base rounded font-medium hover:opacity-90 transition-opacity"
                        >
                            Move
                        </button>
                    </div>
                </div>
            </div>
        )}
        </>
    );
};

export default EditorSidebar;
