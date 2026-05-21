# SPEC-002: Screen Behavior Flags (Pull-to-Refresh & Future Extensions)

**Status:** Draft
**Date:** 2026-05-21
**Author:** Claude Code

---

## 1. Objective
Add screen-level behavior flags to the SDUI BFF API to allow clients to conditionally enable UI interactions (starting with pull-to-refresh). The design must be extensible to support future behaviors like infinite scroll, swipe actions, etc.

---

## 2. Background & Motivation
The current `ScreenResponse` only carries `screen_name` and `components`. Clients need to know if a screen supports interactions like pull-to-refresh without hardcoding screen names. This spec introduces a simple, boolean-based approach for individual behavior flags, with room for structured config later if needed.

---

## 3. Affected Files
- `internal/screen/models.go` — Add `PullToRefresh` field to `ScreenResponse`
- `internal/screen/handler.go` — Set `PullToRefresh: true` for home screen
- `internal/screen/handler.go` — Future screens will set their own flags

---

## 4. Model Changes

### 4.1 Updated `ScreenResponse` (`internal/screen/models.go`)
```go
type ScreenResponse struct {
    ScreenName    string      `json:"screen_name"`
    PullToRefresh bool        `json:"pull_to_refresh"`
    Components    []Component `json:"components"`
}
```

---

## 5. Updated JSON Payload

### 5.1 `GET /api/screen/home` (Updated)
```json
{
  "screen_name": "home_page",
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
    },
    {
      "type": "action_button",
      "id": "btn1",
      "data": {
        "text": "Clique para Recarregar",
        "action_id": "RELOAD_SCREEN"
      }
    }
  ]
}
```

---

## 6. Implementation Checklist
- [x] Update `ScreenResponse` struct with `PullToRefresh` field
- [x] Update `GetHomeScreen` handler to set `PullToRefresh: true`
- [ ] Add `infinite_scroll` flag when needed (future)
- [ ] Document client-side behavior for pull-to-refresh

---

## 7. Success Criteria
1. `go build ./...` compiles without errors
2. `GET /api/screen/home` returns `"pull_to_refresh": true`
3. Other screen endpoints can set `pull_to_refresh` independently
4. JSON field names use snake_case
5. No external dependencies added

---

## 8. Future Extensions
When adding more behaviors (e.g., infinite scroll), follow the same pattern:
```go
type ScreenResponse struct {
    ScreenName     string      `json:"screen_name"`
    PullToRefresh  bool        `json:"pull_to_refresh"`
    InfiniteScroll bool        `json:"infinite_scroll"`
    Components     []Component `json:"components"`
}
```

If the number of flags grows beyond 3-4, consider grouping into a `config` object (see SPEC-002-ADDENDUM if needed).
