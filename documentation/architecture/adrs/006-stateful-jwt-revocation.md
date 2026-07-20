# 006. Stateful JWT Revocation via Blocklist

**Date:** 2026-07-20
**Status:** Accepted

## Context
Our application utilizes JSON Web Tokens (JWTs) for stateless authentication. While JWTs are excellent for horizontal scaling (as the server does not need to look up session state on every request), their stateless nature means they cannot be invalidated before they naturally expire. 

Previously, when a user clicked "Logout", the client simply deleted the cookie. However, if that token was intercepted or stolen, an attacker could continue using it until its 7-day expiration elapsed. This violates basic security requirements for session termination.

## Decision
We decided to implement a **Stateful Server-Side Blocklist** for revoked JWTs.

1. **JTI Injection**: Every JWT is now minted with a unique `jti` (JWT ID) UUID claim.
2. **Database Table**: A new `revoked_tokens` table was created in PostgreSQL, storing the `jti` and its `expires_at` timestamp.
3. **Logout Endpoint**: The `/api/auth/logout` endpoint parses the user's incoming token, extracts the `jti`, and inserts it into the `revoked_tokens` table.
4. **Middleware Verification**: The authentication middleware was updated to check every incoming token's `jti` against the `revoked_tokens` table. If a match is found, the request is rejected with a `401 Unauthorized`.

## Consequences

**Positive:**
- Immediate and definitive session termination upon user logout.
- Prevents the reuse of stolen tokens after a user has actively logged out.
- Fits seamlessly into our existing JWT architecture without completely abandoning the stateless model (valid tokens still avoid complex session lookups, and the blocklist query is highly indexed).

**Negative:**
- Adds a database query to the critical path of every authenticated request, slightly increasing latency and database load.
- Requires maintenance logic (e.g., a background cron job or periodic cleanup) to delete rows from `revoked_tokens` where `expires_at < NOW()`, preventing the table from growing indefinitely.
