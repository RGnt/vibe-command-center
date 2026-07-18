package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"todo-backend/database"
	"todo-backend/models"
	"todo-backend/testutils"

	"github.com/go-chi/chi/v5"
)

func setupIconRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/icons", testutils.AuthContext(1, iconHandler.GetIcons))
	r.Put("/api/icons/{id}", testutils.AuthContext(1, iconHandler.UpdateIcon))
	r.Delete("/api/icons/{id}", testutils.AuthContext(1, iconHandler.DeleteIcon))
	return r
}

func TestGetIcons(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupIconRouter()

	database.DB.Exec(`INSERT INTO icons (user_id, name, url, folder) VALUES (1, 'Test Icon', '/url', 'General')`)

	req, _ := http.NewRequest("GET", "/api/icons", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var icons []models.Icon
	json.NewDecoder(rr.Body).Decode(&icons)
	if len(icons) != 1 {
		t.Errorf("expected 1 icon, got %v", len(icons))
	}
}

func TestUpdateIcon(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupIconRouter()

	var id int
	database.DB.QueryRow(`INSERT INTO icons (user_id, name, url, folder) VALUES (1, 'Test Icon', '/url', 'General') RETURNING id`).Scan(&id)

	icon := models.Icon{
		Name:   "Updated Icon",
		Folder: "New Folder",
	}
	body, _ := json.Marshal(icon)
	req, _ := http.NewRequest("PUT", "/api/icons/"+strconv.Itoa(id), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseIcon models.Icon
	json.NewDecoder(rr.Body).Decode(&responseIcon)
	if responseIcon.Name != "Updated Icon" {
		t.Errorf("expected name to be 'Updated Icon', got %v", responseIcon.Name)
	}
	if responseIcon.Folder != "New Folder" {
		t.Errorf("expected folder to be 'New Folder', got %v", responseIcon.Folder)
	}
}

func TestDeleteIcon(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupIconRouter()

	var id int
	database.DB.QueryRow(`INSERT INTO icons (user_id, name, url, folder) VALUES (1, 'Test Icon', '/url', 'General') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("DELETE", "/api/icons/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	var count int
	database.DB.QueryRow("SELECT COUNT(*) FROM icons WHERE id = $1", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 icons, got %v", count)
	}
}
