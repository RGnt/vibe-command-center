# Backend Code Review Report

**Date:** 2026-07-20  
**Reviewer:** Code Review (Adversarial)  
**Scope:** `backend/` — All Go source files under `cmd/` and `internal/`  
**Approach:** Assume buggy until proven otherwise. Flag every concern regardless of severity.

---

## Summary

The backend is a Chi-based REST API for a personal productivity app with a PostgreSQL + Apache AGE (graph) back-end. The code is generally readable, but a systematic review reveals a number of **critical security flaws**, **concurrency bugs**, **architectural violations**, **missing input validation**, and **silent failure paths** that would cause real harm in a production deployment.

Issues are classified as:
- 🔴 **Critical** — Can cause data loss, security breach, or server crash
- 🟠 **High** — Incorrect behavior under reasonable inputs, potential data corruption
- 🟡 **Medium** — Violations of best practices that degrade reliability or maintainability
- 🔵 **Low** — Stylistic concerns or minor improvements

---

## 1. Security Issues

### 🔴 SEC-01 — Hardcoded JWT Secret Falls Back to Production Code
**File:** `cmd/server/main.go:28`

```go
jwtSecret = "my_super_secret_key_change_me_in_prod" // Default for development
```

If `JWT_SECRET` is not set in the environment, the server silently uses a weak, publicly-known secret string instead of aborting startup. An attacker who knows this string can forge JWT tokens for any user ID and gain full access to the API.

**Fix:** Replace with:
```go
if jwtSecret == "" {
    log.Fatal("JWT_SECRET environment variable is required")
}
```

---

### 🔴 SEC-02 — Cookie `Secure` Flag Hardcoded to `false`
**File:** `internal/handlers/auth_handlers.go:64,97`

```go
Secure: false, // Set to true in production if using HTTPS
```

Both `Register` and `Login` handlers set the auth cookie with `Secure: false`, which allows the token to be transmitted over plain HTTP and is vulnerable to cookie theft via network interception. This is a configuration decision that must never be left as a code constant.

**Fix:** Read this from an environment variable (`IS_HTTPS=true`) and apply it dynamically.

---

### 🔴 SEC-03 — CORS Allows All Origins (`https://*`, `http://*`)
**File:** `cmd/server/main.go:82`

```go
AllowedOrigins: []string{"https://*", "http://*"},
```

This effectively disables CORS protection by allowing any origin. Combined with `AllowCredentials: true`, this is a textbook CORS misconfiguration that enables Cross-Site Request Forgery (CSRF) attacks from any website on the internet.

**Fix:** Set `AllowedOrigins` to an explicit list of known frontend origins, configured via an environment variable.

---

### 🔴 SEC-04 — No Input Validation on User-Supplied Strings
**Files:** `internal/handlers/auth_handlers.go`, `todo_handlers.go`, `wiki_handlers.go`, etc.

`email`, `password`, `title`, `content` and all other user-supplied string fields are passed directly to the database or returned in responses without any sanity checks. There is:
- No maximum length validation (can cause DB column overflow errors or DoS via oversized payloads)
- No email format validation on `Register` (any string is accepted as a valid email)
- No empty-string validation on required fields (a `title: ""` todo can be created)

---

### 🟠 SEC-05 — SQL Injection Vector via Dynamic Label Interpolation in Cypher Query
**File:** `internal/repository/graph_repository.go:72-78`

```go
query := fmt.Sprintf(`
    SELECT * FROM cypher('knowledge_graph', $$
        MERGE (n:%s {id: $id})
        ...
    $$, $1) as (v agtype);
`, label)
```

While `sanitizeLabel` is called before this, the fallback value is `"Entity"` for any label that fails the alphanumeric check. The issue is that the compiled regex only validates the label, but the `label` value originates from the untrusted `models.GraphPayload.Nodes[*].Label` field provided by the agent-harness over HTTP. If the agent harness itself is compromised or the code path is changed and `sanitizeLabel` is bypassed, this becomes a direct Cypher injection vulnerability. The current defense is one function call away from being removed.

---

### 🟡 SEC-06 — File Paths Stored in Database are Not Validated for Path Traversal
**File:** `internal/handlers/library_handlers.go:64, 297`

```go
uploadDir := filepath.Join("storage", "library")
```

`filepath.Join` is used safely here, but the `doc.Filepath` value retrieved from the database is passed directly to `os.Open` in `DownloadDocument`. If an attacker could manipulate the database record (e.g., via a future SQL injection), they could cause arbitrary file reads from the host filesystem.

---

## 2. Concurrency and State Bugs

### 🔴 BUG-01 — Global Database Mutex Not Used at Query Sites
**File:** `internal/database/db.go:16`

```go
var Mutex sync.RWMutex
```

A global `sync.RWMutex` is declared and exported on the `DB` package, implying intent to synchronize access. However, examination of every handler and repository shows that **no repository or handler actually acquires this mutex before querying the database**. The only place it is used is in `testutils/testutils.go` for TRUNCATE operations during tests. The `database/sql.DB` pool is concurrency-safe by itself, so the Mutex is misleading. Either it should be removed (to reduce confusion), or its intended usage pattern must be documented and applied consistently.

---

### 🟠 BUG-02 — `GetTodosByStage` Directly Accesses Global `database.DB`
**File:** `internal/handlers/todo_handlers.go:241`

```go
stageRows, err := database.DB.Query(`SELECT name FROM workflow_stages ...`)
```

The `TodoHandler` directly imports and queries `database.DB` rather than going through its service or repository. This:
1. Bypasses the dependency injection pattern used everywhere else.
2. Makes the handler impossible to unit-test in isolation (it requires a real database).
3. Is an architectural layer violation (handler → database, skipping service/repository).
4. Has a self-admitted comment acknowledging this: *"For now, doing it via DB directly is a slight layer breach but acceptable"* — this is not acceptable in production code.

**Fix:** Inject a `WorkflowService` or `WorkflowRepository` into `TodoHandler` and resolve stages through the proper abstraction.

---

### 🟠 BUG-03 — `ToggleCompleted` is Not Atomic (TOCTOU Race Condition)
**File:** `internal/repository/todo_repository.go:127-136`

```go
err := r.db.QueryRow("SELECT completed FROM todos ...").Scan(&currentCompleted)
...
_, err = r.db.Exec("UPDATE todos SET completed = $1 ...", newCompleted, ...)
```

The SELECT and UPDATE are two separate queries without a transaction. Under concurrent requests (e.g., a user double-clicking rapidly), two goroutines can both read `completed = false` and both set it to `true`, losing the expected toggle behavior. This is a classic TOCTOU (time-of-check-time-of-use) race condition.

**Fix:** Use a single atomic SQL statement:
```sql
UPDATE todos SET completed = NOT completed WHERE id = $1 AND user_id = $2 RETURNING completed
```

---

### 🟠 BUG-04 — `UpsertGraph` Transaction Does Not Rollback on Node/Edge Failure
**File:** `internal/repository/graph_repository.go:80-84`

```go
_, err = tx.Exec(query, string(paramsJSON))
if err != nil {
    log.Printf("Failed to upsert node %s: %v", node.ID, err)
    // Continue with other nodes instead of completely failing
}
```

When a node upsert fails, the error is **only logged and the loop continues**. The transaction is then committed regardless. This means partial, inconsistent graph states can be committed silently. A graph with some nodes but missing their relationships (because node upserts failed) can cause downstream query failures.

---

## 3. Architectural Violations

### 🟠 ARCH-01 — `CreateUser` Contains Complex Business Logic in Repository Layer
**File:** `internal/repository/user_repository.go:32-96`

The `CreateUser` function in the repository layer creates not just a user, but also: a default workflow, four workflow stages, a default project, updates user settings, and creates a default wiki page — all in one massive transaction. This is registration orchestration logic that belongs in the service layer (`auth_service.go`).

As a result:
- The repository is not a pure data-access layer.
- The `UserRepository` interface claims to do one thing (`CreateUser`) but secretly does six.
- It is impossible to reuse or test these operations independently.

---

### 🟡 ARCH-02 — `Claims` Type is Duplicated Between `service` and `middleware` Packages
**File:** `internal/service/auth_service.go:18-21`, `internal/middleware/auth.go:11-14`

The `Claims` struct is defined twice in two separate packages with identical content:

```go
// service/auth_service.go
type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}

// middleware/auth.go
type Claims struct {
    UserID int `json:"user_id"`
    jwt.RegisteredClaims
}
```

This violates DRY (Don't Repeat Yourself). If the JWT claims schema changes (e.g., adding a role field), one definition will be updated while the other silently diverges.

**Fix:** Define `Claims` once in a shared location (e.g., a new `internal/auth` package) and import it in both places.

---

### 🟡 ARCH-03 — Handler Receives Concrete `*sql.DB` via `repository.GraphRepository` Bypass
**File:** `internal/handlers/library_handlers.go:31`

```go
type LibraryHandler struct {
    ...
    graphRepo repository.GraphRepository
}
```

The `LibraryHandler` receives the `GraphRepository` directly, making it responsible for writing graph data after document ingestion. This crosses a domain boundary — a library/document handler should not be responsible for graph persistence. The graph upsert should be triggered from the service layer.

---

## 4. Error Handling Gaps

### 🟠 ERR-01 — `reqBody` Marshal Error Silently Ignored in `IngestDocument`
**File:** `internal/handlers/library_handlers.go:223`

```go
reqBody, _ := json.Marshal(map[string]string{...})
```

The error from `json.Marshal` is discarded with `_`. While marshaling a `map[string]string` to JSON is unlikely to fail, this pattern teaches developers it is acceptable to discard errors from serialization.

---

### 🟠 ERR-02 — `strconv.Atoi` Error Silently Ignored in `GetTodos`
**File:** `internal/handlers/todo_handlers.go:40`

```go
projectID, _ := strconv.Atoi(projectIDStr)
pID = &projectID
```

If `project_id=abc` is passed as a query parameter, `Atoi` silently returns `0`, and the handler will query for `project_id = 0`, which will return an empty result instead of a `400 Bad Request`. The client has no idea the query parameter was malformed.

**Fix:**
```go
projectID, err := strconv.Atoi(projectIDStr)
if err != nil {
    http.Error(w, "Invalid project_id parameter", http.StatusBadRequest)
    return
}
```

---

### 🟡 ERR-03 — Error Detection via String Matching is Fragile
**File:** `internal/handlers/auth_handlers.go:50`

```go
if strings.Contains(err.Error(), "email might already exist") {
```

The handler detects a duplicate email error by parsing the error's string message. The comment acknowledges this: *"Could use strongly typed error here too, exported from repository layer"*. Indeed — `ErrEmailExists` is already exported from `internal/repository/user_repository.go:13` but is wrapped with a different message inside `CreateUser`. The sentinel error is never actually surfaced, making `errors.Is` unusable.

**Fix:** Unwrap the error correctly and use `errors.Is(err, repository.ErrEmailExists)`.

---

### 🟡 ERR-04 — `GetMe` Returns `404 Not Found` for All DB Errors
**File:** `internal/handlers/auth_handlers.go:130`

```go
user, err := h.authService.GetUserByID(userID)
if err != nil {
    http.Error(w, "User not found", http.StatusNotFound)
    return
}
```

Any database error (e.g., connection timeout, deadlock) is incorrectly surfaced as a `404 Not Found` to the client. This makes debugging difficult and gives clients incorrect error semantics. A database timeout should be a `500 Internal Server Error`.

---

## 5. Data Model Issues

### 🟡 DATA-01 — `Diagram.CreatedAt` and `Icon.CreatedAt` are Typed as `string`, Not `time.Time`
**File:** `internal/models/models.go:81-82, 92`

```go
// Diagram
CreatedAt   string `json:"created_at"`
UpdatedAt   string `json:"updated_at"`

// Icon
CreatedAt string `json:"created_at"`
```

These fields are `string` while all other models use `time.Time`. Using `string` for timestamps means:
- No timezone normalization.
- No ability to sort or compare dates in Go code without manual parsing.
- Inconsistent API response format (some dates are ISO 8601 objects, others are raw strings).

---

### 🟡 DATA-02 — `wiki_pages` Table Has No `UNIQUE` Constraint Scoped to Project
**File:** `internal/database/db.go:151`

```sql
UNIQUE(user_id, slug)
```

The unique constraint on wiki slugs is scoped globally to `(user_id, slug)`. This means the same slug cannot appear in two different projects belonging to the same user. This is overly restrictive and prevents re-using natural slug names (like "overview") across projects.

---

### 🟡 DATA-03 — Schema Migrations Handled via Raw `ALTER TABLE` in `createTables()`
**File:** `internal/database/db.go:200-205`

```go
for _, table := range tables {
    _, _ = DB.Exec(`ALTER TABLE ` + table + ` ADD COLUMN IF NOT EXISTS user_id ...`)
}
```

Schema evolution is handled by appending `ADD COLUMN IF NOT EXISTS` statements to the initialization function. This is:
- Not reproducible across environments (adding tables later won't re-run old migrations).
- Not version-controlled (there's no migration history).
- Dangerous in production (running `createTables()` on a live DB modifies schema mid-flight).

**Fix:** Adopt a migration tool such as `golang-migrate/migrate` with versioned SQL migration files.

---

## 6. Missing Functionality

### 🟠 FUNC-01 — `DeleteDocument` Does Not Delete Extracted Images
**File:** `internal/handlers/library_handlers.go:187-193`

```go
if err := h.libraryService.DeleteDocument(id, userID); err != nil { ... }
_ = os.Remove(doc.Filepath)
```

When a document is deleted, only its primary file is removed. Any images extracted during `IngestDocument` and stored as separate `LibraryDocument` records are orphaned in both the database and on disk.

---

### 🟡 FUNC-02 — No Rate Limiting on Auth Endpoints
**File:** `cmd/server/main.go:88-91`

The `Register` and `Login` endpoints have no rate limiting. A brute-force or credential-stuffing attack against `/api/auth/login` is completely unconstrained.

---

### 🟡 FUNC-03 — No Server-Side Request Timeout
**File:** `cmd/server/main.go:161`

```go
log.Fatal(http.ListenAndServe(":8080", r))
```

The HTTP server is started with no read/write timeouts, which means slow clients or misbehaving connections can hold goroutines indefinitely and exhaust the server under moderate load.

**Fix:**
```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      r,
    ReadTimeout:  15 * time.Second,
    WriteTimeout: 30 * time.Second,
    IdleTimeout:  60 * time.Second,
}
log.Fatal(srv.ListenAndServe())
```

---

## Issue Summary Table

| ID | Severity | Category | Location | Title |
|----|----------|----------|----------|-------|
| SEC-01 | 🔴 Critical | Security | `main.go:28` | Hardcoded JWT secret |
| SEC-02 | 🔴 Critical | Security | `auth_handlers.go:64` | Cookie `Secure: false` |
| SEC-03 | 🔴 Critical | Security | `main.go:82` | Wildcard CORS with credentials |
| SEC-04 | 🔴 Critical | Security | Multiple handlers | No input validation |
| SEC-05 | 🟠 High | Security | `graph_repository.go:72` | Cypher label interpolation |
| SEC-06 | 🟡 Medium | Security | `library_handlers.go:152` | File path from DB to `os.Open` |
| BUG-01 | 🔴 Critical | Concurrency | `database/db.go:16` | Unused exported Mutex |
| BUG-02 | 🟠 High | Architecture | `todo_handlers.go:241` | Handler bypasses DI, queries DB directly |
| BUG-03 | 🟠 High | Concurrency | `todo_repository.go:127` | Non-atomic toggle (TOCTOU race) |
| BUG-04 | 🟠 High | Correctness | `graph_repository.go:80` | Partial graph commit on failure |
| ARCH-01 | 🟠 High | Architecture | `user_repository.go:32` | Business logic in repository layer |
| ARCH-02 | 🟡 Medium | Architecture | `service/` and `middleware/` | Duplicate `Claims` struct |
| ARCH-03 | 🟡 Medium | Architecture | `library_handlers.go:31` | Handler crosses domain boundary |
| ERR-01 | 🟠 High | Error Handling | `library_handlers.go:223` | Silenced marshal error |
| ERR-02 | 🟠 High | Error Handling | `todo_handlers.go:40` | Silenced `Atoi` error |
| ERR-03 | 🟡 Medium | Error Handling | `auth_handlers.go:50` | String-based error detection |
| ERR-04 | 🟡 Medium | Error Handling | `auth_handlers.go:130` | DB errors reported as 404 |
| DATA-01 | 🟡 Medium | Data Model | `models/models.go:81` | Timestamps typed as `string` |
| DATA-02 | 🟡 Medium | Data Model | `database/db.go:151` | Overly-scoped slug uniqueness |
| DATA-03 | 🟡 Medium | Data Model | `database/db.go:200` | No migration framework |
| FUNC-01 | 🟠 High | Correctness | `library_handlers.go:187` | Orphaned images on delete |
| FUNC-02 | 🟡 Medium | Security | `main.go:88` | No auth rate limiting |
| FUNC-03 | 🟡 Medium | Reliability | `main.go:161` | No HTTP server timeouts |

---

## Recommended Priority Order

1. **Immediately:** SEC-01 (hardcoded JWT secret), SEC-02 (insecure cookie), SEC-03 (wildcard CORS)
2. **Before next release:** BUG-03 (TOCTOU race), ERR-02 (silent Atoi), FUNC-01 (image orphan on delete), BUG-02 (handler-DB bypass)
3. **Backlog:** ARCH-01 (repository bloat), DATA-03 (migrations), FUNC-02 (rate limiting), FUNC-03 (server timeouts)
