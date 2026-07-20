# Golang Style Guide and Best Practices

This document outlines the standard best practices, formatting rules, and stylistic guidelines for writing idiomatic Go (Golang) code in this project. Following these guidelines ensures that our codebase remains clean, maintainable, and aligned with the broader Go community standards.

## 1. Core Principles

* **Simplicity and Readability:** Go is designed to be easy to read and understand. Favor clear, straightforward code over "clever" one-liners.
* **Idiomatic Go:** Embrace the "Go way" of doing things. Refer to [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments) as primary references.
* **"Clear is better than clever."** (from Go Proverbs)

## 2. Formatting

* **Automated Formatting:** Never format code manually. Always use `gofmt` or `goimports` on your code before committing. This enforces a universal formatting standard across the project.
* **Line Length:** Go has no strict line length limits, but aim to keep lines reasonably short for readability. Avoid unnecessarily long lines when a simple refactor can break them up.

## 3. Naming Conventions

* **Variables and Constants:** Use `camelCase` or `mixedCaps` (e.g., `userName`, `MAX_RETRIES` should be `maxRetries`). Avoid underscores (`_`).
* **Exported Identifiers:** Use `PascalCase` for variables, constants, functions, and types that need to be exported (public) outside the package (e.g., `func ParseData()`).
* **Unexported Identifiers:** Use `camelCase` for variables, constants, functions, and types internal to a package (e.g., `func parseData()`).
* **Receiver Names:** Keep receiver names short and descriptive, typically 1 to 3 letters (e.g., `func (c *Client) Fetch()`). Do not use `me`, `this`, or `self`.
* **Descriptiveness vs. Scope:** The longer the scope or lifetime of a variable, the more descriptive its name should be. Short names (like `i`, `v`, `k`) are fine for short scopes (e.g., loop indices).

## 4. Packages

* **Naming:** Package names should be short, concise, and ideally a single lowercase word (e.g., `http`, `bytes`, `time`). Avoid underscores or mixed caps in package names.
* **Meaningful Names:** Avoid utility "catch-all" package names like `util`, `common`, or `helper`. Group code by its functional domain instead (e.g., `stringset`, `auth`, `metrics`).
* **Avoid Stutter:** Do not repeat the package name in its exported identifiers. (e.g., in package `chubby`, use `chubby.File`, not `chubby.ChubbyFile`).

## 5. Control Structures & Idioms

* **Early Returns (Guard Clauses):** Prefer early returns to deep nesting. If an error occurs, handle it and return immediately, keeping the "happy path" flat and unindented.
  ```go
  // Bad
  if err == nil {
      // Do something
  } else {
      return err
  }

  // Good
  if err != nil {
      return err
  }
  // Do something
  ```
* **Zero Values:** Take advantage of Go's zero values. You don't need to explicitly initialize variables to `0`, `false`, `""`, or `nil` if that is their starting state.
  ```go
  // Prefer
  var count int
  var user User
  
  // Over
  count := 0
  user := User{}
  ```
* **Variable Declarations:** Use `:=` for initializing new variables inside functions when the type is apparent or not strictly necessary to define. Use `var` for declaring zero-value variables.

## 6. Error Handling

* **Explicit Checks:** Always check errors explicitly. Never ignore an error unless you have a documented reason (e.g., using `_ = r.Close()`).
* **Return Errors:** Functions that can fail should return an `error` as their last return value.
* **Avoid Panics:** Do not use `panic` for normal error handling. `panic` should be reserved strictly for truly unrecoverable programming errors or initialization failures.
* **Error Messages:** Error strings should not be capitalized (unless beginning with proper nouns or acronyms) and should not end with punctuation.
  ```go
  // Good: fmt.Errorf("failed to open file: %w", err)
  // Bad:  fmt.Errorf("Failed to open file.")
  ```

## 7. Interfaces

* **Keep Them Small:** Go interfaces should be as small as possible—often just one or two methods. This makes them highly composable. Examples include `io.Reader` and `io.Writer`.
* **"Accept interfaces, return concrete types."** (Generally preferred) Functions should accept the smallest interface necessary to do their job, but return concrete implementations.

## 8. Concurrency

* **Channels over Shared Memory:** Follow the Go proverb: *"Do not communicate by sharing memory; instead, share memory by communicating."* Prefer channels for orchestrating concurrent tasks and passing data between goroutines.
* **Mutexes:** Use `sync.Mutex` or `sync.RWMutex` when you truly need to protect a shared state or struct field from concurrent access.
* **Goroutine Lifecycles:** When starting a goroutine, always know how, when, and why it will stop. Avoid leaking goroutines.

## 9. Documentation

* **Comments:** Every exported (capitalized) identifier in a package should have a doc comment.
* **Format:** The comment should start with the name of the item it describes and be a complete sentence.
  ```go
  // Request represents an HTTP request received by a server.
  type Request struct { ... }
  ```
* **Package Comments:** Each package should have a package comment (usually in `doc.go` or above the `package` declaration in the main file of the package) that introduces the package and provides information relevant to the package as a whole.

## 10. Tools & Static Analysis

* We highly encourage running static analysis tools as part of your development workflow.
* Always run `go vet` to catch potential suspicious constructs.
* Consider integrating `golangci-lint` to enforce these guidelines and catch common errors automatically.
