# ADR 001: Go and Chi for Backend Services

## Status
Accepted

## Context
We need a performant backend language to build RESTful API services capable of supporting concurrent traffic without memory bloat. Additionally, we need a router to handle parameterized URLs and middleware chains effectively.

## Decision
We chose **Go (Golang)** as the primary backend language.
For routing, we chose **go-chi/chi**.

## Rationale
- **Performance**: Go compiles to static binaries and provides exceptional concurrency via goroutines out of the box.
- **Simplicity**: Unlike full-blown frameworks (like NestJS or Spring), Go's standard library combined with `go-chi` is highly idiomatic. 
- **Chi vs Gin/Fiber**: `go-chi` leverages the standard `http.Handler` interfaces, meaning zero lock-in to proprietary context types. It remains extremely lightweight while supporting nested routing and custom middleware natively.

## Consequences
- Developers must manually manage dependency injection and write boiler-plate JSON marshal/unmarshal logic, rather than relying on framework magic.
- Integration tests are simple and don't require heavy framework servers.
