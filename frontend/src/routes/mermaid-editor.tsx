import React from 'react';
import { createFileRoute } from '@tanstack/react-router';
import MermaidEditor from '../components/mermaid-editor/MermaidEditor';

export const Route = createFileRoute('/mermaid-editor')({
  component: MermaidEditorRoute,
});

function MermaidEditorRoute() {
  return (
    <div className="flex-1 h-full overflow-hidden">
        <MermaidEditor />
    </div>
  );
}
