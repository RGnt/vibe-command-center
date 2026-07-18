import React, { useCallback, useRef, useState, useEffect } from 'react';
import {
    ReactFlow,
    Background,
    Controls,
    MiniMap,
    ReactFlowProvider,
    Handle,
    Position
} from '@xyflow/react';
import '@xyflow/react/dist/style.css';
import useMermaidStore from '../../store/mermaidStore';

let id = 0;
const getId = () => `node_${id++}`;

const MermaidNode = ({ data, id }) => {
    const { updateNodeLabel } = useMermaidStore();

    // Check if label contains an image
    const isImage = data.label && data.label.includes('<img');
    let imgSrc = '';
    if (isImage) {
        const match = data.label.match(/src=['"](.*?)['"]/);
        if (match) imgSrc = match[1];
    }

    // Basic styling based on shape
    let style = {
        padding: isImage ? '0' : '10px 20px',
        background: isImage ? 'transparent' : 'var(--color-bg-base)',
        border: isImage ? 'none' : '2px solid var(--color-primary)',
        color: 'var(--color-text-base)',
        minWidth: isImage ? 'auto' : '100px',
        textAlign: 'center',
    };

    if (!isImage) {
        if (data.shape === 'round') style.borderRadius = '20px';
        if (data.shape === 'stadium') style.borderRadius = '50% 50% 50% 50% / 15% 15% 15% 15%';
        if (data.shape === 'circle') {
            style.borderRadius = '50%';
            style.aspectRatio = '1/1';
            style.display = 'flex';
            style.alignItems = 'center';
            style.justifyContent = 'center';
        }
        if (data.shape === 'cylinder') style.borderRadius = '10px';
        if (data.shape === 'rhombus') {
            style.transform = 'skewX(-15deg)';
        }
    }

    return (
        <div style={style}>
            <Handle type="target" position={Position.Top} />
            <div style={(!isImage && data.shape === 'rhombus') ? { transform: 'skewX(15deg)' } : {}}>
                {isImage ? (
                    <img src={imgSrc} alt="Diagram element" style={{ maxWidth: '100px', maxHeight: '100px', borderRadius: '8px' }} />
                ) : (
                    <input
                        type="text"
                        value={data.label || ''}
                        onChange={(e) => updateNodeLabel(id, e.target.value)}
                        className="bg-transparent text-center focus:outline-none w-full"
                        placeholder="Label..."
                    />
                )}
            </div>
            <Handle type="source" position={Position.Bottom} />
        </div>
    );
};

const nodeTypes = {
    mermaidNode: MermaidNode,
};

const EditorCanvas = () => {
    const reactFlowWrapper = useRef(null);
    const [reactFlowInstance, setReactFlowInstance] = useState(null);
    
    const { nodes, edges, onNodesChange, onEdgesChange, onConnect, addNode } = useMermaidStore();

    useEffect(() => {
        const handleAddImage = (e) => {
            const newNodeId = getId();
            addNode({
                id: newNodeId,
                type: 'mermaidNode',
                position: { x: 100, y: 100 },
                data: { label: `<img src='${e.detail.url}' width='100' height='100' />`, shape: 'rectangle' }
            });
        };
        window.addEventListener('addImageNode', handleAddImage);
        return () => window.removeEventListener('addImageNode', handleAddImage);
    }, [addNode]);

    const onDragOver = useCallback((event) => {
        event.preventDefault();
        event.dataTransfer.dropEffect = 'move';
    }, []);

    const onDrop = useCallback(
        (event) => {
            event.preventDefault();

            const type = event.dataTransfer.getData('application/reactflow');
            const shape = event.dataTransfer.getData('application/reactflow/shape');
            const url = event.dataTransfer.getData('application/reactflow/url');

            if (typeof type === 'undefined' || !type) {
                return;
            }

            const position = reactFlowInstance.screenToFlowPosition({
                x: event.clientX,
                y: event.clientY,
            });

            const newNodeId = getId();
            const label = url ? `<img src='${url}' width='100' height='100' />` : `Node ${newNodeId.split('_')[1]}`;
            
            const newNode = {
                id: newNodeId,
                type,
                position,
                data: { label, shape },
            };

            addNode(newNode);
        },
        [reactFlowInstance, addNode],
    );

    return (
        <div className="flex-1 h-full" ref={reactFlowWrapper}>
            <ReactFlow
                nodes={nodes}
                edges={edges}
                onNodesChange={onNodesChange}
                onEdgesChange={onEdgesChange}
                onConnect={onConnect}
                onInit={setReactFlowInstance}
                onDrop={onDrop}
                onDragOver={onDragOver}
                nodeTypes={nodeTypes}
                fitView
            >
                <Background variant="dots" gap={12} size={1} />
                <Controls />
                <MiniMap />
            </ReactFlow>
        </div>
    );
};

// Wrap with provider
export default () => (
    <ReactFlowProvider>
        <EditorCanvas />
    </ReactFlowProvider>
);
