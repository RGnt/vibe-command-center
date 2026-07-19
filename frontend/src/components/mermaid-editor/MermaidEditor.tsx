import React, { useState, useEffect, useRef, useCallback } from 'react';
import EditorPanel from './EditorPreview';
import EditorCanvas from './EditorCanvas';
import useMermaidStore from '../../store/mermaidStore';

const MIN_LEFT_PCT = 20;
const MAX_LEFT_PCT = 70;
const DEFAULT_LEFT_PCT = 35;

const MermaidEditor = () => {
    const { diagramId, diagramName, diagramType, mermaidCode, diagramExplanation, setDiagramName, loadDiagram } = useMermaidStore();
    const [savedDiagrams, setSavedDiagrams] = useState<Array<{ id: number; name: string; diagram_type: string; code: string; explanation: string }>>([]);
    const [leftPct, setLeftPct] = useState(DEFAULT_LEFT_PCT);
    const isDraggingDivider = useRef(false);
    const containerRef = useRef<HTMLDivElement>(null);

    const fetchDiagrams = async () => {
        try {
            const res = await fetch('/api/diagrams');
            if (res.ok) setSavedDiagrams(await res.json() || []);
        } catch (err) {
            console.error('Failed to fetch diagrams', err);
        }
    };

    useEffect(() => { fetchDiagrams(); }, []);

    const handleSave = async () => {
        let finalName = diagramName;
        if (!diagramId && diagramName === 'Untitled Diagram') {
            const newName = prompt('Enter diagram name:', diagramName);
            if (!newName) return;
            finalName = newName;
            setDiagramName(newName);
        }

        const payload = { name: finalName, diagram_type: diagramType, code: mermaidCode, explanation: diagramExplanation };

        try {
            const url = diagramId ? `/api/diagrams/${diagramId}` : '/api/diagrams';
            const method = diagramId ? 'PUT' : 'POST';
            const res = await fetch(url, { method, headers: { 'Content-Type': 'application/json' }, body: JSON.stringify(payload) });
            if (res.ok) {
                const data = await res.json();
                loadDiagram(data);
                fetchDiagrams();
                alert('Diagram saved successfully!');
            } else {
                alert('Failed to save diagram');
            }
        } catch (err) {
            console.error(err);
            alert('Error saving diagram');
        }
    };

    // ── Resizable divider ────────────────────────────────────────────────────────
    const startDividerDrag = useCallback((e: React.MouseEvent) => {
        e.preventDefault();
        isDraggingDivider.current = true;

        const onMouseMove = (ev: MouseEvent) => {
            if (!isDraggingDivider.current || !containerRef.current) return;
            const rect = containerRef.current.getBoundingClientRect();
            const pct = ((ev.clientX - rect.left) / rect.width) * 100;
            setLeftPct(Math.min(MAX_LEFT_PCT, Math.max(MIN_LEFT_PCT, pct)));
        };

        const onMouseUp = () => {
            isDraggingDivider.current = false;
            window.removeEventListener('mousemove', onMouseMove);
            window.removeEventListener('mouseup', onMouseUp);
        };

        window.addEventListener('mousemove', onMouseMove);
        window.addEventListener('mouseup', onMouseUp);
    }, []);

    return (
        <div className="flex flex-col h-screen bg-bg-base text-text-base overflow-hidden">
            {/* ── Top bar ─────────────────────────────────────────────────────── */}
            <div className="h-14 border-b border-border flex items-center justify-between px-6 bg-bg-subtle z-[60] shrink-0">
                <div className="flex items-center gap-4">
                    <input
                        value={diagramName}
                        onChange={(e) => setDiagramName(e.target.value)}
                        className="bg-transparent border-none focus:outline-none focus:ring-0 font-bold text-xl w-48 text-text-base"
                    />
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
                            e.target.value = '';
                        }}
                        className="bg-bg-base border border-border text-text-base text-sm rounded-lg focus:ring-primary focus:border-primary block p-2"
                        value=""
                    >
                        <option value="" disabled>Load Diagram...</option>
                        {savedDiagrams.map(d => (
                            <option key={d.id} value={d.id}>{d.name}</option>
                        ))}
                    </select>
                </div>
            </div>

            {/* ── Split pane ───────────────────────────────────────────────────── */}
            <div ref={containerRef} className="flex-1 flex overflow-hidden">
                {/* Left pane: code editor */}
                <div style={{ width: `${leftPct}%` }} className="flex flex-col shrink-0 min-w-0">
                    <EditorPanel />
                </div>

                {/* Draggable divider */}
                <div
                    className="w-1 shrink-0 bg-border hover:bg-primary cursor-col-resize transition-colors relative group"
                    onMouseDown={startDividerDrag}
                    title="Drag to resize"
                >
                    {/* Visual grip dots */}
                    <div className="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 flex flex-col gap-1 opacity-40 group-hover:opacity-80 transition-opacity pointer-events-none">
                        <div className="w-1 h-1 rounded-full bg-text-muted" />
                        <div className="w-1 h-1 rounded-full bg-text-muted" />
                        <div className="w-1 h-1 rounded-full bg-text-muted" />
                    </div>
                </div>

                {/* Right pane: live canvas */}
                <div style={{ width: `${100 - leftPct}%` }} className="flex flex-col min-w-0 overflow-hidden">
                    <EditorCanvas />
                </div>
            </div>
        </div>
    );
};

export default MermaidEditor;
