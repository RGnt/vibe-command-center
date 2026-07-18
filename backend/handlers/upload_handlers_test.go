package handlers

import (
	"bytes"
	"encoding/json"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"todo-backend/database"
	"todo-backend/testutils"

	"github.com/go-chi/chi/v5"
)

func setupUploadRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/api/upload", testutils.AuthContext(1, uploadHandler.UploadFile))
	return r
}

func TestUploadFile(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupUploadRouter()

	// Clean up uploads directory if it exists
	os.RemoveAll("uploads")

	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "testimage.png")
	part.Write([]byte("fake image data"))
	writer.Close()

	req, _ := http.NewRequest("POST", "/api/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusOK, rr.Code, rr.Body.String())
	}

	var resp map[string]interface{}
	json.NewDecoder(rr.Body).Decode(&resp)

	if resp["url"] == nil {
		t.Errorf("expected url in response")
	}

	// Verify DB entry
	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM icons WHERE user_id = 1").Scan(&count)
	if count != 1 {
		t.Errorf("expected 1 icon in db, got %v", count)
	}

	// Clean up
	os.RemoveAll("uploads")
}
