# Wiki Improvements Specification

## Goal
Implement missing "usual Wiki" features: a formal hierarchy (folders/parent-child) and document history (revisions).

## Selected Approach
We will use **Option 1 (True Parent-Child Hierarchy)** for its robust handling of renames and structural changes, alongside a dedicated table for tracking historical edits.

## 1. Database Schema (`backend/internal/database/db.go`)
- **Hierarchy**: Add `parent_id INTEGER REFERENCES wiki_pages(id) ON DELETE CASCADE` to the `wiki_pages` table.
- **Revisions**: Create a `wiki_page_revisions` table:
  ```sql
  CREATE TABLE IF NOT EXISTS wiki_page_revisions (
      id SERIAL PRIMARY KEY,
      wiki_page_id INTEGER REFERENCES wiki_pages(id) ON DELETE CASCADE,
      content TEXT NOT NULL,
      created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
  )
  ```

## 2. Models (`backend/internal/models/models.go`)
- Add `ParentID *int \`json:"parent_id,omitempty"\`` to `WikiPage`.
- Create `WikiPageRevision` struct.

## 3. Repository Layer (`backend/internal/repository/wiki_repository.go`)
- Update `Create`, `Update`, `GetAllByUserID`, and `GetBySlugAndUserID` to scan and insert `parent_id`.
- Update the `Update` method to automatically insert the *old* content into `wiki_page_revisions` before modifying the page.
- Add `GetRevisions(wikiID int)` to fetch history for a given page.

## 4. Service Layer (`backend/internal/service/wiki_service.go`)
- Add `GetWikiRevisions(wikiID int, userID int)` method.
- Add cycle detection logic when assigning a `ParentID` to prevent infinite loops in the tree.

## 5. Handlers & Routing
- Update JSON decoders in `CreateWiki` and `UpdateWiki` to accept `parent_id`.
- Add `GetWikiRevisions` handler.
- Register `r.Get("/api/wikis/{id}/revisions", wikiHandler.GetWikiRevisions)` in `main.go`.
