package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo-backend/database"
	"todo-backend/models"
	"todo-backend/testutils"

	"github.com/go-chi/chi/v5"
)

func setupUserRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/user/settings", testutils.AuthContext(1, userHandler.GetUserSettings))
	r.Put("/api/user/settings", testutils.AuthContext(1, userHandler.UpdateUserSettings))
	return r
}

// TestGetUserSettings ...
func TestGetUserSettings(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupUserRouter()

	_, _ = database.DB.Exec(`INSERT INTO user_settings (user_id, theme, default_project_id) VALUES (1, 'dark', NULL)`)

	req, _ := http.NewRequest("GET", "/api/user/settings", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var settings models.UserSettings
	_ = json.NewDecoder(rr.Body).Decode(&settings)
	if settings.Theme != "dark" {
		t.Errorf("expected theme dark, got %v", settings.Theme)
	}
}

// TestUpdateUserSettings ...
func TestUpdateUserSettings(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupUserRouter()

	_, _ = database.DB.Exec(`INSERT INTO user_settings (user_id, theme, default_project_id) VALUES (1, 'dark', NULL)`)

	newSettings := models.UserSettings{
		Theme: "light",
	}
	body, _ := json.Marshal(newSettings)
	req, _ := http.NewRequest("PUT", "/api/user/settings", bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseSettings models.UserSettings
	_ = json.NewDecoder(rr.Body).Decode(&responseSettings)
	if responseSettings.Theme != "light" {
		t.Errorf("expected theme to be 'light', got %v", responseSettings.Theme)
	}
}
