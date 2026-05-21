# SPEC-004: Generic Screen Endpoint for SDUI

**Status:** Implemented
**Date:** 2026-05-21
**Author:** Claude Code

---

## 1. Objective
Migrate from screen-specific endpoints (`/api/screen/home`) to a generic screen endpoint (`/api/screen/{screenId}`) that serves dynamic UI definitions based on a screen identifier. The server controls what each screen renders; the client only displays.

---

## 2. Motivation
The mobile app is migrating from per-screen endpoints to a single generic SDUI endpoint. This simplifies client code and allows the server to control screen composition without requiring app updates.

---

## 3. Technical Details

### 3.1 Endpoint
```
GET /api/screen/{screenId}
```

- `screenId` (path parameter): identifier of the screen to render (e.g., "home", "profile", "settings")
- No API versioning — the app has not launched yet, no real users to break

### 3.2 Response Format
Same `ScreenResponse` contract already defined in `internal/screen/models.go`:

```json
{
  "screen_name": "home",
  "pull_to_refresh": true,
  "components": [
    {
      "type": "hero_banner",
      "id": "banner_1",
      "data": {
        "title": "Descubra o SDUI no Android",
        "image_url": "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f",
        "action_id": "NAVIGATE_TO_DETAILS"
      }
    }
  ]
}
```

### 3.3 Routing
Use Go 1.22+ ServeMux pattern matching:
```go
mux.HandleFunc("GET /api/screen/{screenId}", screenHandler.GetScreen)
```

### 3.4 Handler Behavior
- Parse `screenId` from URL path using `r.PathValue("screenId")`
- Route to screen-specific builder based on `screenId`
- Return 404 for unknown `screenId` values
- Initial implementation: hardcode payloads per screenId (simple map or switch case)
- Compatibility: `screenId="home"` must return identical response to the current `/api/screen/home` endpoint

### 3.5 Screen Registry (Initial)
Hardcode known screens in the handler or a simple registry:

| screenId | screen_name | Description |
|----------|-------------|-------------|
| `home` | `home` | Home screen with hero banner and action button |

---

## 4. Implementation Checklist
- [x] `specs/004-api-generic-endpoint.md` - This specification
- [x] `internal/server/server.go` - Add route `GET /api/screen/{screenId}`
- [x] `internal/screen/handler.go` - Add `GetScreen` method for generic screen handling
- [x] Backward compatibility: `screenId="home"` returns identical response
- [x] `go build ./...` compiles without errors
- [x] `GET /api/screen/home` returns HTTP 200 with correct JSON
- [x] `GET /api/screen/invalid` returns HTTP 404
- [x] `GET /api/screen/{screenId}` responds correctly for valid screenIds

---

## 5. Success Criteria
1. `GET /api/screen/home` returns HTTP 200 with the same JSON as before
2. `GET /api/screen/{any_valid_id}` returns HTTP 200 with correct screen payload
3. `GET /api/screen/{unknown_id}` returns HTTP 404
4. No API versioning (no `/api/v2/...`)
5. No external dependencies beyond standard library
6. Code follows idiomatic Go patterns
