# Wiki Improvements Execution Report

**Date**: 2026-07-20  
**Scope**: Wiki Feature Expansion (Hierarchy & History)  

## Summary
The user requested that the existing Wiki system be expanded to include "usual Wiki features" such as document history and a formal directory hierarchy capable of supporting structures like `instructions/prompts/<file>.md`.

This report documents what was built and why.

---

## What Was Added & Changed

### 1. Document History & Revisions
* **What**: Added an automatic revision tracking system. Every time a Wiki page is updated, the previous version of its contents is securely preserved in a new database table.
* **Why**: Version control is a core feature of any Wiki. It allows users to view how documents evolved over time, recover accidentally deleted information, and maintain an audit log of changes without requiring Git.
* **How**: 
  * Created the `wiki_page_revisions` PostgreSQL table.
  * Modified the `Update` method in `WikiRepository`. It now uses a database transaction to fetch the existing page content, insert it into `wiki_page_revisions` (if it differs from the incoming content), and then update the active `wiki_pages` row.
  * Added two new API endpoints (`GET /api/wikis/{id}/revisions` and `GET /api/wikis/{id}/revisions/{revID}`) to allow the frontend to request historical versions.

### 2. Formal Page Hierarchy (Folders & Subpages)
* **What**: Added a True Parent-Child tree structure to the Wiki. 
* **Why**: A flat Wiki quickly becomes unmanageable as the document count grows. A hierarchy allows pages to act as "folders" or parent documents (e.g., an "Architecture" parent page containing multiple "ADR" child pages), facilitating intuitive UI rendering (like a collapsible file tree) and logical organization.
* **How**:
  * Added a `parent_id` foreign key column to the `wiki_pages` table, referencing another `wiki_pages(id)`. This allows infinite nesting depths.
  * Exposed `ParentID` through the `models.WikiPage` struct and JSON payloads so the frontend can send/receive the tree structure.
  * Implemented strict **Cycle Detection** in `WikiService.UpdateWiki`. Before a page's `parent_id` is saved, the backend walks up the ancestry tree to ensure the new parent isn't actually a descendant of the current page, preventing infinite circular loops.

### 3. Path-Based URL Resolution
* **What**: Updated the routing engine to natively support slashes in Wiki identifiers.
* **Why**: The user specifically requested paths like `instruction/prompts/file.md`. Standard REST routers treat `/` as a delimiter, breaking if a parameter contains a slash.
* **How**:
  * Updated the Chi router definition in `main.go` from `r.Get("/api/wikis/{slug}", ...)` to `r.Get("/api/wikis/{slug:*}", ...)`.
  * This catch-all route means a request to `/api/wikis/instruction/prompts/my-file.md` correctly maps the entire `instruction/prompts/my-file.md` string into the `slug` parameter, allowing the database to look it up flawlessly.

## Status
All tasks are completed and compiling cleanly. The backend is now fully capable of powering a modern, hierarchical Wiki application with document history.
