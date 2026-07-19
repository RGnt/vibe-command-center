export const exportToMermaid = (nodes, edges, diagramType) => {
    let code = `${diagramType}\n`;

    // Order nodes by X coordinate to maintain left-to-right participant order
    const sortedNodes = [...nodes].sort((a, b) => a.position.x - b.position.x);

    const nodeStrs = sortedNodes.map(node => {
        const shape = node.data.shape || 'participant';
        const label = node.data.label || node.id;
        
        // Output participant or actor
        // If label is different from id, use "as"
        if (label !== node.id) {
            return `    ${shape} ${node.id} as ${label}`;
        }
        return `    ${shape} ${node.id}`;
    });

    // Edges (messages) should ideally be ordered top-to-bottom.
    // If edges don't have distinct Y coordinates (since React Flow might not guarantee it),
    // we just map them in order of creation.
    const edgeStrs = edges.map(edge => {
        const type = edge.data?.sequenceType || '->'; // Default to solid line without arrow
        
        if (edge.label) {
            return `    ${edge.source}${type}${edge.target}: ${edge.label}`;
        }
        return `    ${edge.source}${type}${edge.target}`;
    });

    code += nodeStrs.join('\n');
    if (nodes.length > 0 && edges.length > 0) code += '\n';
    code += edgeStrs.join('\n');

    return code;
};

export const importFromMermaid = (code, diagramType) => {
    const nodes = [];
    const edges = [];

    const lines = code.split('\n').map(l => l.trim()).filter(l => l.length > 0);
    const nodeMap = new Map();
    let edgeIdCounter = 1;

    let xOffset = 0;

    const addNode = (id, label, shape = 'participant') => {
        if (!nodeMap.has(id)) {
            const newNode = {
                id,
                type: 'mermaidNode',
                position: { x: xOffset, y: 100 },
                data: { label: label || id, shape }
            };
            xOffset += 200; // Layout them horizontally
            nodeMap.set(id, newNode);
            nodes.push(newNode);
        } else if (label && nodeMap.get(id).data.label === id) {
            nodeMap.get(id).data.label = label;
            nodeMap.get(id).data.shape = shape;
        }
    };

    for (const line of lines) {
        if (line === 'sequenceDiagram') continue;

        let match;

        // Check for participant/actor declarations
        if ((match = line.match(/^(participant|actor)\s+([a-zA-Z0-9_]+)(?:\s+as\s+(.*))?$/))) {
            const [, shape, id, label] = match;
            addNode(id, label, shape);
            continue;
        }

        // Check for messages
        // Match arrows like ->, -->, ->>, -->>, -x, --x, -), --)
        const arrowMatch = line.match(/^([a-zA-Z0-9_]+)\s*(->|-->|->>|-->>|-x|--x|-\)|--\))\s*([a-zA-Z0-9_]+)(?::\s*(.*))?$/);
        if (arrowMatch) {
            const [, source, arrow, target, label] = arrowMatch;
            addNode(source, null, 'participant');
            addNode(target, null, 'participant');

            edges.push({
                id: `edge_${edgeIdCounter++}`,
                source,
                target,
                label: label || '',
                type: 'smoothstep',
                data: { sequenceType: arrow }
            });
            continue;
        }
    }

    return { nodes, edges, diagramType };
};
