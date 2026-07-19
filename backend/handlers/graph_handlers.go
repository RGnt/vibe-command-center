package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/repository"
)

type GraphHandler struct {
	graphRepo repository.GraphRepository
}

func NewGraphHandler(graphRepo repository.GraphRepository) *GraphHandler {
	return &GraphHandler{
		graphRepo: graphRepo,
	}
}

func (h *GraphHandler) GetGraph(w http.ResponseWriter, r *http.Request) {
	graphData, err := h.graphRepo.GetGraph()
	if err != nil {
		http.Error(w, "Failed to fetch graph data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(graphData)
}
