import { describe, it, expect } from 'vitest';
import { exportToMermaid, importFromMermaid } from './flowchart';

describe('Flowchart Parser', () => {
    describe('exportToMermaid', () => {
        it('should correctly export basic nodes and edges', () => {
            const nodes = [
                { id: 'A', data: { label: 'Node A', shape: 'rectangle' } },
                { id: 'B', data: { label: 'Node B', shape: 'circle' } }
            ];
            const edges = [
                { source: 'A', target: 'B', label: 'Goes to' }
            ];
            const result = exportToMermaid(nodes, edges, 'graph TD');
            expect(result).toContain('graph TD');
            expect(result).toContain('A[Node A]');
            expect(result).toContain('B((Node B))');
            expect(result).toContain('A -->|Goes to| B');
        });
    });

    describe('importFromMermaid', () => {
        it('should correctly parse basic mermaid string into nodes and edges', () => {
            const code = `graph TD
    A[Node A]
    B((Node B))
    A -->|Goes to| B`;
            
            const result = importFromMermaid(code, 'graph TD');
            expect(result.nodes).toHaveLength(2);
            expect(result.edges).toHaveLength(1);
            expect(result.nodes.find(n => n.id === 'A').data.label).toBe('Node A');
            expect(result.nodes.find(n => n.id === 'B').data.shape).toBe('circle');
            expect(result.edges[0].label).toBe('Goes to');
        });
    });
});
