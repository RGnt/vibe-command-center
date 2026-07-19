import React, { useEffect, useRef, useState, useCallback } from 'react';
import useMermaidStore from '../../store/mermaidStore';
import mermaid from 'mermaid';

mermaid.initialize({
    startOnLoad: false,
    theme: 'dark',
    securityLevel: 'loose',
    fontFamily: 'Inter, sans-serif'
});

interface Transform {
    scale: number;
    x: number;
    y: number;
}

const SCALE_MIN = 0.2;
const SCALE_MAX = 5.0;
const SCALE_STEP = 0.15;

const EditorCanvas = () => {
    const { mermaidCode, selectedNodeId, setSelectedNodeId, appendEdge } = useMermaidStore();
    const containerRef = useRef<HTMLDivElement>(null);
    const wrapperRef = useRef<HTMLDivElement>(null);
    const [error, setError] = useState<string | null>(null);

    // Zoom / pan transform
    const [transform, setTransform] = useState<Transform>({ scale: 1, x: 0, y: 0 });

    // Pan state (managed via refs to avoid stale closures in event listeners)
    const isPanning = useRef(false);
    const panStart = useRef({ x: 0, y: 0 });
    const transformRef = useRef(transform);
    transformRef.current = transform;

    // Drag-to-connect state
    const [handlePos, setHandlePos] = useState<{ x: number; y: number } | null>(null);
    const [isDragging, setIsDragging] = useState(false);
    const [dragStartPos, setDragStartPos] = useState({ x: 0, y: 0 });
    const [mousePos, setMousePos] = useState({ x: 0, y: 0 });
    const [dragStartId, setDragStartId] = useState<string | null>(null);

    // ── Render Mermaid ──────────────────────────────────────────────────────────
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
                    svgNodes.forEach((node, index) => {
                        // Resolve a stable ID for this node: element id → parent g[id] → inner text → index
                        let rawId = node.id || node.getAttribute('data-id') || '';
                        if (!rawId) {
                            const pg = (node as Element).closest('g[id]') as SVGGElement | null;
                            rawId = pg?.id || '';
                        }
                        if (!rawId) {
                            const textEl = node.querySelector('text');
                            rawId = textEl?.textContent?.trim() || '';
                        }
                        // Strip mermaid-generated prefixes
                        let cleanId = rawId;
                        const svgPrefixMatch = cleanId.match(/^mermaid-svg-[^-]+-(.*)$/);
                        if (svgPrefixMatch) cleanId = svgPrefixMatch[1];
                        if (cleanId.startsWith('flowchart-')) cleanId = cleanId.substring(10);
                        cleanId = cleanId.replace(/-\d+$/, '');
                        const finalId = cleanId || `seq-node-${index}`;

                        // Stamp the resolved id so the handle-finder can look it up reliably
                        (node as HTMLElement).dataset.resolvedNodeId = finalId;
                        (node as HTMLElement).style.cursor = 'pointer';

                        node.addEventListener('click', (e) => {
                            e.stopPropagation();
                            setSelectedNodeId(finalId);
                        });
                    });
                }
            } catch (err: unknown) {
                if (isMounted) {
                    setError((err as Error).message || 'Syntax Error');
                }
            }
        };

        renderDiagram();
        return () => { isMounted = false; };
    }, [mermaidCode, setSelectedNodeId]);

    // ── Connection handle position ──────────────────────────────────────────────
    useEffect(() => {
        if (!containerRef.current || !selectedNodeId) {
            setHandlePos(null);
            return;
        }

        // Find by the data-resolved-node-id attribute stamped at render time
        const nodes = Array.from(containerRef.current.querySelectorAll('.node, .actor, .classGroup, .state'));
        const targetNode = nodes.find(n =>
            (n as HTMLElement).dataset.resolvedNodeId === selectedNodeId
        ) as HTMLElement | undefined;

        if (targetNode) {
            const rect = targetNode.getBoundingClientRect();
            setHandlePos({ x: rect.right + 10, y: rect.top + rect.height / 2 });
            targetNode.style.filter = 'drop-shadow(0 0 8px var(--color-primary)) brightness(1.2)';
        }

        return () => {
            if (targetNode) targetNode.style.filter = '';
        };
    }, [selectedNodeId, mermaidCode]);

    // ── Drag-to-connect tracking ────────────────────────────────────────────────
    useEffect(() => {
        if (!isDragging) return;

        const handleMouseMove = (e: MouseEvent) => setMousePos({ x: e.clientX, y: e.clientY });
        const handleMouseUp = (e: MouseEvent) => {
            setIsDragging(false);

            const elements = document.elementsFromPoint(e.clientX, e.clientY);
            let targetNode: Element | null = null;
            for (const el of elements) {
                const nodeEl = el.closest('.node, .actor, .classGroup, .state');
                if (nodeEl) { targetNode = nodeEl; break; }
            }

            if (targetNode) {
                // Look up by data-resolved-node-id stamped at render
                const resolvedId = (targetNode as HTMLElement).dataset.resolvedNodeId || '';
                if (resolvedId && resolvedId !== dragStartId) {
                    appendEdge(dragStartId!, resolvedId);
                }
            }
            // No drop-to-empty-space menu — just end the drag
        };

        window.addEventListener('mousemove', handleMouseMove);
        window.addEventListener('mouseup', handleMouseUp);
        return () => {
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseup', handleMouseUp);
        };
    }, [isDragging, dragStartId, appendEdge]);

    // ── Zoom via mouse wheel ────────────────────────────────────────────────────
    const handleWheel = useCallback((e: React.WheelEvent<HTMLDivElement>) => {
        e.preventDefault();
        const delta = e.deltaY < 0 ? SCALE_STEP : -SCALE_STEP;
        setTransform(prev => {
            const next = Math.min(SCALE_MAX, Math.max(SCALE_MIN, prev.scale + delta));
            // Keep zoom centred on cursor relative to the wrapper
            const rect = wrapperRef.current?.getBoundingClientRect();
            if (!rect) return { ...prev, scale: next };
            const cursorX = e.clientX - rect.left;
            const cursorY = e.clientY - rect.top;
            const scaleRatio = next / prev.scale;
            return {
                scale: next,
                x: cursorX - scaleRatio * (cursorX - prev.x),
                y: cursorY - scaleRatio * (cursorY - prev.y),
            };
        });
    }, []);

    // ── Pan via background mouse drag ───────────────────────────────────────────
    const handleCanvasMouseDown = (e: React.MouseEvent<HTMLDivElement>) => {
        // Only pan when clicking directly on the canvas background (not a node/handle)
        if ((e.target as Element).closest('.node, .actor, .classGroup, .state, .connection-handle')) return;
        if (isDragging) return;
        isPanning.current = true;
        panStart.current = { x: e.clientX - transformRef.current.x, y: e.clientY - transformRef.current.y };
    };

    useEffect(() => {
        const handleMouseMove = (e: MouseEvent) => {
            if (!isPanning.current) return;
            setTransform(prev => ({
                ...prev,
                x: e.clientX - panStart.current.x,
                y: e.clientY - panStart.current.y,
            }));
        };
        const handleMouseUp = () => { isPanning.current = false; };

        window.addEventListener('mousemove', handleMouseMove);
        window.addEventListener('mouseup', handleMouseUp);
        return () => {
            window.removeEventListener('mousemove', handleMouseMove);
            window.removeEventListener('mouseup', handleMouseUp);
        };
    }, []);

    // ── Zoom buttons ────────────────────────────────────────────────────────────
    const zoomIn = () => setTransform(prev => ({ ...prev, scale: Math.min(SCALE_MAX, +(prev.scale + SCALE_STEP).toFixed(2)) }));
    const zoomOut = () => setTransform(prev => ({ ...prev, scale: Math.max(SCALE_MIN, +(prev.scale - SCALE_STEP).toFixed(2)) }));
    const zoomReset = () => setTransform({ scale: 1, x: 0, y: 0 });

    return (
        <div
            ref={wrapperRef}
            className="relative w-full h-full bg-bg-base overflow-hidden select-none"
            onWheel={handleWheel}
            onMouseDown={handleCanvasMouseDown}
            onClick={(e) => {
                if ((e.target as Element).closest('.connection-handle')) return;
                // Don't clear selection if the click landed on an SVG node element
                if ((e.target as Element).closest('.node, .actor, .classGroup, .state')) return;
                setSelectedNodeId(null);
            }}
            style={{ cursor: isPanning.current ? 'grabbing' : 'grab' }}
        >
            {/* Transformed SVG container */}
            <div
                style={{
                    transform: `translate(${transform.x}px, ${transform.y}px) scale(${transform.scale})`,
                    transformOrigin: '0 0',
                    willChange: 'transform',
                    position: 'absolute',
                    top: 0,
                    left: 0,
                    width: '100%',
                    height: '100%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                }}
            >
                {error ? (
                    <div className="bg-red-900/20 text-red-400 p-4 rounded border border-red-900/50 max-w-lg">
                        <h3 className="font-bold mb-2">Syntax Error</h3>
                        <pre className="text-xs overflow-auto whitespace-pre-wrap">{error}</pre>
                    </div>
                ) : (
                    <div
                        ref={containerRef}
                        className="mermaid-container"
                    />
                )}
            </div>

            {/* Zoom controls — top-right corner */}
            <div className="absolute top-3 right-3 z-30 flex items-center gap-1 bg-bg-panel border border-border rounded-lg shadow-md p-1">
                <button
                    aria-label="Zoom In"
                    onClick={(e) => { e.stopPropagation(); zoomIn(); }}
                    className="w-8 h-8 flex items-center justify-center rounded text-text-muted hover:text-text-base hover:bg-bg-subtle transition-colors text-lg font-bold"
                >+</button>
                <button
                    aria-label="Zoom Out"
                    onClick={(e) => { e.stopPropagation(); zoomOut(); }}
                    className="w-8 h-8 flex items-center justify-center rounded text-text-muted hover:text-text-base hover:bg-bg-subtle transition-colors text-lg font-bold"
                >−</button>
                <div className="w-px h-5 bg-border mx-0.5" />
                <button
                    aria-label="Reset Zoom"
                    onClick={(e) => { e.stopPropagation(); zoomReset(); }}
                    className="w-8 h-8 flex items-center justify-center rounded text-text-muted hover:text-text-base hover:bg-bg-subtle transition-colors text-sm"
                    title="Reset zoom"
                >⟳</button>
                <span className="text-text-muted text-xs px-1 min-w-[3rem] text-center tabular-nums">
                    {Math.round(transform.scale * 100)}%
                </span>
            </div>

            {/* Connection Handle */}
            {handlePos && !isDragging && (
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
        </div>
    );
};

export default EditorCanvas;
