BINARY_NAME=api
MODULE_NAME=api-context-sdui
CMD_PATH=./cmd/api
BUILD_DIR=./bin

.PHONY: build run test test-coverage clean fmt vet tidy check dev spec

build:
	go build -o $(BUILD_DIR)/$(BINARY_NAME) $(CMD_PATH)

run:
	go run $(CMD_PATH)

test:
	go test -v ./...

test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf $(BUILD_DIR) coverage.out coverage.html

fmt:
	gofmt -w .

vet:
	go vet ./...

tidy:
	go mod tidy

check: fmt vet test

dev:
	go run $(CMD_PATH) & sleep 2; curl -s http://localhost:8080/api/screen/home | head -20; kill %1 2>/dev/null || true

spec:
	@echo "Specs available:" && ls -1 specs/
