package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"regexp"

	"todo-backend/internal/models"
)

// GraphRepository ...
type GraphRepository interface {
	UpsertGraph(payload models.GraphPayload) error
	GetGraph() (models.GraphPayload, error)
}

type graphRepository struct {
	db *sql.DB
}

// NewGraphRepository ...
func NewGraphRepository(db *sql.DB) GraphRepository {
	return &graphRepository{db: db}
}

// sanitizeLabel ensures the label only contains alphanumeric chars to prevent cypher injection
func sanitizeLabel(label string) string {
	re := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
	if !re.MatchString(label) {
		return "Entity" // Fallback label
	}
	return label
}

// UpsertGraph ...
func (r *graphRepository) UpsertGraph(payload models.GraphPayload) error {
	// Execute in a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	// Make sure search path is set for this tx
	_, err = tx.Exec(`SET search_path = ag_catalog, "$user", public;`)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	// 1. Upsert Nodes
	for _, node := range payload.Nodes {
		label := sanitizeLabel(node.Label)
		if label == "" {
			label = "Entity"
		}
		
		// Ensure properties is not nil
		if node.Properties == nil {
			node.Properties = make(map[string]interface{})
		}

		// Prepare parameter map
		params := map[string]interface{}{
			"id":         node.ID,
			"properties": node.Properties,
		}
		paramsJSON, _ := json.Marshal(params)

		// AGE allows passing agtype params as the 3rd argument
		query := fmt.Sprintf(`
			SELECT * FROM cypher('knowledge_graph', $$
				MERGE (n:%s {id: $id})
				SET n += $properties
				RETURN n
			$$, $1) as (v agtype);
		`, label)

		_, err = tx.Exec(query, string(paramsJSON))
		if err != nil {
			log.Printf("Failed to upsert node %s: %v", node.ID, err)
			// Continue with other nodes instead of completely failing
		}
	}

	// 2. Upsert Edges
	for _, edge := range payload.Edges {
		label := sanitizeLabel(edge.Label)
		if label == "" {
			label = "RELATED_TO"
		}

		if edge.Properties == nil {
			edge.Properties = make(map[string]interface{})
		}

		params := map[string]interface{}{
			"source":     edge.Source,
			"target":     edge.Target,
			"properties": edge.Properties,
		}
		paramsJSON, _ := json.Marshal(params)

		// Match source and target nodes, then MERGE the relationship
		query := fmt.Sprintf(`
			SELECT * FROM cypher('knowledge_graph', $$
				MATCH (a {id: $source}), (b {id: $target})
				MERGE (a)-[e:%s]->(b)
				SET e += $properties
				RETURN e
			$$, $1) as (v agtype);
		`, label)

		_, err = tx.Exec(query, string(paramsJSON))
		if err != nil {
			log.Printf("Failed to upsert edge %s->%s: %v", edge.Source, edge.Target, err)
		}
	}

	return tx.Commit()
}

// GetGraph ...
func (r *graphRepository) GetGraph() (models.GraphPayload, error) {
	var payload models.GraphPayload
	payload.Nodes = make([]models.GraphNode, 0)
	payload.Edges = make([]models.GraphEdge, 0)

	// Make sure search path is set
	_, err := r.db.Exec(`SET search_path = ag_catalog, "$user", public;`)
	if err != nil {
		return payload, err
	}

	// 1. Fetch Nodes
	// We extract id, label, and properties
	nodeQuery := `
		SELECT * FROM cypher('knowledge_graph', $$
			MATCH (n)
			RETURN properties(n).id, label(n), properties(n)
		$$) as (id agtype, label agtype, props agtype);
	`
	rows, err := r.db.Query(nodeQuery)
	if err != nil {
		return payload, fmt.Errorf("failed to query nodes: %w", err)
	}
	defer func() { _ = rows.Close() }()

	for rows.Next() {
		var idAg, labelAg, propsAg string
		if err := rows.Scan(&idAg, &labelAg, &propsAg); err != nil {
			log.Println("Error scanning node row:", err)
			continue
		}

		// agtype strings are wrapped in double quotes (e.g. '"Person"'). We must unmarshal them.
		var idStr, labelStr string
		_ = json.Unmarshal([]byte(idAg), &idStr)
		_ = json.Unmarshal([]byte(labelAg), &labelStr)

		var props map[string]interface{}
		_ = json.Unmarshal([]byte(propsAg), &props)

		if idStr == "" {
			continue // skip nodes without our custom ID
		}

		payload.Nodes = append(payload.Nodes, models.GraphNode{
			ID:         idStr,
			Label:      labelStr,
			Properties: props,
		})
	}

	// 2. Fetch Edges
	edgeQuery := `
		SELECT * FROM cypher('knowledge_graph', $$
			MATCH (n)-[e]->(m)
			RETURN properties(n).id, properties(m).id, label(e), properties(e)
		$$) as (source agtype, target agtype, label agtype, props agtype);
	`
	edgeRows, err := r.db.Query(edgeQuery)
	if err != nil {
		return payload, fmt.Errorf("failed to query edges: %w", err)
	}
	defer func() { _ = edgeRows.Close() }()

	for edgeRows.Next() {
		var srcAg, tgtAg, labelAg, propsAg string
		if err := edgeRows.Scan(&srcAg, &tgtAg, &labelAg, &propsAg); err != nil {
			log.Println("Error scanning edge row:", err)
			continue
		}

		var srcStr, tgtStr, labelStr string
		_ = json.Unmarshal([]byte(srcAg), &srcStr)
		_ = json.Unmarshal([]byte(tgtAg), &tgtStr)
		_ = json.Unmarshal([]byte(labelAg), &labelStr)

		var props map[string]interface{}
		_ = json.Unmarshal([]byte(propsAg), &props)

		if srcStr == "" || tgtStr == "" {
			continue
		}

		payload.Edges = append(payload.Edges, models.GraphEdge{
			Source:     srcStr,
			Target:     tgtStr,
			Label:      labelStr,
			Properties: props,
		})
	}

	return payload, nil
}
