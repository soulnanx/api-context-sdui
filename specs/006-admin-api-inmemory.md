# SPEC-006: Admin API with In-Memory Screen Storage

**Status:** Implemented
**Date:** 2026-05-21
**Author:** Claude Code

---

## 1. Objective
Replace hardcoded screen definitions with an in-memory screen registry and expose admin endpoints to create, read, update, and delete screen configurations at runtime. This allows the API return to be configured without code changes or restarts (data survives only in memory; restarting the server resets to initial state).

---

## 2. Motivation
Currently, all screen definitions are hardcoded in `internal/screen/handler.go` via `buildScreen()` and `buildHomeScreen()`. Adding or changing a screen requires a code change and redeploy. An admin API backed by in-memory storage lets us configure screen responses dynamically through HTTP requests.

---

## 3. Technical Details

### 3.1 In-Memory Screen Storage

A package-level `sync.RWMutex`-protected map that stores screen configurations by screen ID:

```go
var screenStore = map[string]ScreenResponse{
    "home": { /* initial home screen, same as current hardcoded */ },
}
```

- No external dependencies (stdlib only, per project constraints)
- `sync.RWMutex` for concurrent read/write safety
- Server restart clears all data; initial state is the current "home" screen
- No persistence to disk

### 3.2 Admin Endpoints

All admin endpoints are mounted under `/api/admin/screen/`.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/screen/{screenId}` | Create a new screen with the given ID |
| `GET` | `/api/admin/screen/{screenId}` | Retrieve a screen configuration |
| `PUT` | `/api/admin/screen/{screenId}` | Update an existing screen (full replace) |
| `DELETE` | `/api/admin/screen/{screenId}` | Delete a screen configuration |

#### POST /api/admin/screen/{screenId}
Request body (same shape as `ScreenResponse` but without `screen_name` — it comes from the URL):

```json
{
  "pull_to_refresh": true,
  "components": [
    {
      "type": "hero_banner",
      "id": "b1",
      "data": {
        "title": "Descubra o SDUI no Android",
        "image_url": "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f",
        "action_id": "NAVIGATE_TO_DETAILS"
      }
    }
  ]
}
```

Response: HTTP 201 Created on success, 409 if screen already exists.

#### GET /api/admin/screen/{screenId}
Response: HTTP 200 with full `ScreenResponse` JSON (including `screen_name`). HTTP 404 if not found.

#### PUT /api/admin/screen/{screenId}
Request body: same as POST. Fully replaces the screen configuration.
Response: HTTP 200 on success, 404 if screen does not exist.

#### DELETE /api/admin/screen/{screenId}
Response: HTTP 204 on success, 404 if screen does not exist.

### 3.3 Changes to Existing GET /api/screen/{screenId}

The existing `GetScreen` handler must read from the in-memory store instead of the hardcoded `buildScreen()` switch case:

- Look up `screenId` in `screenStore`
- Return 404 if not found
- Return the stored `ScreenResponse` as JSON

The `buildScreen()` and `buildHomeScreen()` functions are removed.

### 3.4 Routing (server.go)

Add admin routes to the existing mux:

```go
mux.HandleFunc("POST /api/admin/screen/{screenId}", screenHandler.CreateScreen)
mux.HandleFunc("GET /api/admin/screen/{screenId}", screenHandler.GetScreenConfig)
mux.HandleFunc("PUT /api/admin/screen/{screenId}", screenHandler.UpdateScreen)
mux.HandleFunc("DELETE /api/admin/screen/{screenId}", screenHandler.DeleteScreen)
```

---

## 4. Implementation Checklist
- [x] `specs/006-admin-api-inmemory.md` — This specification
- [x] `internal/screen/storage.go` — Add in-memory screen store with RWMutex
- [x] `internal/screen/handler.go` — Remove `buildScreen()` and `buildHomeScreen()`; rewrite `GetScreen` to read from store; add `CreateScreen`, `GetScreenConfig`, `UpdateScreen`, `DeleteScreen` methods
- [x] `internal/server/server.go` — Register admin routes
- [x] Pre-populate store with the current "home" screen as initial state
- [x] `go build ./...` compiles without errors
- [x] `make check` passes (fmt, vet)
- [x] `GET /api/screen/home` returns HTTP 200 with correct JSON (backward compatibility)
- [x] `POST /api/admin/screen/test` + `GET /api/screen/test` returns the created screen
- [x] `PUT /api/admin/screen/home` updates the home screen
- [x] `DELETE /api/admin/screen/test` deletes the screen; subsequent GET returns 404

---

## 5. Success Criteria
1. In-memory store replaces all hardcoded screen definitions
2. `GET /api/screen/{screenId}` reads from the store (HTTP 200 / 404)
3. Admin API can create, read, update, and delete screens via HTTP
4. Existing `GET /api/screen/home` continues to work (backward compatible)
5. Concurrent reads/writes are safe (`sync.RWMutex`)
6. No external dependencies beyond standard library
7. All JSON field names use `snake_case`
8. `make check` passes
