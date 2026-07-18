package handlers

import (
	"encoding/json"
	"net/http"
	"todo-backend/database"
	"todo-backend/middleware"
	"github.com/go-chi/chi/v5"
	"log"
)

type Icon struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name"`
	URL       string `json:"url"`
	Folder    string `json:"folder"`
	CreatedAt string `json:"created_at"`
}

func GetIcons(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	database.Mutex.RLock()
	rows, err := database.DB.Query(`SELECT id, user_id, name, url, folder, created_at FROM icons WHERE user_id = $1 ORDER BY folder ASC, created_at DESC`, userID)
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
		if err := rows.Scan(&i.ID, &i.UserID, &i.Name, &i.URL, &i.Folder, &i.CreatedAt); err != nil {
			continue
		}
		icons = append(icons, i)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(icons)
}

func UpdateIcon(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

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
		_, err = database.DB.Exec(`UPDATE icons SET name = $1, folder = $2 WHERE id = $3 AND user_id = $4`, req.Name, req.Folder, id, userID)
	} else if req.Name != "" {
		_, err = database.DB.Exec(`UPDATE icons SET name = $1 WHERE id = $2 AND user_id = $3`, req.Name, id, userID)
	} else if req.Folder != "" {
		_, err = database.DB.Exec(`UPDATE icons SET folder = $1 WHERE id = $2 AND user_id = $3`, req.Folder, id, userID)
	}
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to update icon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func DeleteIcon(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id := chi.URLParam(r, "id")

	database.Mutex.Lock()
	_, err := database.DB.Exec(`DELETE FROM icons WHERE id = $1 AND user_id = $2`, id, userID)
	database.Mutex.Unlock()

	if err != nil {
		http.Error(w, "Failed to delete icon", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
