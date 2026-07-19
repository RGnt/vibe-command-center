import { createFileRoute } from '@tanstack/react-router'
import { useQuery } from '@tanstack/react-query'
import { SigmaContainer, ControlsContainer, ZoomControl, FullScreenControl } from '@react-sigma/core'
import { useLoadGraph } from '@react-sigma/core'
import { useWorkerLayoutForceAtlas2 } from '@react-sigma/layout-forceatlas2'
import { MultiDirectedGraph } from 'graphology'
import { useEffect } from 'react'
import '@react-sigma/core/lib/style.css'

export const Route = createFileRoute('/graph')({
  component: GraphView,
})

// Types matching Go backend payload
type GraphNode = {
  id: string
  label: string
  properties: Record<string, any>
}

type GraphEdge = {
  source: string
  target: string
  label: string
  properties: Record<string, any>
}

type GraphPayload = {
  nodes: GraphNode[]
  edges: GraphEdge[]
}

// Colors for node types
const LABEL_COLORS: Record<string, string> = {
  Person: '#3b82f6', // blue-500
  Organization: '#ef4444', // red-500
  Concept: '#10b981', // emerald-500
  Location: '#eab308', // yellow-500
  Event: '#8b5cf6', // violet-500
  Entity: '#6b7280', // gray-500
}

function LoadGraph({ payload }: { payload: GraphPayload }) {
  const loadGraph = useLoadGraph()
  const { start, stop } = useWorkerLayoutForceAtlas2({ settings: { slowDown: 10 } })

  useEffect(() => {
    const graph = new MultiDirectedGraph()

    // Add Nodes
    payload.nodes.forEach((node) => {
      if (!graph.hasNode(node.id)) {
        const color = LABEL_COLORS[node.label] || LABEL_COLORS['Entity']
        graph.addNode(node.id, {
          x: Math.random(),
          y: Math.random(),
          size: 15,
          label: node.properties.name || node.id,
          color: color,
        })
      }
    })

    // Add Edges
    payload.edges.forEach((edge) => {
      if (graph.hasNode(edge.source) && graph.hasNode(edge.target)) {
        // graphology doesn't allow duplicate multi-edges between same nodes without MultiGraph, 
        // we'll just ignore duplicates or use addEdge with generated keys
        try {
          graph.addEdge(edge.source, edge.target, {
            size: 2,
            label: edge.label,
            color: '#cbd5e1' // slate-300
          })
        } catch (e) {
          // ignore duplicate edge error
        }
      }
    })

    loadGraph(graph)
    
    // Run layout algorithm for a short time to stabilize
    start()
    const timer = setTimeout(() => {
      stop()
    }, 2000)

    return () => {
      clearTimeout(timer)
      stop()
    }
  }, [payload, loadGraph, start, stop])

  return null
}

function GraphView() {
  const { data: graphData, isLoading, isError } = useQuery<GraphPayload>({
    queryKey: ['knowledge_graph'],
    queryFn: async () => {
      const res = await fetch('/api/graph')
      if (!res.ok) throw new Error('Failed to fetch graph data')
      return res.json()
    },
  })

  if (isLoading) {
    return (
      <div className="flex h-[calc(100vh-4rem)] items-center justify-center">
        <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-blue-600"></div>
      </div>
    )
  }

  if (isError) {
    return (
      <div className="flex h-[calc(100vh-4rem)] items-center justify-center text-red-500">
        Failed to load the knowledge graph.
      </div>
    )
  }

  return (
    <div className="h-[calc(100vh-4rem)] w-full flex flex-col">
      <div className="p-4 bg-white dark:bg-gray-800 border-b dark:border-gray-700 flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-gray-800 dark:text-gray-100">Knowledge Graph Exploration</h1>
          <p className="text-sm text-gray-500 dark:text-gray-400">
            Interactive visualization of extracted entities and relationships
          </p>
        </div>
        <div className="flex gap-2 text-sm">
          {Object.entries(LABEL_COLORS).map(([label, color]) => (
            <div key={label} className="flex items-center gap-1 bg-gray-50 dark:bg-gray-900 px-2 py-1 rounded">
              <div className="w-3 h-3 rounded-full" style={{ backgroundColor: color }}></div>
              <span className="text-gray-700 dark:text-gray-300">{label}</span>
            </div>
          ))}
        </div>
      </div>
      
      <div className="flex-1 bg-slate-50 dark:bg-slate-900 relative">
        <SigmaContainer style={{ height: "100%", width: "100%" }} settings={{ allowInvalidContainer: true, labelRenderedSizeThreshold: 12, renderEdgeLabels: true, defaultEdgeType: "arrow" }}>
          {graphData && <LoadGraph payload={graphData} />}
          <ControlsContainer position={"bottom-right"}>
            <ZoomControl />
            <FullScreenControl />
          </ControlsContainer>
        </SigmaContainer>
      </div>
    </div>
  )
}
