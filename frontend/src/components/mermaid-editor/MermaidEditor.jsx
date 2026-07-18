import React, { useState, useEffect } from 'react';
import EditorSidebar from './EditorSidebar';
import EditorCanvas from './EditorCanvas';
import EditorPreview from './EditorPreview';
import useMermaidStore from '../../store/mermaidStore';

const diagramTypes = [
    { value: 'graph TD', label: 'Top-Down Graph (graph TD)' },
    { value: 'graph LR', label: 'Left-Right Graph (graph LR)' },
    { value: 'flowchart TD', label: 'Top-Down Flowchart (flowchart TD)' },
    { value: 'sequenceDiagram', label: 'Sequence Diagram (sequenceDiagram)' },
    { value: 'classDiagram', label: 'Class Diagram (classDiagram)' },
    { value: 'stateDiagram-v2', label: 'State Diagram (stateDiagram-v2)' },
];

const MermaidEditor = () => {
    const store = useMermaidStore();
    const { diagramId, diagramName, diagramType, mermaidCode, diagramExplanation, setDiagramType, setDiagramName, loadDiagram } = store;
    
    const [savedDiagrams, setSavedDiagrams] = useState([]);

    const fetchDiagrams = async () => {
        try {
            const res = await fetch('/api/diagrams');
            if (res.ok) {
                setSavedDiagrams(await res.json() || []);
            }
        } catch (err) {
            console.error("Failed to fetch diagrams", err);
        }
    };

    useEffect(() => {
        fetchDiagrams();
    }, []);

    const handleSave = async () => {
        let finalName = diagramName;
        if (!diagramId && diagramName === 'Untitled Diagram') {
            const newName = prompt("Enter diagram name:", diagramName);
            if (!newName) return;
            finalName = newName;
            setDiagramName(newName);
        }

        const payload = {
            name: finalName,
            diagram_type: diagramType,
            code: mermaidCode,
            explanation: diagramExplanation
        };

        try {
            const url = diagramId ? `/api/diagrams/${diagramId}` : '/api/diagrams';
            const method = diagramId ? 'PUT' : 'POST';
            
            const res = await fetch(url, {
                method,
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(payload)
            });
            
            if (res.ok) {
                const data = await res.json();
                loadDiagram(data);
                fetchDiagrams();
                alert("Diagram saved successfully!");
            } else {
                alert("Failed to save diagram");
            }
        } catch (err) {
            console.error(err);
            alert("Error saving diagram");
        }
    };

    return (
        <div className="flex h-screen bg-bg-base text-text-base">
            <EditorSidebar />
            <div className="flex-1 flex flex-col">
                <div className="h-14 border-b border-border flex items-center justify-between px-6 bg-bg-subtle">
                    <div className="flex items-center gap-4">
                        <h1 className="text-xl font-bold">
                            <input 
                                value={diagramName} 
                                onChange={(e) => setDiagramName(e.target.value)} 
                                className="bg-transparent border-none focus:outline-none focus:ring-0 font-bold w-48"
                            />
                        </h1>
                        <button 
                            onClick={handleSave}
                            className="bg-primary text-bg-base px-3 py-1 rounded text-sm font-medium hover:bg-primary/90 transition-colors"
                        >
                            Save
                        </button>
                    </div>

                    <div className="flex items-center gap-4">
                        <select 
                            onChange={(e) => {
                                const id = e.target.value;
                                if (!id) return;
                                const d = savedDiagrams.find(x => x.id.toString() === id);
                                if (d) loadDiagram(d);
                                e.target.value = "";
                            }}
                            className="bg-bg-base border border-border text-text-base text-sm rounded-lg focus:ring-primary focus:border-primary block p-2"
                            value=""
                        >
                            <option value="" disabled>Load Diagram...</option>
                            {savedDiagrams.map(d => (
                                <option key={d.id} value={d.id}>{d.name}</option>
                            ))}
                        </select>
                        <select 
                            value={diagramType} 
                            onChange={(e) => setDiagramType(e.target.value)}
                            className="bg-bg-base border border-border text-text-base text-sm rounded-lg focus:ring-primary focus:border-primary block p-2"
                        >
                            {diagramTypes.map(type => (
                                <option key={type.value} value={type.value}>{type.label}</option>
                            ))}
                        </select>
                    </div>
                </div>
                <div className="flex-1 relative">
                    <EditorCanvas />
                </div>
            </div>
            <EditorPreview />
        </div>
    );
};

export default MermaidEditor;
