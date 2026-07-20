# Authentication Security Remediation Report

**Date**: 2026-07-20
**Scope**: Backend Authentication Service and Handlers

## Executive Summary

Following a comprehensive backend code review, 11 critical and high-priority security vulnerabilities were identified in the authentication system. The previous implementation suffered from several architectural flaws, including stateless JWTs that could not be revoked on logout, long-lived access tokens, vulnerability to enumeration attacks, and poor cookie security practices. 

This report outlines what was changed, why it was necessary, and how it was implemented.

---

## Detailed Changes

### 1. Stateful JWT Revocation (AUTH-01)
* **Why it was needed**: JWTs are stateless. Previously, when a user logged out, the token was deleted from the client-side cookie, but the token itself remained valid until it expired (7 days later). If an attacker intercepted the token, they could continue using it even after the user explicitly logged out.
* **How it was changed**: We implemented a server-side blocklist for revoked tokens.
  * Added a `jti` (JWT ID - a UUID) to every issued access token.
  * Created a `revoked_tokens` PostgreSQL table.
  * Updated the `/api/auth/logout` endpoint to parse the incoming token, extract its `jti`, and insert it into the `revoked_tokens` table.
  * Updated the Auth Middleware to check if the incoming `jti` exists in the blocklist, rejecting the request with a `401 Unauthorized` if found.

### 2. Refresh Token Rotation Architecture (AUTH-09)
* **Why it was needed**: Access tokens lived for 7 days, which is dangerously long if a token is compromised. A secure system should use short-lived access tokens paired with long-lived refresh tokens.
* **How it was changed**: 
  * Reduced the JWT access token lifespan from 7 days to **15 minutes**.
  * Created a `refresh_tokens` PostgreSQL table.
  * On login/register, the server now generates a cryptographically secure random 32-byte string (Refresh Token). This token is hashed using SHA-256 and stored in the database, while the raw token is returned to the client as an `HttpOnly` cookie.
  * Added a new `/api/auth/refresh` endpoint. When the access token expires, the client calls this endpoint. The server hashes the refresh token from the cookie, checks the database, deletes the old refresh token, and issues a brand new access and refresh token pair (Token Rotation).

### 3. Account Enumeration Prevention (AUTH-05, AUTH-06, AUTH-11)
* **Why it was needed**: The registration endpoint returned a `409 Conflict` if an email was already registered, allowing an attacker to quickly check if a specific user (e.g., a CEO or target) had an account. Furthermore, the response time was much faster for existing emails (since the server skipped the slow bcrypt hashing), exposing a timing vulnerability. Finally, emails were not normalized, leading to potential duplicate accounts or login failures depending on casing.
* **How it was changed**:
  * **Normalization (AUTH-11)**: All emails are now trimmed and lowercased during Login and Registration.
  * **Status Code (AUTH-05)**: The registration endpoint now unconditionally returns a `202 Accepted` with a generic message ("If this email is not already registered, your account has been created") when an email already exists.
  * **Timing Equalization (AUTH-06)**: The `AuthService` now computes a dummy bcrypt hash on startup. If an email already exists during registration, the server executes `bcrypt.CompareHashAndPassword` against the dummy hash to ensure the response time is indistinguishable from a successful registration.

### 4. Rate Limiting Proxy Support (AUTH-04)
* **Why it was needed**: The rate limiter was strictly reading `r.RemoteAddr`. If the backend was placed behind a reverse proxy (like Nginx or a cloud load balancer), all requests would appear to come from the proxy's IP, meaning all users would share a single rate limit bucket and quickly get locked out.
* **How it was changed**: Created a `realIP()` helper function that checks the `X-Real-IP` and `X-Forwarded-For` HTTP headers to extract the true client IP before falling back to `r.RemoteAddr`.

### 5. Cookie Security Hardening (AUTH-03, AUTH-08, AUTH-10)
* **Why it was needed**: Cookies lacked strict cross-site protections and the `Logout` handler accidentally cleared cookies without the `Secure` flag in HTTPS environments. Additionally, the JSON response body leaked the raw token, encouraging insecure client-side storage (like localStorage).
* **How it was changed**:
  * **SameSite Strict (AUTH-08)**: Upgraded all auth cookies (`token`, `refresh_token`) from `SameSite=Lax` to `SameSite=Strict` to eliminate cross-site request forgery (CSRF) vectors.
  * **Secure Flag Sync (AUTH-03)**: Ensured the `Logout` handler dynamically respects the `IS_HTTPS` environment variable, just like `Login`.
  * **Token Removal (AUTH-10)**: Removed the `token` field from the JSON `AuthResponse` struct. The frontend must now rely entirely on the secure `HttpOnly` cookies.

### 6. Cryptographic Hardening (AUTH-02, AUTH-07)
* **Why it was needed**: The JWT parsing logic did not explicitly enforce the expected signing algorithm (`HS256`), making it vulnerable to "none" algorithm attacks or public-key confusion attacks. Additionally, the default bcrypt cost was hardcoded and potentially too low for modern hardware.
* **How it was changed**:
  * **Algorithm Assertion (AUTH-02)**: The JWT `KeyFunc` now explicitly checks `if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok` and rejects the token if it doesn't match `HS256`.
  * **Configurable Bcrypt (AUTH-07)**: Introduced the `BCRYPT_COST` environment variable with bounds checking. It defaults to `12` (up from standard defaults) for production, but allows tests to run at cost `4` for speed.

---

## Architectural Decision Records (ADRs) Generated
1. **[ADR-006: Stateful JWT Revocation via Blocklist](../architecture/adrs/006-stateful-jwt-revocation.md)**
2. **[ADR-007: Refresh Token Rotation architecture](../architecture/adrs/007-refresh-token-rotation.md)**
