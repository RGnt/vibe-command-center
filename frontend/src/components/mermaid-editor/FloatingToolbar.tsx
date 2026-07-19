import React, { useState } from 'react';
import useMermaidStore from '../../store/mermaidStore';

const FloatingToolbar = () => {
    const { diagramType, appendCode } = useMermaidStore();
    const [isOpen, setIsOpen] = useState(false);

    let id = Date.now();
    const getId = () => `node_${id++}`;

    const handleAddNode = (nodeType, shape) => {
        const newNodeId = getId();
        let fragment = '';

        // Construct standard string based on shape
        if (diagramType.startsWith('graph') || diagramType.startsWith('flowchart')) {
            const label = `Node ${newNodeId.toString().substr(-4)}`;
            if (shape === 'rectangle') fragment = `${newNodeId}[${label}]`;
            else if (shape === 'round') fragment = `${newNodeId}(${label})`;
            else if (shape === 'stadium') fragment = `${newNodeId}([${label}])`;
            else if (shape === 'cylinder') fragment = `${newNodeId}[(${label})]`;
            else if (shape === 'circle') fragment = `${newNodeId}((${label}))`;
            else if (shape === 'rhombus') fragment = `${newNodeId}{${label}}`;
        } else if (diagramType === 'sequenceDiagram') {
            const label = `Node ${newNodeId.toString().substr(-4)}`;
            if (shape === 'actor') fragment = `actor ${newNodeId}`;
            else if (shape === 'participant') fragment = `participant ${newNodeId}`;
            else if (shape === 'note') fragment = `Note over ${newNodeId}: A note`;
        } else if (diagramType === 'classDiagram') {
            if (shape === 'class') fragment = `class ${newNodeId}`;
            else if (shape === 'interface') fragment = `class ${newNodeId} {\n  <<interface>>\n}`;
        } else if (diagramType === 'stateDiagram-v2') {
            if (shape === 'initial') fragment = `[*] --> ${newNodeId}`;
            else if (shape === 'state') fragment = `state ${newNodeId}`;
            else if (shape === 'final') fragment = `${newNodeId} --> [*]`;
        }

        if (fragment) {
            appendCode(fragment);
        }

        setIsOpen(false);
    };

    const renderFlowchartPalette = () => (
        <div className="grid grid-cols-2 gap-2 p-2 w-64">
            <button onClick={() => handleAddNode('mermaidNode', 'rectangle')} className="border border-border hover:border-primary p-2 rounded text-sm text-center text-text-base">Process (Rectangle)</button>
            <button onClick={() => handleAddNode('mermaidNode', 'round')} className="border border-border hover:border-primary p-2 rounded-full text-sm text-center text-text-base">Round (Action)</button>
            <button onClick={() => handleAddNode('mermaidNode', 'stadium')} className="border border-border hover:border-primary p-2 rounded-full text-sm text-center text-text-base">Stadium (Terminal)</button>
            <button onClick={() => handleAddNode('mermaidNode', 'cylinder')} className="border border-border hover:border-primary p-2 rounded text-sm text-center text-text-base">Database (Cylinder)</button>
            <button onClick={() => handleAddNode('mermaidNode', 'circle')} className="border border-border hover:border-primary p-2 rounded-full aspect-square text-sm flex items-center justify-center text-text-base">Circle</button>
            <button onClick={() => handleAddNode('mermaidNode', 'rhombus')} className="border border-border hover:border-primary p-2 transform -skew-x-12 text-sm text-center text-text-base">Decision (Rhombus)</button>
        </div>
    );

    const renderSequencePalette = () => (
        <div className="grid grid-cols-1 gap-2 p-2 w-48">
            <button onClick={() => handleAddNode('sequenceNode', 'actor')} className="border border-border hover:border-green-500 p-2 rounded text-sm text-center text-text-base">Actor</button>
            <button onClick={() => handleAddNode('sequenceNode', 'participant')} className="border border-border hover:border-blue-500 p-2 rounded text-sm text-center text-text-base">Participant</button>
            <button onClick={() => handleAddNode('sequenceNode', 'note')} className="border border-border hover:border-yellow-500 p-2 rounded text-sm text-center text-text-base">Note</button>
        </div>
    );

    const renderClassPalette = () => (
        <div className="grid grid-cols-1 gap-2 p-2 w-48">
            <button onClick={() => handleAddNode('classNode', 'class')} className="border border-border hover:border-indigo-500 p-2 rounded text-sm text-center text-text-base">Class</button>
            <button onClick={() => handleAddNode('classNode', 'interface')} className="border border-border hover:border-purple-500 p-2 rounded text-sm text-center text-text-base">Interface</button>
        </div>
    );

    const renderStatePalette = () => (
        <div className="grid grid-cols-1 gap-2 p-2 w-48">
            <button onClick={() => handleAddNode('stateNode', 'initial')} className="border border-border hover:border-gray-500 p-2 rounded-full w-10 h-10 self-center mx-auto bg-gray-800"></button>
            <button onClick={() => handleAddNode('stateNode', 'state')} className="border border-border hover:border-red-500 p-2 rounded-full text-sm text-center text-text-base">State</button>
            <button onClick={() => handleAddNode('stateNode', 'final')} className="border-4 border-gray-500 hover:border-gray-400 p-2 rounded-full w-10 h-10 self-center mx-auto bg-gray-800"></button>
        </div>
    );

    const renderDynamicPalette = () => {
        if (diagramType.startsWith('graph') || diagramType.startsWith('flowchart')) return renderFlowchartPalette();
        if (diagramType === 'sequenceDiagram') return renderSequencePalette();
        if (diagramType === 'classDiagram') return renderClassPalette();
        if (diagramType === 'stateDiagram-v2') return renderStatePalette();
        return renderFlowchartPalette();
    };

    return (
        <div className="absolute top-4 left-1/2 -translate-x-1/2 z-50 flex items-start flex-col gap-2">
            <div className="bg-bg-panel border border-border rounded-lg shadow-lg p-1 flex items-center justify-center gap-2 w-full">
                <button
                    onClick={() => setIsOpen(!isOpen)}
                    className="flex items-center gap-2 bg-primary text-bg-base px-4 py-2 rounded font-medium hover:bg-primary/90 transition-colors"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><path d="M5 12h14" /><path d="M12 5v14" /></svg>
                    Add Shape
                </button>
            </div>

            {isOpen && (
                <div className="bg-bg-panel border border-border rounded-lg shadow-xl overflow-hidden self-center">
                    {renderDynamicPalette()}
                </div>
            )}
        </div>
    );
};

export default FloatingToolbar;
