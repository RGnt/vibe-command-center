# ADR 004: Playwright for End-to-End Testing

## Status
Accepted

## Context
Testing UI-heavy applications (drag and drop, complex state transitions) is brittle when using solely DOM-based component testing (like Jest/React Testing Library). We need a tool to run the application exactly as a user sees it in a real browser.

## Decision
We chose **Playwright** as our End-to-End (E2E) testing framework.

## Rationale
- **Isolation**: Playwright supports browser contexts, ensuring tests don't pollute each other's cookies or localStorage.
- **Speed**: Capable of parallel test execution out-of-the-box.
- **Resilience**: Features automatic waiting (e.g. `toBeVisible()`), heavily reducing the need for arbitrary `sleep()` calls which cause flaky tests.
- **Docker Integration**: Our suite boots up a completely fresh, isolated set of backend and database containers before running tests against them, guaranteeing that the test data is clean and isolated from development data.

## Consequences
- Requires a Node environment to run tests, even if the backend is Go.
- E2E tests are slower than unit tests, requiring them to be run separately via `docker-compose.test.yml`.
