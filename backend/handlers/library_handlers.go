package handlers

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
	"todo-backend/middleware"
	"todo-backend/models"
	"todo-backend/repository"
	"todo-backend/service"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type LibraryHandler struct {
	libraryService service.LibraryService
	wikiService    service.WikiService
	graphRepo      repository.GraphRepository
}

func NewLibraryHandler(libraryService service.LibraryService, wikiService service.WikiService, graphRepo repository.GraphRepository) *LibraryHandler {
	return &LibraryHandler{
		libraryService: libraryService,
		wikiService:    wikiService,
		graphRepo:      graphRepo,
	}
}

func (h *LibraryHandler) UploadDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// 50 MB limit
	r.ParseMultipartForm(50 << 20)

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to retrieve file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(handler.Filename))

	// Ensure secure library directory exists
	uploadDir := filepath.Join("storage", "library")
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		http.Error(w, "Failed to create library directory", http.StatusInternalServerError)
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

	size, err := io.Copy(dest, file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	mimeType := handler.Header.Get("Content-Type")
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	doc := models.LibraryDocument{
		UserID:       userID,
		OriginalName: handler.Filename,
		Filename:     newFilename,
		Filepath:     filePath,
		Size:         size,
		MimeType:     mimeType,
	}

	createdDoc, err := h.libraryService.CreateDocument(doc)
	if err != nil {
		log.Println("Failed to insert document to DB:", err)
		http.Error(w, "Failed to save document data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdDoc)
}

func (h *LibraryHandler) ListDocuments(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	docs, err := h.libraryService.GetDocumentsByUserID(userID)
	if err != nil {
		http.Error(w, "Failed to fetch documents", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(docs)
}

func (h *LibraryHandler) DownloadDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.libraryService.GetDocumentByID(id, userID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	file, err := os.Open(doc.Filepath)
	if err != nil {
		http.Error(w, "File not found on disk", http.StatusNotFound)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename=\""+doc.OriginalName+"\"")
	w.Header().Set("Content-Type", doc.MimeType)
	w.Header().Set("Content-Length", strconv.FormatInt(doc.Size, 10))

	io.Copy(w, file)
}

func (h *LibraryHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.libraryService.GetDocumentByID(id, userID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	if err := h.libraryService.DeleteDocument(id, userID); err != nil {
		http.Error(w, "Failed to delete document", http.StatusInternalServerError)
		return
	}

	os.Remove(doc.Filepath)
	w.WriteHeader(http.StatusNoContent)
}

func (h *LibraryHandler) IngestDocument(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid document ID", http.StatusBadRequest)
		return
	}

	doc, err := h.libraryService.GetDocumentByID(id, userID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	// The agent-harness has /app/workspace mounted to the host project root (Gemma4).
	// The backend saves files to /app/storage/library/... inside its container, which is mounted from host ./data/storage/library/...
	// So agent-harness sees the file at /app/workspace/data/storage/library/...
	
	agentFilepath := "data/" + doc.Filepath // 'data/storage/library/filename.pdf'

	reqBody, _ := json.Marshal(map[string]string{
		"filepath": agentFilepath,
		"original_filename": doc.OriginalName,
	})

	agentHarnessURL := os.Getenv("AGENT_HARNESS_URL")
	if agentHarnessURL == "" {
		agentHarnessURL = "http://agent-harness:8000"
	}

	resp, err := http.Post(agentHarnessURL+"/api/v1/convert", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		http.Error(w, "Failed to communicate with agent harness", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		http.Error(w, fmt.Sprintf("Agent harness failed: %s", string(bodyBytes)), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	scanner := bufio.NewScanner(resp.Body)
	// docling returns very large lines for markdown output, especially with base64 images, increase buffer size
	buf := make([]byte, 0, 64*1024)
	scanner.Buffer(buf, 100*1024*1024) // up to 100MB line buffer

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")
			var event map[string]interface{}
			if err := json.Unmarshal([]byte(data), &event); err != nil {
				continue
			}

			status, _ := event["status"].(string)
			if status == "processing" || status == "extracting_graph" {
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
			} else if status == "error" {
				fmt.Fprintf(w, "data: %s\n\n", data)
				flusher.Flush()
				return
			} else if status == "complete" {
				markdown, _ := event["markdown"].(string)

				// Find and extract Base64 embedded images
				re := regexp.MustCompile(`!\[(.*?)\]\(data:image/([a-zA-Z]+);base64,([^)]+)\)`)
				matches := re.FindAllStringSubmatch(markdown, -1)
				
				for _, match := range matches {
					if len(match) == 4 {
						fullMatch := match[0]
						caption := match[1]
						ext := match[2]
						b64Data := match[3]

						imgBytes, err := base64.StdEncoding.DecodeString(b64Data)
						if err != nil {
							continue
						}

						// Save it as a LibraryDocument
						uploadDir := filepath.Join("storage", "library")
						os.MkdirAll(uploadDir, os.ModePerm)
						newFilename := uuid.New().String() + "." + ext
						filePath := filepath.Join(uploadDir, newFilename)

						dest, err := os.Create(filePath)
						if err == nil {
							dest.Write(imgBytes)
							dest.Close()

							docRecord := models.LibraryDocument{
								UserID:       userID,
								OriginalName: caption + "." + ext,
								Filename:     newFilename,
								Filepath:     filePath,
								Size:         int64(len(imgBytes)),
								MimeType:     "image/" + ext,
							}

							createdDoc, err := h.libraryService.CreateDocument(docRecord)
							if err == nil {
								// Replace the fullMatch with the new URL
								newURL := fmt.Sprintf("![%s](/api/library/%d/download)", caption, createdDoc.ID)
								markdown = strings.Replace(markdown, fullMatch, newURL, 1)
							}
						}
					}
				}

				// Extract and save the graph if present
				if graphData, ok := event["graph"]; ok {
					graphBytes, _ := json.Marshal(graphData)
					var graphPayload models.GraphPayload
					if err := json.Unmarshal(graphBytes, &graphPayload); err == nil {
						err := h.graphRepo.UpsertGraph(graphPayload)
						if err != nil {
							log.Println("Error upserting graph:", err)
						}
					}
				}

				// Create Wiki Page
				baseName := strings.TrimSuffix(doc.OriginalName, filepath.Ext(doc.OriginalName))
				slug := strings.ToLower(strings.ReplaceAll(baseName, " ", "-")) + "-" + strconv.FormatInt(time.Now().Unix(), 10)

				newWiki := models.WikiPage{
					Title:    baseName,
					Slug:     slug,
					Content:  markdown,
					Category: "Ingested Documents",
				}

				createdWiki, err := h.wikiService.CreateWiki(userID, newWiki)
				if err != nil {
					errMsg, _ := json.Marshal(map[string]string{"status": "error", "detail": "Failed to create wiki page"})
					fmt.Fprintf(w, "data: %s\n\n", string(errMsg))
					flusher.Flush()
					return
				}

				successMsg, _ := json.Marshal(map[string]interface{}{
					"status": "complete",
					"wiki":   createdWiki,
				})
				fmt.Fprintf(w, "data: %s\n\n", string(successMsg))
				flusher.Flush()
				return
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading agent stream:", err)
		errMsg, _ := json.Marshal(map[string]string{"status": "error", "detail": "Stream reading failed"})
		fmt.Fprintf(w, "data: %s\n\n", string(errMsg))
		flusher.Flush()
	}
}
