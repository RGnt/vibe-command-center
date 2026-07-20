package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"todo-backend/internal/database"
	"todo-backend/internal/models"
	"todo-backend/internal/testutils"

	"github.com/go-chi/chi/v5"
)

func setupWikiRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/api/wikis", testutils.AuthContext(1, wikiHandler.GetWikis))
	r.Get("/api/wikis/{slug}", testutils.AuthContext(1, wikiHandler.GetWiki))
	r.Post("/api/wikis", testutils.AuthContext(1, wikiHandler.CreateWiki))
	r.Put("/api/wikis/{id}", testutils.AuthContext(1, wikiHandler.UpdateWiki))
	r.Delete("/api/wikis/{id}", testutils.AuthContext(1, wikiHandler.DeleteWiki))
	return r
}

// TestCreateWiki ...
func TestCreateWiki(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWikiRouter()

	wiki := models.WikiPage{
		Title:   "Test Wiki",
		Content: "Wiki Content",
	}
	body, _ := json.Marshal(wiki)
	req, _ := http.NewRequest("POST", "/api/wikis", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusCreated {
		t.Errorf("Expected status code %d, got %d. Body: %s", http.StatusCreated, rr.Code, rr.Body.String())
	}

	var responseWiki models.WikiPage
	_ = json.NewDecoder(rr.Body).Decode(&responseWiki)
	if responseWiki.Title != "Test Wiki" {
		t.Errorf("expected title to be 'Test Wiki', got %v", responseWiki.Title)
	}
	if responseWiki.Slug != "test-wiki" {
		t.Errorf("expected slug to be 'test-wiki', got %v", responseWiki.Slug)
	}
}

// TestGetWikis ...
func TestGetWikis(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWikiRouter()

	_, _ = database.DB.Exec(`INSERT INTO wiki_pages (user_id, title, slug, content) VALUES (1, 'Test Wiki', 'test-wiki', 'Content')`)

	req, _ := http.NewRequest("GET", "/api/wikis", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var wikis []models.WikiPage
	_ = json.NewDecoder(rr.Body).Decode(&wikis)
	if len(wikis) != 1 {
		t.Errorf("expected 1 wiki, got %v", len(wikis))
	}
}

// TestGetWiki ...
func TestGetWiki(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWikiRouter()

	_, _ = database.DB.Exec(`INSERT INTO wiki_pages (user_id, title, slug, content) VALUES (1, 'Test Wiki', 'test-wiki', 'Content')`)

	req, _ := http.NewRequest("GET", "/api/wikis/test-wiki", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var w models.WikiPage
	_ = json.NewDecoder(rr.Body).Decode(&w)
	if w.Slug != "test-wiki" {
		t.Errorf("expected slug test-wiki, got %v", w.Slug)
	}
}

// TestUpdateWiki ...
func TestUpdateWiki(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWikiRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO wiki_pages (user_id, title, slug, content) VALUES (1, 'Test Wiki', 'test-wiki', 'Content') RETURNING id`).Scan(&id)

	wiki := models.WikiPage{
		Title:   "Updated Wiki",
		Content: "Updated Content",
		Slug:    "updated-wiki",
	}
	body, _ := json.Marshal(wiki)
	req, _ := http.NewRequest("PUT", "/api/wikis/"+strconv.Itoa(id), bytes.NewBuffer(body))
	
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, rr.Code)
	}

	var responseWiki models.WikiPage
	_ = json.NewDecoder(rr.Body).Decode(&responseWiki)
	if responseWiki.Title != "Updated Wiki" {
		t.Errorf("expected title to be 'Updated Wiki', got %v", responseWiki.Title)
	}
}

// TestDeleteWiki ...
func TestDeleteWiki(t *testing.T) {
	testutils.ClearDB()
	setupTestUser()
	r := setupWikiRouter()

	var id int
	_ = database.DB.QueryRow(`INSERT INTO wiki_pages (user_id, title, slug, content) VALUES (1, 'Test Wiki', 'test-wiki', 'Content') RETURNING id`).Scan(&id)

	req, _ := http.NewRequest("DELETE", "/api/wikis/"+strconv.Itoa(id), nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Errorf("Expected status code %d, got %d", http.StatusNoContent, rr.Code)
	}

	var count int
	_ = database.DB.QueryRow("SELECT COUNT(*) FROM wiki_pages WHERE id = $1", id).Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 wikis, got %v", count)
	}
}
