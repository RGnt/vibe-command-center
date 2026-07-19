import React, { useState, useEffect } from 'react';
import useMermaidStore from '../../store/mermaidStore';
import Editor from '@monaco-editor/react';

const EditorPanel = () => {
    const { mermaidCode, updateFromCode, diagramExplanation, setDiagramExplanation, diagramType, setDiagramType } = useMermaidStore();
    const [localCode, setLocalCode] = useState(mermaidCode);

    // Sync local code when mermaidCode changes from outside (e.g. Load Diagram)
    useEffect(() => {
        setLocalCode(mermaidCode);
    }, [mermaidCode]);

    const handleApply = () => {
        updateFromCode(localCode);
    };

    return (
        <div className="h-full flex flex-col bg-bg-panel border-r border-border overflow-hidden">
            {/* Header */}
            <div className="px-4 py-3 border-b border-border bg-bg-subtle shrink-0">
                <h2 className="text-text-base font-semibold text-base">Mermaid Source</h2>
                <p className="text-text-muted text-xs mt-0.5">Edit code, then click Apply</p>
            </div>

            {/* Content */}
            <div className="flex-1 flex flex-col gap-3 overflow-y-auto p-4 min-h-0">
                {/* Diagram Type */}
                <div className="shrink-0 flex flex-col gap-1.5">
                    <label className="text-text-base font-medium text-xs uppercase tracking-wide">Diagram Type</label>
                    <select
                        value={diagramType}
                        onChange={(e) => setDiagramType(e.target.value)}
                        className="w-full bg-bg-base border border-border text-text-base text-sm rounded-lg focus:ring-primary focus:border-primary p-2"
                    >
                        <option value="sequenceDiagram">Sequence Diagram</option>
                    </select>
                </div>

                {/* Editor heading + Apply button */}
                <div className="shrink-0 flex justify-between items-center">
                    <label className="text-text-base font-medium text-xs uppercase tracking-wide">Editor</label>
                    <button
                        onClick={handleApply}
                        className="bg-primary text-bg-base px-3 py-1.5 rounded text-xs font-semibold hover:bg-primary/90 transition-colors"
                    >
                        Apply Code
                    </button>
                </div>

                {/* Monaco Editor — grows to fill available space */}
                <div className="flex-1 min-h-[200px] border border-border rounded-lg overflow-hidden">
                    <Editor
                        height="100%"
                        defaultLanguage="markdown"
                        theme="vs-dark"
                        value={localCode}
                        onChange={(value) => setLocalCode(value || '')}
                        options={{
                            minimap: { enabled: false },
                            scrollBeyondLastLine: false,
                            wordWrap: 'on',
                            lineNumbers: 'on',
                            fontSize: 13,
                        }}
                    />
                </div>

                {/* Explanation */}
                <div className="shrink-0 flex flex-col gap-1.5 pb-2">
                    <label className="text-text-base font-medium text-xs uppercase tracking-wide">Explanation (optional)</label>
                    <textarea
                        value={diagramExplanation}
                        onChange={(e) => setDiagramExplanation(e.target.value)}
                        className="w-full h-20 p-3 rounded-lg bg-bg-base border border-border text-text-base text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
                        placeholder="Explain what this diagram is about..."
                    />
                </div>
            </div>
        </div>
    );
};

export default EditorPanel;
