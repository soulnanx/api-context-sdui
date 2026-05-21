# SPEC-005: Makefile for Development Workflow

**Status:** Draft
**Date:** 2026-05-21
**Author:** Claude Code

---

## 1. Objective
Create a Makefile to simplify common development tasks (build, run, test, lint) and standardize the development workflow.

---

## 2. Motivation
The project lacks automation for repetitive tasks. A Makefile reduces friction when building, testing, and running the server — especially for newcomers who shouldn't need to remember Go CLI flags or directory paths.

---

## 3. Technical Details

### 3.1 Makefile Targets

| Target | Description |
|--------|-------------|
| `build` | Compile binary to `bin/api` |
| `run` | Run server directly via `go run` |
| `test` | Run all tests with verbose output |
| `test-coverage` | Generate HTML coverage report (`coverage.html`) |
| `clean` | Remove `bin/`, `coverage.out`, `coverage.html` |
| `fmt` | Format code with `gofmt` |
| `vet` | Run `go vet` on all packages |
| `tidy` | Run `go mod tidy` |
| `check` | Run `fmt`, `vet`, and `test` sequentially |
| `dev` | Start server and curl the home endpoint to verify it's up |
| `spec` | List available specs in `specs/` |

### 3.2 Configuration

```makefile
BINARY_NAME=api
CMD_PATH=./cmd/api
BUILD_DIR=./bin
```

### 3.3 Post-Implementation: Update README

After the Makefile is implemented and verified, update `README.md`:
- Add a new **"Makefile Commands"** section after the "Quick Start" section
- List all available `make` targets with a one-line description each
- Keep the existing `go run` and `go build` examples (don't remove them)

---

## 4. Implementation Checklist
- [ ] `Makefile` - Create file in project root with all targets above
- [ ] `make build` - compiles to `bin/api`
- [ ] `make run` - starts server on `:8080`
- [ ] `make test` - runs `go test -v ./...`
- [ ] `make test-coverage` - generates `coverage.html`
- [ ] `make clean` - removes build artifacts
- [ ] `make fmt` - formats code
- [ ] `make vet` - runs static analysis
- [ ] `make tidy` - tidies go.mod
- [ ] `make check` - runs fmt, vet, test
- [ ] `make dev` - starts server and hits `/api/screen/home`
- [ ] `make spec` - lists files in `specs/`
- [ ] Update `README.md` with Makefile commands section
- [ ] `make check` passes with no errors

---

## 5. Success Criteria
1. All Makefile targets execute without errors
2. `make build` produces `bin/api` that runs correctly
3. `make test` passes all existing tests
4. `make check` is a one-stop command for CI-quality checks
5. README is updated with a Makefile commands reference
6. No external dependencies added (Makefile uses only Go CLI and standard Unix tools)
