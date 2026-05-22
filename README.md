# Go BFF for Server-Driven UI (SDUI)

> A lightweight, idiomatic Go backend-for-frontend that serves dynamic, polymorphic screen definitions for Server-Driven UI applications.

[![Go Version](https://img.shields.io/badge/Go-1.22+-00ADD8?style=flat&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)
[![GitHub](https://img.shields.io/badge/GitHub-soulnanx/api--context--sdui-black?style=flat&logo=github)](https://github.com/soulnanx/api-context-sdui)

---

## 📖 Overview

This project demonstrates a **Server-Driven UI (SDUI)** architecture where the backend defines the UI structure as data, and the client (Android, iOS, Web) renders it dynamically. The BFF (Backend-for-Frontend) serves screen configurations with polymorphic components using `json.RawMessage` for flexible, decoupled component schemas.

### Key Features
- ✅ **Zero external dependencies** - Pure Go standard library
- ✅ **Polymorphic components** - `json.RawMessage` for flexible schemas
- ✅ **Go 1.22+ routing** - Native `net/http.ServeMux` with pattern matching
- ✅ **SPEC Driven Development** - Specifications in `specs/`
- ✅ **Graceful shutdown** - Proper signal handling and timeouts
- ✅ **Standard Go layout** - Idiomatic project structure
- ✅ **Admin API** - Dynamic screen management via in-memory CRUD endpoints

---

## 🏗️ Architecture

### Polymorphic Component Design

The core idea: **Components are polymorphic**. The `Component` envelope carries metadata (`type`, `id`), while `data` contains type-specific JSON deferred via `json.RawMessage`.

```
ScreenResponse
├── screen_name: "home_page"
└── components: [
    {
      "type": "hero_banner",
      "id": "b1",
      "data": { ... }  ← HeroBanner schema
    },
    {
      "type": "action_button",
      "id": "btn1",
      "data": { ... }  ← ActionButton schema
    }
]
```

### Request Flow

```
Client (Android/iOS/Web)
    ↓ GET /api/screen/home
Server (Go BFF)
    ↓ Marshal component structs → json.RawMessage
ScreenResponse (polymorphic JSON)
    ↓
Client renders UI dynamically
```

---

## 📁 Project Structure

```
api-context-sdui/
├── cmd/
│   └── api/
│       └── main.go              # Entry point, server startup
├── internal/
│   ├── component/
│   │   └── models.go           # Component schemas (HeroBanner, ActionButton)
│   ├── screen/
│   │   ├── handler.go          # HTTP handlers (screen + admin)
│   │   ├── storage.go         # In-memory screen store with RWMutex
│   │   └── models.go          # Core SDUI envelopes (Component, ScreenResponse)
│   └── server/
│       └── server.go           # Routing (Go 1.22 mux) + middleware
├── specs/
│   └── 001-sdui-bff-api.md    # Technical specification (SPEC Driven)
├── go.mod
├── .gitignore
└── README.md
```

---

## 🚀 Quick Start

### Prerequisites
- Go 1.22+ ([install guide](https://golang.org/doc/install))

### Run Locally

```bash
# Clone the repository
git clone https://github.com/soulnanx/api-context-sdui.git
cd api-context-sdui

# Run the server
go run cmd/api/main.go
```

Server starts on `http://localhost:8080` with graceful shutdown on `SIGINT`/`SIGTERM`.

### Build Binary

```bash
go build -o bin/api ./cmd/api
./bin/api
```

---

## 🛠️ Makefile Commands

The project includes a Makefile to simplify common development tasks:

| Target | Description |
|--------|-------------|
| `make build` | Compile binary to `bin/api` |
| `make run` | Run server directly via `go run` |
| `make test` | Run all tests with verbose output |
| `make test-coverage` | Generate HTML coverage report (`coverage.html`) |
| `make clean` | Remove `bin/`, `coverage.out`, `coverage.html` |
| `make fmt` | Format code with `gofmt` |
| `make vet` | Run `go vet` on all packages |
| `make tidy` | Run `go mod tidy` |
| `make check` | Run `fmt`, `vet`, and `test` sequentially |
| `make dev` | Start server and verify with curl |
| `make spec` | List available specs in `specs/` |

Example:
```bash
# Build and run
make build
./bin/api

# Or use check before committing
make check
```

---

## 📡 API Reference

### GET /api/screen/{screenId}

Returns a screen configuration by its identifier. Reads from the in-memory screen store; screens can be created/updated/deleted via the Admin API.

**Request:**
```bash
curl http://localhost:8080/api/screen/home
```

**Parameters:**
- `screenId` (path): Screen identifier (e.g., `home`, `profile`, `settings`)

**Response (200 OK):**
```json
{
  "screen_name": "home",
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

**Response (404 Not Found):**
```json
{
  "error": "screen not found"
}
```

### Admin API (In-Memory Screen Management)

All admin endpoints are mounted under `/api/admin/screen/`. Screens are stored in memory; a server restart resets all data to the initial `home` screen.

| Method | Path | Description |
|--------|------|-------------|
| `POST` | `/api/admin/screen/{screenId}` | Create a new screen (409 if exists) |
| `GET` | `/api/admin/screen/{screenId}` | Retrieve screen config (404 if not found) |
| `PUT` | `/api/admin/screen/{screenId}` | Update existing screen (404 if not found) |
| `DELETE` | `/api/admin/screen/{screenId}` | Delete screen (404 if not found) |

#### POST /api/admin/screen/{screenId}

Request body (no `screen_name` — it is taken from the URL):
```json
{
  "pull_to_refresh": true,
  "components": [
    {
      "type": "hero_banner",
      "id": "b1",
      "data": {
        "title": "My Screen",
        "image_url": "https://example.com/img.jpg",
        "action_id": "NAVIGATE"
      }
    }
  ]
}
```

Response: HTTP 201 Created with the created screen (includes `screen_name`), 409 if screen already exists.

#### GET /api/admin/screen/{screenId}

Response: HTTP 200 with full `ScreenResponse` JSON, 404 if not found.

#### PUT /api/admin/screen/{screenId}

Request body: same as POST. Fully replaces the screen configuration.
Response: HTTP 200 with updated screen, 404 if not found.

#### DELETE /api/admin/screen/{screenId}

Response: HTTP 204 No Content on success, 404 if not found.

---

## ➕ Adding a New Component

To add a new UI component (e.g., `text_input`):

### 1. Define the Component Schema

Edit `internal/component/models.go`:

```go
// Add new struct
type TextInput struct {
    Label    string `json:"label"`
    Placeholder string `json:"placeholder"`
    ActionID string `json:"action_id"`
}
```

### 2. Register in Screen Handler

Edit `internal/screen/handler.go`:

```go
func (h *ScreenHandler) GetHomeScreen(w http.ResponseWriter, r *http.Request) {
    // ... existing components ...

    textInput := component.TextInput{
        Label:       "Email",
        Placeholder: "user@example.com",
        ActionID:    "SUBMIT_EMAIL",
    }
    textInputData, _ := json.Marshal(textInput)

    resp := ScreenResponse{
        ScreenName: "home_page",
        Components: []Component{
            // ... existing ...
            {Type: "text_input", ID: "input1", Data: textInputData},
        },
    }
    // ...
}
```

### 3. Update the SPEC

Document the new component in `specs/001-sdui-bff-api.md`.

---

## 📱 Android Client Integration

### Network Configuration

Add to `res/xml/network_security_config.xml`:

```xml
<?xml version="1.0" encoding="utf-8"?>
<network-security-config>
    <domain-config cleartextTrafficPermitted="true">
        <domain includeSubdomains="true">10.0.2.2</domain>
    </domain-config>
</network-security-config>
```

Reference in `AndroidManifest.xml`:

```xml
<application
    android:networkSecurityConfig="@xml/network_security_config">
    <uses-permission android:name="android.permission.INTERNET" />
</application>
```

### Access URLs

| Environment | URL |
|------------|-----|
| Emulator | `http://10.0.2.2:8080/api/screen/home` |
| Physical Device (same Wi-Fi) | `http://192.168.3.30:8080/api/screen/home` |

### Kotlin Example (Retrofit)

```kotlin
data class Component(
    val type: String,
    val id: String,
    val data: Map<String, Any>
)

data class ScreenResponse(
    val screen_name: String,
    val components: List<Component>
)

interface SduiApiService {
    @GET("api/screen/{screen}")
    suspend fun getScreen(@Path("screen") screen: String): ScreenResponse
}

// Create Retrofit instance
val retrofit = Retrofit.Builder()
    .baseUrl("http://10.0.2.2:8080/")  // Emulator
    .addConverterFactory(GsonConverterFactory.create())
    .build()
```

---

## 📐 SPEC Driven Development

This project follows **SPEC Driven Development**. All features start with a specification in the `specs/` directory:

```
specs/
├── 001-sdui-bff-api.md       # Initial BFF specification (Implemented)
├── 002-*.md                  # Additional specs (Implemented)
├── 004-*.md                  # Additional specs (Implemented)
├── 005-makefile.md           # Makefile workflow (Implemented)
└── 006-admin-api-inmemory.md # Admin API with in-memory storage (Implemented)
```

Each spec includes:
- Objective and constraints
- Technical stack
- Project structure
- Architectural decisions
- Target JSON payloads
- Success criteria

**Workflow:**
1. Write spec in `specs/`
2. Implement following spec exactly
3. Validate against success criteria
4. Mark spec as "Implemented"

---

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Test endpoint manually
curl http://localhost:8080/api/screen/home | jq .
```

---

## 🛠️ Development

### Project Layout

Following [Standard Go Project Layout](https://github.com/golang-standards/project-layout):

- `cmd/` - Main applications
- `internal/` - Private application code (component, screen, server)
- `specs/` - Technical specifications

### Code Style

- Idiomatic Go (effective go)
- Standard library only
- `json.RawMessage` for polymorphism
- snake_case JSON tags

---

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

---

## 🤝 Contributing

1. Check `specs/` for existing specifications
2. For new features, create a new spec in `specs/`
3. Implement following the spec
4. Update documentation
5. Submit a pull request

---

## 📚 References

- [Server-Driven UI Pattern](https://martinfowler.com/articles/responsive-rest-api.html)
- [Go Standard Library net/http](https://pkg.go.dev/net/http)
- [Effective Go](https://golang.org/doc/effective_go)
- [Standard Go Project Layout](https://github.com/golang-standards/project-layout)

---

<p align="center">
  Built with ❤️ using Go standard library only
</p>
