import { exportToMermaid as exportFlowchart, importFromMermaid as importFlowchart } from './parsers/flowchart';
import { exportToMermaid as exportSequence, importFromMermaid as importSequence } from './parsers/sequence';

export const exportToMermaid = (nodes, edges, diagramType = 'graph TD') => {
    if (diagramType.startsWith('graph') || diagramType.startsWith('flowchart')) {
        return exportFlowchart(nodes, edges, diagramType);
    }
    
    if (diagramType === 'sequenceDiagram') {
        return exportSequence(nodes, edges, diagramType);
    }
    
    // Fallback for unimplemented parsers
    return `${diagramType}\n`;
};

export const importFromMermaid = (code) => {
    const lines = code.split('\n').map(l => l.trim()).filter(l => l.length > 0);
    let diagramType = 'graph TD';
    
    if (lines.length > 0) {
        diagramType = lines[0];
    }

    if (diagramType.startsWith('graph ') || diagramType.startsWith('flowchart ')) {
        return importFlowchart(code, diagramType);
    }

    if (diagramType === 'sequenceDiagram') {
        return importSequence(code, diagramType);
    }

    // Fallback for unimplemented parsers
    return { nodes: [], edges: [], diagramType };
};
