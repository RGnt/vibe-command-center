# ADR 002: PostgreSQL and database/sql

## Status
Accepted

## Context
We need a robust, relational database to guarantee ACID compliance for our project structuring and kanban board logic.
Additionally, we need a way to interact with the database from our Go backend.

## Decision
We chose **PostgreSQL** as our primary datastore. We also configured the Docker container to build with `pgvector` and Apache `AGE` to future-proof the application for semantic search and graph relationships.
For Go integration, we chose the standard library **database/sql** using the `pq` driver, eschewing heavy ORMs like GORM.

## Rationale
- **PostgreSQL**: Industry standard, incredibly reliable, and handles JSON columns which is useful for storing arbitrary settings or diagram payloads.
- **Extensions**: Adding vector and graph capabilities immediately prevents needing secondary databases if the application scale requires AI search or graph-traversal of linked notes.
- **Raw SQL (`database/sql`)**: 
  - ORMs add cognitive overhead and obfuscate query performance.
  - By writing raw SQL within Repository structs, we retain 100% control over indexing, query execution paths, and transactional locking.

## Consequences
- Requires manual schema management and structural querying.
- Scanning query results into structs requires explicit boilerplate.
