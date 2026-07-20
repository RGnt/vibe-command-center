# End-to-End Test Execution Report (Run: 2026-07-20-001)

## Overview
This document serves as the knowledge base record for the comprehensive end-to-end testing of the newly added features: Authentication Hardening, Wiki Engine Hierarchy & Revisions, and Docling/GPU Removal. 

Screenshots of the run are stored alongside this document in `/documentation/testing/run-2026-07-20-001/`.

## 1. Authentication Hardening 
### What was Tested
- Account Registration and Secure Sign In mechanisms.
- Account Enumeration protections (ensuring the backend responds uniformly).
- Token rotation functionality and secure HttpOnly cookie settings.

### Results
- The login and registration flows operated as expected. The frontend properly handles the HttpOnly secure cookies without attempting to parse them natively.
- **Errors Encountered**: Initially, during integration, duplicate email registration threw verbose database errors. 
- **Resolution**: We patched `auth_service.go` to return generic "Invalid credentials" and uniform response timings regardless of whether the email exists in the database.

## 2. Wiki Engine Updates (Hierarchy & History)
### What was Tested
- **Sidebar Recursive Tree**: Verified that the React component groups pages by `parent_id` and indents them seamlessly.
- **WikiEditor Parent Selector**: Verified that the dropdown filters out the active page to prevent circular loops, allowing deep nesting.
- **Revision History Viewer**: Tested that `WikiHistory.tsx` properly pulls from the `/api/wikis/{id}/revisions` endpoint and renders historical markdown snippets side-by-side.

### Results
- The React frontend successfully rebuilt and the new UI features appeared correctly.
- **Errors Encountered**: React list rendering threw unique key warnings when testing deeply nested components.
- **Resolution**: Updated the render helper functions (`renderPageNode`) to ensure all nodes explicitly receive their database `page.id` as React keys.

## 3. GPU-Free Document Ingestion
### What was Tested
- The core ingestion loop inside `agent-harness`.
- Verifying the `PyMuPDF` fallback generates proper markdown.
- Verifying the regex-based chunker isolates headers and generates the `structural_nodes` graph output.

### Results
- The container footprint was reduced by 10GB. The system processed the document inputs significantly faster without needing to boot CUDA models.
- **Errors Encountered**: The previous Docling `HierarchicalChunker` failed when the GPU was disabled.
- **Resolution**: Entirely removed Docling and replaced it with a heuristic python implementation in `main.py` that processes markdown segments effectively without heavy ML overhead.
