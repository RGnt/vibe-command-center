package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/database"
	"github.com/go-chi/chi/v5"
	"log"
)

type Diagram struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	DiagramType  string `json:"diagram_type"`
	Code         string `json:"code"`
	Explanation  string `json:"explanation"`
	CreatedAt    string `json:"created_at"`
	UpdatedAt    string `json:"updated_at"`
}

func GetDiagrams(w http.ResponseWriter, r *http.Request) {
	database.Mutex.RLock()
	rows, err := database.DB.Query(`SELECT id, name, diagram_type, code, explanation, created_at, updated_at FROM diagrams ORDER BY updated_at DESC`)
	database.Mutex.RUnlock()

	if err != nil {
		log.Println("Error fetching diagrams:", err)
		http.Error(w, "Error fetching diagrams", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var diagrams []Diagram
	for rows.Next() {
		var d Diagram
		if err := rows.Scan(&d.ID, &d.Name, &d.DiagramType, &d.Code, &d.Explanation, &d.CreatedAt, &d.UpdatedAt); err != nil {
			continue
		}
		diagrams = append(diagrams, d)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(diagrams)
}

func GetDiagram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	database.Mutex.RLock()
	var d Diagram
	err := database.DB.QueryRow(`SELECT id, name, diagram_type, code, explanation, created_at, updated_at FROM diagrams WHERE id = $1`, id).
		Scan(&d.ID, &d.Name, &d.DiagramType, &d.Code, &d.Explanation, &d.CreatedAt, &d.UpdatedAt)
	database.Mutex.RUnlock()

	if err != nil {
		http.Error(w, "Diagram not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(d)
}

func CreateDiagram(w http.ResponseWriter, r *http.Request) {
	var req Diagram
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	err := database.DB.QueryRow(
		`INSERT INTO diagrams (name, diagram_type, code, explanation) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at`,
		req.Name, req.DiagramType, req.Code, req.Explanation).
		Scan(&req.ID, &req.CreatedAt, &req.UpdatedAt)
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to create diagram", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func UpdateDiagram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req Diagram
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	err := database.DB.QueryRow(
		`UPDATE diagrams SET name = $1, diagram_type = $2, code = $3, explanation = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 RETURNING updated_at`,
		req.Name, req.DiagramType, req.Code, req.Explanation, id).
		Scan(&req.UpdatedAt)
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to update diagram", http.StatusInternalServerError)
		return
	}
	
	req.ID = 0 // not strictly necessary, but good to know
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(req)
}

func DeleteDiagram(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	database.Mutex.Lock()
	_, err := database.DB.Exec(`DELETE FROM diagrams WHERE id = $1`, id)
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to delete diagram", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
