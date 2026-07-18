package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/database"
	"github.com/go-chi/chi/v5"
	"log"
)

type Icon struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	Folder    string `json:"folder"`
	CreatedAt string `json:"created_at"`
}

func GetIcons(w http.ResponseWriter, r *http.Request) {
	database.Mutex.RLock()
	rows, err := database.DB.Query(`SELECT id, name, url, folder, created_at FROM icons ORDER BY folder ASC, created_at DESC`)
	database.Mutex.RUnlock()

	if err != nil {
		log.Println("Error fetching icons:", err)
		http.Error(w, "Error fetching icons", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var icons []Icon
	for rows.Next() {
		var i Icon
		if err := rows.Scan(&i.ID, &i.Name, &i.URL, &i.Folder, &i.CreatedAt); err != nil {
			continue
		}
		icons = append(icons, i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(icons)
}

func UpdateIcon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req struct {
		Name   string `json:"name"`
		Folder string `json:"folder"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	var err error
	// Update fields selectively depending on what was provided
	if req.Name != "" && req.Folder != "" {
		_, err = database.DB.Exec(`UPDATE icons SET name = $1, folder = $2 WHERE id = $3`, req.Name, req.Folder, id)
	} else if req.Name != "" {
		_, err = database.DB.Exec(`UPDATE icons SET name = $1 WHERE id = $2`, req.Name, id)
	} else if req.Folder != "" {
		_, err = database.DB.Exec(`UPDATE icons SET folder = $1 WHERE id = $2`, req.Folder, id)
	}
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to update icon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteIcon(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	database.Mutex.Lock()
	_, err := database.DB.Exec(`DELETE FROM icons WHERE id = $1`, id)
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to delete icon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
