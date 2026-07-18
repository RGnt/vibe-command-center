package handlers

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"todo-backend/database"
	"todo-backend/middleware"
	"todo-backend/models"

	"github.com/go-chi/chi/v5"
)

// GetWikis fetches all wiki pages for a given project, or global wikis if project_id is omitted
func GetWikis(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	projectIDStr := r.URL.Query().Get("project_id")
	
	var rows *sql.Rows
	var err error

	if projectIDStr != "" {
		projectID, _ := strconv.Atoi(projectIDStr)
		rows, err = database.DB.Query(`SELECT id, project_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE project_id = $1 AND user_id = $2 ORDER BY category, title`, projectID, userID)
	} else {
		rows, err = database.DB.Query(`SELECT id, project_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE project_id IS NULL AND user_id = $1 ORDER BY category, title`, userID)
	}

	if err != nil {
		http.Error(w, "Failed to fetch wikis", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var pages []models.WikiPage
	for rows.Next() {
		var p models.WikiPage
		var projectID sql.NullInt64
		var createdAt, updatedAt time.Time

		err := rows.Scan(&p.ID, &projectID, &p.Category, &p.Title, &p.Slug, &p.Content, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, "Failed to scan wiki page", http.StatusInternalServerError)
			return
		}

		p.UserID = userID
		p.CreatedAt = createdAt
		p.UpdatedAt = updatedAt
		
		if projectID.Valid {
			pid := int(projectID.Int64)
			p.ProjectID = &pid
		}

		pages = append(pages, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pages)
}

// GetWiki gets a specific wiki page by slug
func GetWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	slug := chi.URLParam(r, "slug")

	database.Mutex.RLock()
	defer database.Mutex.RUnlock()

	var p models.WikiPage
	var projectID sql.NullInt64
	var createdAt, updatedAt time.Time

	err := database.DB.QueryRow(`SELECT id, project_id, category, title, slug, content, created_at, updated_at FROM wiki_pages WHERE slug = $1 AND user_id = $2`, slug, userID).
		Scan(&p.ID, &projectID, &p.Category, &p.Title, &p.Slug, &p.Content, &createdAt, &updatedAt)
		
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Wiki page not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch wiki page", http.StatusInternalServerError)
		}
		return
	}

	p.UserID = userID
	p.CreatedAt = createdAt
	p.UpdatedAt = updatedAt

	if projectID.Valid {
		pid := int(projectID.Int64)
		p.ProjectID = &pid
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// CreateWiki creates a new wiki page
func CreateWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var p models.WikiPage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	var id int64
	query := `INSERT INTO wiki_pages (user_id, project_id, category, title, slug, content) VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := database.DB.QueryRow(query, userID, p.ProjectID, p.Category, p.Title, p.Slug, p.Content).Scan(&id)
	if err != nil {
		http.Error(w, "Failed to create wiki page (slug might not be unique)", http.StatusInternalServerError)
		return
	}

	p.ID = int(id)
	p.UserID = userID
	p.CreatedAt = time.Now()
	p.UpdatedAt = p.CreatedAt

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(p)
}

// UpdateWiki updates a wiki page
func UpdateWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid wiki ID", http.StatusBadRequest)
		return
	}

	var p models.WikiPage
	if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	query := `UPDATE wiki_pages SET category = $1, title = $2, slug = $3, content = $4, updated_at = CURRENT_TIMESTAMP WHERE id = $5 AND user_id = $6`
	_, err = database.DB.Exec(query, p.Category, p.Title, p.Slug, p.Content, id, userID)
	if err != nil {
		http.Error(w, "Failed to update wiki page", http.StatusInternalServerError)
		return
	}

	p.ID = id
	p.UserID = userID
	p.UpdatedAt = time.Now()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(p)
}

// DeleteWiki deletes a wiki page
func DeleteWiki(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	id, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		http.Error(w, "Invalid wiki ID", http.StatusBadRequest)
		return
	}

	database.Mutex.Lock()
	defer database.Mutex.Unlock()

	_, err = database.DB.Exec("DELETE FROM wiki_pages WHERE id = $1 AND user_id = $2", id, userID)
	if err != nil {
		http.Error(w, "Failed to delete wiki page", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
