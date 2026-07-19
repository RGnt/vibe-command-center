import { create } from 'zustand';

const useMermaidStore = create((set, get) => ({
    diagramId: null,
    diagramName: 'Untitled Diagram',
    diagramExplanation: '',
    diagramType: 'sequenceDiagram',
    mermaidCode: 'sequenceDiagram\n  actor Alice\n  actor Bob\n',
    selectedNodeId: null,

    setDiagramId: (id) => set({ diagramId: id }),
    setDiagramName: (name) => set({ diagramName: name }),
    setDiagramExplanation: (exp) => set({ diagramExplanation: exp }),
    setSelectedNodeId: (id) => set({ selectedNodeId: id }),
    
    setDiagramType: (type) => {
        let defaultCode = 'graph TD\n';
        if (type === 'sequenceDiagram') defaultCode = 'sequenceDiagram\n';
        else if (type === 'classDiagram') defaultCode = 'classDiagram\n';
        else if (type === 'stateDiagram-v2') defaultCode = 'stateDiagram-v2\n';
        
        set({ diagramType: type, mermaidCode: defaultCode });
    },

    appendCode: (fragment) => {
        set((state) => ({ mermaidCode: state.mermaidCode + '\n' + fragment }));
    },
    
    appendEdge: (sourceId, targetId) => {
        const type = get().diagramType;
        if (type === 'sequenceDiagram') {
            get().appendCode(`${sourceId} ->> ${targetId}: message`);
            return;
        }
        let edgeSyntax = '-->';
        if (type === 'classDiagram') edgeSyntax = '<|--';
        
        get().appendCode(`${sourceId} ${edgeSyntax} ${targetId}`);
    },

    updateFromCode: (code) => {
        set({ mermaidCode: code });
    },

    loadDiagram: (diagram) => {
        set({
            diagramId: diagram.id,
            diagramName: diagram.name,
            diagramExplanation: diagram.explanation,
            diagramType: diagram.diagram_type,
            mermaidCode: diagram.code
        });
    }
}));

export default useMermaidStore;
