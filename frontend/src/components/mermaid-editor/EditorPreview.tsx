import React, { useState, useEffect } from 'react';
import useMermaidStore from '../../store/mermaidStore';

import Editor from '@monaco-editor/react';

const EditorPreview = () => {
    const { mermaidCode, updateFromCode, diagramExplanation, setDiagramExplanation, diagramType, setDiagramType } = useMermaidStore();
    const [isCollapsed, setIsCollapsed] = useState(false);
    const [localCode, setLocalCode] = useState(mermaidCode);

    // Sync local code if mermaidCode is updated from elsewhere (e.g. Toolbar, Load)
    useEffect(() => {
        setLocalCode(mermaidCode);
    }, [mermaidCode]);

    const handleApply = () => {
        updateFromCode(localCode);
    };

    if (isCollapsed) {
        return (
            <div className="absolute top-20 left-4 z-40">
                <button
                    onClick={() => setIsCollapsed(false)}
                    className="bg-bg-panel border border-border shadow-lg p-3 rounded-lg text-text-base hover:bg-bg-subtle transition-colors flex items-center gap-2"
                    aria-label="Toggle Code Editor"
                    data-testid="expand-code-btn"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="16 18 22 12 16 6"></polyline><polyline points="8 6 2 12 8 18"></polyline></svg>
                    <span className="font-medium">Code &lt;/&gt;</span>
                </button>
            </div>
        );
    }

    return (
        <div className="absolute top-20 left-4 z-40 w-[640px] max-h-[calc(100vh-120px)] bg-bg-panel border border-border rounded-xl shadow-2xl flex flex-col overflow-hidden">
            <div className="p-4 border-b border-border flex justify-between items-center bg-bg-subtle">
                <div>
                    <h2 className="text-text-base font-semibold text-lg">Mermaid Source</h2>
                    <p className="text-text-muted text-xs mt-1">Manual sync with canvas</p>
                </div>
                <button
                    onClick={() => setIsCollapsed(true)}
                    className="p-2 text-text-muted hover:text-text-base hover:bg-bg-base rounded-lg transition-colors"
                    aria-label="Toggle Code Editor"
                    data-testid="collapse-code-btn"
                >
                    <svg xmlns="http://www.w3.org/2000/svg" width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round"><polyline points="15 18 9 12 15 6"></polyline></svg>
                </button>
            </div>

            <div className="flex-1 p-4 flex flex-col gap-4 overflow-y-auto">
                <div className="flex flex-col gap-1.5 shrink-0 border-b border-border pb-3">
                    <h3 className="text-text-base font-medium text-sm">Diagram Type</h3>
                    <select
                        value={diagramType}
                        onChange={(e) => setDiagramType(e.target.value)}
                        className="w-full bg-bg-base border border-border text-text-base text-sm rounded-lg focus:ring-primary focus:border-primary block p-2"
                    >
                        <option value="sequenceDiagram">Sequence Diagram (sequenceDiagram)</option>
                    </select>
                </div>

                <div className="flex justify-between items-center">
                    <h3 className="text-text-base font-medium text-sm">Editor</h3>
                    <button
                        onClick={handleApply}
                        className="bg-primary text-bg-base px-3 py-1.5 rounded text-sm font-medium hover:bg-primary/90 transition-colors"
                    >
                        Apply Code
                    </button>
                </div>

                <div className="w-full h-64 border border-border rounded-lg overflow-hidden shrink-0">
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
                        }}
                    />
                </div>

                <div className="border-t border-border pt-4 shrink-0 pb-4">
                    <h3 className="text-text-base font-medium mb-2 text-sm">Explanation (Optional)</h3>
                    <textarea
                        value={diagramExplanation}
                        onChange={(e) => setDiagramExplanation(e.target.value)}
                        className="w-full h-24 p-3 rounded-lg bg-bg-base border border-border text-text-base text-sm focus:outline-none focus:ring-2 focus:ring-primary focus:border-transparent resize-none"
                        placeholder="Explain what this diagram is about..."
                    />
                </div>
            </div>
        </div>
    );
};

export default EditorPreview;
