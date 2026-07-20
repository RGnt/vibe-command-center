# 007. Refresh Token Rotation Architecture

**Date:** 2026-07-20
**Status:** Accepted

## Context
Our application previously utilized a single JWT access token with a 7-day lifespan. This presented a significant security risk: if an attacker compromised an access token (e.g., via XSS, network interception, or local access), they would have unfettered access to the user's account for an entire week. 

While we implemented a JWT blocklist (see ADR-006) to allow explicit revocation on logout, this does not protect users who simply close their browser without logging out, or users who are unaware their token has been stolen.

## Decision
We decided to implement a **Refresh Token Rotation Architecture**.

1. **Short-Lived Access Tokens**: JWT access tokens have been reduced from a 7-day lifespan to **15 minutes**.
2. **Long-Lived Refresh Tokens**: Upon login/registration, the server generates a cryptographically secure 32-byte string. This is the Refresh Token, which lives for 30 days.
3. **Secure Storage**: The Refresh Token is hashed using SHA-256 before being stored in the `refresh_tokens` PostgreSQL table. The raw token is returned to the client as an `HttpOnly` cookie.
4. **Token Rotation Endpoint**: A new `/api/auth/refresh` endpoint was added. When the access token expires, the client calls this endpoint. The server:
   - Reads the raw refresh token from the cookie.
   - Hashes it and looks it up in the database.
   - If valid, the server deletes the *old* refresh token from the database.
   - Generates a brand new access token (15 minutes) and a brand new refresh token (30 days).
   - Returns the new tokens as cookies.
5. **Logout Invalidation**: The `/api/auth/logout` endpoint now deletes the refresh token from the database and clears both cookies.

## Consequences

**Positive:**
- **Drastically Reduced Attack Surface**: A stolen access token is now only useful for a maximum of 15 minutes.
- **Defense in Depth**: Even if an attacker steals the refresh token, they must use it before the legitimate user does. Because we implemented token rotation (deleting the old token on refresh), if the attacker uses it, the legitimate user will be unexpectedly logged out on their next refresh, serving as an indicator of compromise.
- **Database Security**: Because refresh tokens are hashed via SHA-256 in the database, a database dump does not expose usable refresh tokens to an attacker.

**Negative:**
- Increased architectural complexity. The client must now handle `401 Unauthorized` responses gracefully by calling the `/api/auth/refresh` endpoint and retrying the failed request.
- The `/api/auth/refresh` endpoint introduces additional database writes (deleting the old token, inserting the new one) every 15 minutes per active user.
