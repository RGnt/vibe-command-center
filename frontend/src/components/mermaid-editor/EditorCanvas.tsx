import React, { useEffect, useRef, useState } from 'react';
import useMermaidStore from '../../store/mermaidStore';
import mermaid from 'mermaid';

mermaid.initialize({
    startOnLoad: false,
    theme: 'dark',
    securityLevel: 'loose',
    fontFamily: 'Inter, sans-serif'
});

const EditorCanvas = () => {
    const { mermaidCode, selectedNodeId, setSelectedNodeId, appendEdge, appendCode, diagramType } = useMermaidStore();
    const containerRef = useRef(null);
    const [error, setError] = useState(null);
    
    // Drag state
    const [handlePos, setHandlePos] = useState(null);
    const [isDragging, setIsDragging] = useState(false);
    const [dragStartPos, setDragStartPos] = useState({ x: 0, y: 0 });
    const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
    const [dragStartId, setDragStartId] = useState(null);
    
    // Menu state
    const [menuPos, setMenuPos] = useState(null);
    const [menuSourceId, setMenuSourceId] = useState(null);

    useEffect(() => {
        let isMounted = true;
        const renderDiagram = async () => {
            if (!containerRef.current) return;
            try {
                containerRef.current.innerHTML = '';
                const { svg } = await mermaid.render(`mermaid-svg-${Date.now()}`, mermaidCode);
                
                if (isMounted) {
                    containerRef.current.innerHTML = svg;
                    setError(null);
                    
                    const svgNodes = containerRef.current.querySelectorAll('.node, .actor, .classGroup, .state');
                    svgNodes.forEach((node) => {
                        node.style.cursor = 'pointer';
                        node.addEventListener('click', (e) => {
                            e.stopPropagation();
                            const id = node.id || node.getAttribute('data-id') || node.getAttribute('id');
                            let cleanId = id;
                            const svgPrefixMatch = cleanId.match(/^mermaid-svg-[^-]+-(.*)$/);
                            if (svgPrefixMatch) cleanId = svgPrefixMatch[1];
                            if (cleanId.startsWith('flowchart-')) cleanId = cleanId.substring(10);
                            cleanId = cleanId.replace(/-\d+$/, '');
                            
                            setSelectedNodeId(cleanId);
                        });
                    });
                }
            } catch (err) {
                if (isMounted) {
                    setError(err.message || 'Syntax Error');
                }
            }
        };

        renderDiagram();
        return () => { isMounted = false; };
    }, [mermaidCode, setSelectedNodeId]);

    // Update handle position
    useEffect(() => {
        if (!containerRef.current || !selectedNodeId) {
            setHandlePos(null);
            return;
        }
        
        // Find the node. Mermaid flowchart nodes might have "-xx" postfix, so we use startsWith or includes
        const nodes = Array.from(containerRef.current.querySelectorAll('.node, .actor, .classGroup, .state'));
        const targetNode = nodes.find(n => n.id === selectedNodeId || n.id.includes(`-${selectedNodeId}-`) || n.id.startsWith(`flowchart-${selectedNodeId}-`));
        
        if (targetNode) {
            const rect = targetNode.getBoundingClientRect();
            setHandlePos({
                x: rect.right + 10,
                y: rect.top + rect.height / 2
            });
            targetNode.classList.add('selected-node');
            targetNode.style.filter = 'drop-shadow(0 0 8px var(--color-primary)) brightness(1.2)';
        }
        
        return () => {
            if (targetNode) {
                targetNode.classList.remove('selected-node');
                targetNode.style.filter = '';
            }
        };
    }, [selectedNodeId, mermaidCode]); // re-run if code changes so handle re-anchors

    // Drag tracking
    useEffect(() => {
        if (!isDragging) return;
        
        const handleMouseMove = (e) => setMousePos({ x: e.clientX, y: e.clientY });
        const handleMouseUp = (e) => {
            setIsDragging(false);
            
            // Check if dropped on a node
            const elements = document.elementsFromPoint(e.clientX, e.clientY);
            let targetNode = null;
            for (const el of elements) {
                const nodeEl = el.closest('.node, .actor, .classGroup, .state');
                if (nodeEl) {
                    targetNode = nodeEl;
                    break;
                }
            }
            
            if (targetNode) {
                const id = targetNode.id || targetNode.getAttribute('data-id') || targetNode.getAttribute('id');
                let cleanId = id;
                const svgPrefixMatch = cleanId.match(/^mermaid-svg-[^-]+-(.*)$/);
                if (svgPrefixMatch) cleanId = svgPrefixMatch[1];
                if (cleanId.startsWith('flowchart-')) cleanId = cleanId.substring(10);
                cleanId = cleanId.replace(/-\d+$/, '');
                
                if (cleanId && cleanId !== dragStartId) {
                    appendEdge(dragStartId, cleanId);
                    return;
                }
            }
            
            // Dropped on empty space
            setMenuPos({ x: e.clientX, y: e.clientY });
            setMenuSourceId(dragStartId);
        };
        
        window.addEventListener('mousemove', handleMouseMove);
        window.addEventListener('mouseup', handleMouseUp);
        return () => {
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseup', handleMouseUp);
        };
    }, [isDragging, dragStartId, appendEdge]);

    const handleMenuSelect = (shape) => {
        const newNodeId = `node_${Date.now()}`;
        let fragment = '';
        
        if (diagramType.startsWith('graph') || diagramType.startsWith('flowchart')) {
            const label = `Node ${newNodeId.toString().substr(-4)}`;
            if (shape === 'rectangle') fragment = `${newNodeId}[${label}]`;
            else if (shape === 'round') fragment = `${newNodeId}(${label})`;
            else if (shape === 'cylinder') fragment = `${newNodeId}[(${label})]`;
        } else if (diagramType === 'sequenceDiagram') {
            if (shape === 'actor') fragment = `actor ${newNodeId}`;
            else if (shape === 'participant') fragment = `participant ${newNodeId}`;
        } else if (diagramType === 'classDiagram') {
            if (shape === 'class') fragment = `class ${newNodeId}`;
        } else if (diagramType === 'stateDiagram-v2') {
            if (shape === 'state') fragment = `state ${newNodeId}`;
        }
        
        if (fragment) {
            appendCode(fragment);
            appendEdge(menuSourceId, newNodeId);
        }
        
        setMenuPos(null);
        setMenuSourceId(null);
    };

    return (
        <div 
            className="w-full h-full bg-bg-base overflow-auto flex items-center justify-center p-8 relative"
            onClick={(e) => {
                if (e.target.closest('.connection-menu')) return;
                setSelectedNodeId(null);
                setMenuPos(null);
            }}
        >
            {error ? (
                <div className="bg-red-900/20 text-red-400 p-4 rounded border border-red-900/50 max-w-lg mt-24">
                    <h3 className="font-bold mb-2">Syntax Error</h3>
                    <pre className="text-xs overflow-auto whitespace-pre-wrap">{error}</pre>
                </div>
            ) : (
                <div 
                    ref={containerRef} 
                    className="mermaid-container w-full h-full flex items-center justify-center min-h-[500px]"
                />
            )}
            
            {/* Connection Handle */}
            {handlePos && !isDragging && !menuPos && (
                <div
                    className="fixed w-4 h-4 bg-primary rounded-full cursor-grab border-2 border-bg-base z-40 connection-handle shadow-md hover:scale-125 transition-transform"
                    style={{ left: handlePos.x, top: handlePos.y, transform: 'translate(0, -50%)' }}
                    onMouseDown={(e) => {
                        e.stopPropagation();
                        setIsDragging(true);
                        setDragStartPos({ x: e.clientX, y: e.clientY });
                        setMousePos({ x: e.clientX, y: e.clientY });
                        setDragStartId(selectedNodeId);
                    }}
                />
            )}
            
            {/* Drag Line */}
            {isDragging && (
                <svg className="fixed inset-0 pointer-events-none z-50" style={{ width: '100vw', height: '100vh' }}>
                    <line 
                        x1={dragStartPos.x} y1={dragStartPos.y} 
                        x2={mousePos.x} y2={mousePos.y} 
                        stroke="var(--color-primary)" 
                        strokeWidth="3" 
                        strokeDasharray="5,5" 
                        className="opacity-70"
                    />
                </svg>
            )}
            
            {/* Context Menu */}
            {menuPos && (
                <div 
                    className="fixed z-50 bg-bg-panel border border-border rounded-lg shadow-xl p-2 flex flex-col gap-1 connection-menu"
                    style={{ left: menuPos.x, top: menuPos.y }}
                >
                    <div className="text-xs text-text-muted px-2 py-1 mb-1 font-medium border-b border-border">Create & Connect</div>
                    {diagramType.startsWith('graph') && (
                        <>
                            <button onClick={() => handleMenuSelect('rectangle')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">Process (Rectangle)</button>
                            <button onClick={() => handleMenuSelect('round')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">Round (Action)</button>
                            <button onClick={() => handleMenuSelect('cylinder')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">Database (Cylinder)</button>
                        </>
                    )}
                    {diagramType === 'sequenceDiagram' && (
                        <>
                            <button onClick={() => handleMenuSelect('actor')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">Actor</button>
                            <button onClick={() => handleMenuSelect('participant')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">Participant</button>
                        </>
                    )}
                    {diagramType === 'classDiagram' && (
                        <button onClick={() => handleMenuSelect('class')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">Class</button>
                    )}
                    {diagramType === 'stateDiagram-v2' && (
                        <button onClick={() => handleMenuSelect('state')} className="text-left px-3 py-1.5 hover:bg-bg-subtle rounded text-sm text-text-base">State</button>
                    )}
                </div>
            )}
        </div>
    );
};

export default EditorCanvas;
