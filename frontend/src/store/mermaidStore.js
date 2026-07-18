import { create } from 'zustand';
import {
    addEdge,
    applyNodeChanges,
    applyEdgeChanges,
} from '@xyflow/react';
import { exportToMermaid, importFromMermaid } from '../utils/mermaidParser';

const useMermaidStore = create((set, get) => ({
    nodes: [],
    edges: [],
    diagramId: null,
    diagramName: 'Untitled Diagram',
    diagramExplanation: '',
    diagramType: 'graph TD',
    mermaidCode: 'graph TD\n',

    setDiagramId: (id) => set({ diagramId: id }),
    setDiagramName: (name) => set({ diagramName: name }),
    setDiagramExplanation: (exp) => set({ diagramExplanation: exp }),
    
    setDiagramType: (type) => {
        set({ diagramType: type });
        get().generateCode();
    },

    onNodesChange: (changes) => {
        set({
            nodes: applyNodeChanges(changes, get().nodes),
        });
        get().generateCode();
    },
    onEdgesChange: (changes) => {
        set({
            edges: applyEdgeChanges(changes, get().edges),
        });
        get().generateCode();
    },
    onConnect: (connection) => {
        set({
            edges: addEdge({ ...connection, type: 'smoothstep' }, get().edges),
        });
        get().generateCode();
    },
    addNode: (node) => {
        set({
            nodes: [...get().nodes, node],
        });
        get().generateCode();
    },
    updateNodeLabel: (id, label) => {
        set({
            nodes: get().nodes.map((n) => {
                if (n.id === id) {
                    return { ...n, data: { ...n.data, label } };
                }
                return n;
            }),
        });
        get().generateCode();
    },
    updateNodeShape: (id, shape) => {
        set({
            nodes: get().nodes.map((n) => {
                if (n.id === id) {
                    return { ...n, data: { ...n.data, shape } };
                }
                return n;
            }),
        });
        get().generateCode();
    },
    generateCode: () => {
        const { nodes, edges, diagramType } = get();
        const code = exportToMermaid(nodes, edges, diagramType);
        set({ mermaidCode: code });
    },
    updateFromCode: (code) => {
        try {
            const { nodes, edges, diagramType } = importFromMermaid(code);
            set({ nodes, edges, mermaidCode: code, diagramType });
        } catch (error) {
            console.error("Failed to parse mermaid code:", error);
            // We only update code on error so user can keep typing, but nodes won't change
            set({ mermaidCode: code });
        }
    },
    loadDiagram: (diagram) => {
        set({
            diagramId: diagram.id,
            diagramName: diagram.name,
            diagramExplanation: diagram.explanation,
            diagramType: diagram.diagram_type
        });
        get().updateFromCode(diagram.code);
    }
}));

export default useMermaidStore;
