# Incident Report: Backend Container Crash Loop

**Date**: 2026-07-20  
**Scope**: Docker Compose & Backend Environment Initialization

## What Was Fixed
The `gemma4-backend-1` Docker container was continuously failing to start (crashing immediately upon boot), preventing the backend API from serving requests. This issue has been fully resolved, and the backend is now starting and running correctly.

## Why It Needed to be Fixed
During a recent security remediation (specifically **SEC-01**), we hardened the backend's startup sequence. The application was previously falling back to an insecure, hardcoded JSON Web Token (JWT) secret if the `JWT_SECRET` environment variable was missing. We modified the backend to call a fatal panic (`log.Fatal`) on startup if this secret is missing, guaranteeing the application never runs in an insecure state.

However, the `docker-compose.yml` configuration was completely missing the `environment` injection for these newly required variables. Because Docker Compose did not pass the `JWT_SECRET` from the host's `.env` file into the container, the backend container booted with an empty secret, triggered the new security safeguard, and immediately crashed.

## How It Was Fixed
The `docker-compose.yml` file was updated to explicitly pass the necessary configuration variables into the `backend` service.

**Changes made to `docker-compose.yml`:**
```yaml
  backend:
    build:
      context: ./backend
      dockerfile: Dockerfile
    environment:
      - DATABASE_URL=postgres://${POSTGRES_USER:-todo_user}:${POSTGRES_PASSWORD:-todo_password}@db:5432/${POSTGRES_DB:-todo_db}?sslmode=disable
      - JWT_SECRET=${JWT_SECRET}                # <--- ADDED
      - BCRYPT_COST=${BCRYPT_COST:-12}          # <--- ADDED
      - IS_HTTPS=${IS_HTTPS:-false}             # <--- ADDED
```

By adding these lines, Docker Compose now reads the values from your `.env` file and safely injects them into the backend container during startup. Following a `docker compose up -d` to restart the stack, the backend successfully verified the presence of the `JWT_SECRET` and booted successfully.
