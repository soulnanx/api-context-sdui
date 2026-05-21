# SPEC-001: Go BFF for Server-Driven UI (SDUI)

**Status:** Implemented
**Date:** 2026-05-20
**Author:** Claude Code

---

## 1. Objective
Build a scalable Backend-for-Frontend (BFF) API in Go to serve dynamic, polymorphic screen definitions for a Server-Driven UI application. The project must be highly maintainable, ready to support new UI components, and follow idiomatic Go design patterns and community standards.

---

## 2. Technical Stack & Constraints
- **Language Version:** Go 1.22+ (utilizing the enhanced `net/http` ServeMux pattern matching).
- **Dependencies:** Standard library only. No external routers or frameworks (e.g., Gin, Chi) are allowed.
- **Polymorphism Strategy:** `json.RawMessage` for the component payload to defer serialization and keep component definitions decoupled.

---

## 3. Project Structure (Standard Go Layout)
The project must be structured according to the community standard project layout:

```
├── cmd/
│   └── api/
│       └── main.go         # Application entry point, initialization, and server startup
├── internal/
│   ├── component/
│   │   └── models.go       # Structs for specific UI components (Banner, Button, etc.)
│   ├── screen/
│   │   ├── handler.go      # HTTP handlers for serving screen configurations
│   │   └── models.go       # Core SDUI envelopes (ScreenResponse, Base Component)
│   └── server/
│       └── server.go       # Server routing configuration and middleware setup
└── specs/
    └── 001-sdui-bff-api.md # This specification
```

---

## 4. Architectural Design & Polymorphism

### 4.1 Core Contract (`internal/screen/models.go`)
The response represents a clean UI tree wrapper. The `Component` struct separates global metadata from the specific content using `json.RawMessage`:

```go
type Component struct {
    Type string          `json:"type"`
    ID   string          `json:"id"`
    Data json.RawMessage `json:"data"` // Deferred parsing for specific schemas
}

type ScreenResponse struct {
    ScreenName string      `json:"screen_name"`
    Components []Component `json:"components"`
}
```

### 4.2 Component Definition Strategy (`internal/component/models.go`)
Every new UI component type must have its own standalone struct representation. For this PoC, implement:

**hero_banner Data Schema:**
- Title (string)
- ImageURL (string)
- ActionID (string)

**action_button Data Schema:**
- Text (string)
- ActionID (string)

---

## 5. Target JSON Reference Payload
The endpoint `GET /api/screen/home` must accurately output the following contract:

```json
{
  "screen_name": "home_page",
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
- [x] `go.mod` - Module initialization with Go 1.22+
- [x] `internal/component/models.go` - HeroBanner and ActionButton structs
- [x] `internal/screen/models.go` - Component and ScreenResponse containers
- [x] `internal/screen/handler.go` - ScreenHandler with GetHomeScreen method
- [x] `internal/server/server.go` - Server routing with Go 1.22 mux and middleware
- [x] `cmd/api/main.go` - Entry point with graceful shutdown and timeouts
- [ ] Compile and test endpoint
- [ ] Create GitHub repository and push

---

## 7. Success Criteria
1. `go build ./...` compiles without errors
2. `GET /api/screen/home` returns HTTP 200 with exact JSON payload above
3. All JSON field names use snake_case
4. No external dependencies beyond standard library
5. Code follows idiomatic Go patterns
