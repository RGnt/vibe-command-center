package models

// GraphNode represents an entity in the knowledge graph
type GraphNode struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"` // e.g. Person, Concept, Organization
	Properties map[string]interface{} `json:"properties"`
}

// GraphEdge represents a relationship in the knowledge graph
type GraphEdge struct {
	Source     string                 `json:"source"` // ID of the source node
	Target     string                 `json:"target"` // ID of the target node
	Label      string                 `json:"label"`  // e.g. RELATED_TO, DEPENDS_ON
	Properties map[string]interface{} `json:"properties"`
}

// GraphPayload is the complete graph extracted from a document
type GraphPayload struct {
	Nodes []GraphNode `json:"nodes"`
	Edges []GraphEdge `json:"edges"`
}
