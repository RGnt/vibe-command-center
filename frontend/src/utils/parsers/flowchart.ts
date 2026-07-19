import dagre from 'dagre';

const shapeMap = {
    'rectangle': { left: '[', right: ']' },
    'round': { left: '(', right: ')' },
    'rhombus': { left: '{', right: '}' },
    'stadium': { left: '([', right: '])' },
    'cylinder': { left: '[(', right: ')]' },
    'circle': { left: '((', right: '))' },
};

export const exportToMermaid = (nodes, edges, diagramType) => {
    let code = `${diagramType}\n`;

    const nodeStrs = nodes.map(node => {
        const shape = node.data.shape || 'rectangle';
        const label = node.data.label || node.id;
        const delimiters = shapeMap[shape] || shapeMap['rectangle'];
        return `    ${node.id}${delimiters.left}${label}${delimiters.right}`;
    });

    const edgeStrs = edges.map(edge => {
        let linkStr = '-->';
        if (edge.label) {
            linkStr = `-->|${edge.label}|`;
        }
        return `    ${edge.source} ${linkStr} ${edge.target}`;
    });

    code += nodeStrs.join('\n');
    if (nodes.length > 0 && edges.length > 0) code += '\n';
    code += edgeStrs.join('\n');

    return code;
};

export const getLayoutedElements = (nodes, edges, direction = 'TB') => {
    const dagreGraph = new dagre.graphlib.Graph();
    dagreGraph.setDefaultEdgeLabel(() => ({}));
    dagreGraph.setGraph({ rankdir: direction, nodesep: 60, ranksep: 100 });

    nodes.forEach((node) => {
        dagreGraph.setNode(node.id, { width: 150, height: 50 });
    });

    edges.forEach((edge) => {
        dagreGraph.setEdge(edge.source, edge.target);
    });

    dagre.layout(dagreGraph);

    const layoutedNodes = nodes.map((node) => {
        const nodeWithPosition = dagreGraph.node(node.id);
        return {
            ...node,
            position: {
                x: nodeWithPosition.x - 75,
                y: nodeWithPosition.y - 25,
            },
        };
    });

    return { nodes: layoutedNodes, edges };
};

export const importFromMermaid = (code, diagramType) => {
    const nodes = [];
    const edges = [];
    let direction = 'TD';

    if (diagramType.startsWith('graph ') || diagramType.startsWith('flowchart ')) {
        direction = diagramType.split(' ')[1] || 'TD';
    }

    const lines = code.split('\n').map(l => l.trim()).filter(l => l.length > 0);
    const nodeMap = new Map();

    const addNode = (id, label, shape = 'rectangle') => {
        if (!nodeMap.has(id)) {
            const newNode = {
                id,
                type: 'mermaidNode',
                position: { x: 0, y: 0 },
                data: { label: label || id, shape }
            };
            nodeMap.set(id, newNode);
            nodes.push(newNode);
        } else if (label && nodeMap.get(id).data.label === id) {
            nodeMap.get(id).data.label = label;
            nodeMap.get(id).data.shape = shape;
        }
    };

    const parseNodeSyntax = (str) => {
        let match;
        if ((match = str.match(/^([a-zA-Z0-9_]+)\[\((.*?)\)\]$/))) return { id: match[1], label: match[2], shape: 'cylinder' };
        if ((match = str.match(/^([a-zA-Z0-9_]+)\(\[(.*?)\]\)$/))) return { id: match[1], label: match[2], shape: 'stadium' };
        if ((match = str.match(/^([a-zA-Z0-9_]+)\(\((.*?)\)\)$/))) return { id: match[1], label: match[2], shape: 'circle' };
        if ((match = str.match(/^([a-zA-Z0-9_]+)\[(.*?)\]$/))) return { id: match[1], label: match[2], shape: 'rectangle' };
        if ((match = str.match(/^([a-zA-Z0-9_]+)\((.*?)\)$/))) return { id: match[1], label: match[2], shape: 'round' };
        if ((match = str.match(/^([a-zA-Z0-9_]+)\{(.*?)\}$/))) return { id: match[1], label: match[2], shape: 'rhombus' };
        
        if (str.match(/^[a-zA-Z0-9_]+$/)) return { id: str, label: str, shape: 'rectangle' };
        return null;
    };

    lines.forEach(line => {
        if (line.startsWith('graph ') || line.startsWith('flowchart ')) return;

        const edgeMatch = line.match(/^(.*?)\s*-->(\|([^|]+)\|)?\s*(.*?)$/);
        if (edgeMatch) {
            const sourceStr = edgeMatch[1].trim();
            const edgeLabel = edgeMatch[3] ? edgeMatch[3].trim() : null;
            const targetStr = edgeMatch[4].trim();

            const sourceParsed = parseNodeSyntax(sourceStr);
            const targetParsed = parseNodeSyntax(targetStr);

            if (sourceParsed) addNode(sourceParsed.id, sourceParsed.label, sourceParsed.shape);
            if (targetParsed) addNode(targetParsed.id, targetParsed.label, targetParsed.shape);

            if (sourceParsed && targetParsed) {
                edges.push({
                    id: `e-${sourceParsed.id}-${targetParsed.id}-${edges.length}`,
                    source: sourceParsed.id,
                    target: targetParsed.id,
                    label: edgeLabel,
                    type: 'smoothstep'
                });
            }
        } else {
            const nodeParsed = parseNodeSyntax(line);
            if (nodeParsed) {
                addNode(nodeParsed.id, nodeParsed.label, nodeParsed.shape);
            }
        }
    });

    const layout = getLayoutedElements(nodes, edges, direction);
    return { ...layout, diagramType };
};
