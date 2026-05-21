# Project Guidelines

## Overview
Go BFF for Server-Driven UI (SDUI). Serves dynamic, polymorphic screen definitions via JSON. The server controls what each screen renders; the client (Android, iOS, Web) only displays.

## Stack
- **Language:** Go 1.22+ (native `net/http.ServeMux` with pattern matching)
- **Dependencies:** Standard library only — no external packages
- **Architecture:** Server-Driven UI with `json.RawMessage` for polymorphic components

## Project Structure
```
cmd/api/main.go          # Entry point, server startup
internal/
  component/models.go    # Component schemas (HeroBanner, ActionButton, etc.)
  screen/
    handler.go           # HTTP handlers (GetScreen, buildScreen, buildHomeScreen)
    models.go            # Core SDUI envelopes (Component, ScreenResponse)
  server/
    server.go            # Routing + middleware
specs/                   # SPEC Driven Development — all features start here
Makefile                  # Dev workflow (build, run, test, check, etc.)
```

## Key Patterns
- **SPEC Driven:** Every feature starts with a spec in `specs/` before implementation
- **Polymorphic components:** `Component` envelope (`type`, `id`, `data`) where `data` is `json.RawMessage`
- **Generic endpoint:** `GET /api/screen/{screenId}` — no per-screen endpoints, no API versioning
- **ScreenResponse:** `screen_name`, `pull_to_refresh`, `components[]`
- **JSON:** All field names use `snake_case`
- **Screen registry:** New screens added via `buildScreen()` switch case in `handler.go`

## Workflow
1. Write/update spec in `specs/`
2. Implement following the spec exactly
3. `make check` (fmt, vet, test) before committing
4. Update `README.md` if API changes
5. Commit and push to `origin/main` on GitHub

## Current State
- Specs 001, 002, 004: Implemented
- Spec 005 (Makefile): Implemented, needs commit
- No test files yet (`*_test.go` missing in all packages)
