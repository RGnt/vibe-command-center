import React from 'react';
import useMermaidStore from '../../store/mermaidStore';
import MarkdownRenderer from '../MarkdownRenderer';

const EditorPreview = () => {
    const { mermaidCode, updateFromCode } = useMermaidStore();

    return (
        <div className="w-96 bg-bg-base border-l border-border flex flex-col h-full">
            <div className="p-4 border-b border-border">
                <h2 className="text-text-base font-semibold text-lg">Mermaid Source</h2>
                <p className="text-text-muted text-sm">Edit code here to update the canvas, or copy this code to use elsewhere.</p>
            </div>
            
            <div className="flex-1 p-4 flex flex-col gap-4 overflow-y-auto">
                <textarea
                    value={mermaidCode}
                    onChange={(e) => updateFromCode(e.target.value)}
                    className="w-full h-48 p-3 rounded-lg bg-bg-subtle border border-border text-text-base font-mono text-sm focus:outline-none focus:border-primary resize-y"
                    spellCheck="false"
                />

                <div className="border-t border-border pt-4">
                    <h3 className="text-text-base font-medium mb-2">Explanation (Optional)</h3>
                    <textarea
                        value={useMermaidStore((state) => state.diagramExplanation)}
                        onChange={(e) => useMermaidStore.getState().setDiagramExplanation(e.target.value)}
                        className="w-full h-24 p-3 rounded-lg bg-bg-subtle border border-border text-text-base text-sm focus:outline-none focus:border-primary resize-y"
                        placeholder="Explain what this diagram is about..."
                    />
                </div>
                
                <div className="border-t border-border pt-4">
                    <h3 className="text-text-base font-medium mb-2">Live Render</h3>
                    <div className="bg-bg-subtle rounded-lg border border-border p-2">
                        {/* We use MarkdownRenderer with a mermaid block */}
                        <MarkdownRenderer content={`\`\`\`mermaid\n${mermaidCode}\n\`\`\``} />
                    </div>
                </div>
            </div>
        </div>
    );
};

export default EditorPreview;
