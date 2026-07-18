package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"todo-backend/database"
	"todo-backend/middleware"

	"github.com/google/uuid"
)

// UploadFile handles multipart form uploads and saves files to ./uploads
func UploadFile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 10 MB limit
	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Validate extension
	ext := strings.ToLower(filepath.Ext(handler.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" && ext != ".webp" && ext != ".svg" {
		http.Error(w, "Invalid file type. Only images (including SVG) are allowed.", http.StatusBadRequest)
		return
	}

	// Ensure uploads directory exists
	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create uploads directory", http.StatusInternalServerError)
		return
	}

	// Create a unique filename
	newFilename := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, newFilename)

	dest, err := os.Create(filePath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}
	defer dest.Close()

	if _, err := io.Copy(dest, file); err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	url := "/uploads/" + newFilename
	name := strings.TrimSuffix(handler.Filename, ext)

	// Save to database
	database.Mutex.Lock()
	var id int
	err = database.DB.QueryRow(`INSERT INTO icons (user_id, name, url) VALUES ($1, $2, $3) RETURNING id`, userID, name, url).Scan(&id)
	database.Mutex.Unlock()

	if err != nil {
		log.Println("Failed to insert icon to DB:", err)
		http.Error(w, "Failed to save icon data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
		"url": url,
		"name": name,
	})
}
