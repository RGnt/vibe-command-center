# ADR 005: JWT (JSON Web Tokens) for Authentication

## Status
Accepted

## Context
We need a method to authenticate users, restrict API access, and uniquely identify users across API requests without maintaining expensive server-side session state.

## Decision
We chose **Standard JWT (JSON Web Tokens)** managed manually via the `golang-jwt/jwt` library.

## Rationale
- **Stateless**: Tokens carry all necessary context (`user_id`, `exp`) cryptographically signed via `HS256`. The server does not need to perform a database lookup merely to verify session validity.
- **Decoupled**: Works perfectly in an architecture where the frontend (React) and backend (Go API) are strictly separated. The React app simply stores the token and attaches it to the `Authorization: Bearer` header.
- **Simplicity**: Building it natively allowed us to bypass the complexity or external dependencies of 3rd party providers (e.g. Auth0, Better-auth) while keeping the database schema fully under our control.

## Consequences
- JWTs cannot easily be invalidated server-side before their expiration time without implementing a token blocklist. We accepted this trade-off by keeping expiration times reasonable (1 week) and prioritizing architectural simplicity.
- The React application is responsible for securely storing the token (currently in memory / localStorage).
