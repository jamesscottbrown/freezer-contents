# Default recipe: full build
default: build

# Full build: install deps, build UI, and compile Go binary
build: ui-install ui-build go-build

# Install UI dependencies
ui-install:
    cd ui && npm install

# Build UI only
ui-build:
    cd ui && npm run build

# Build Go binary (static, no CGO)
go-build:
    CGO_ENABLED=0 go build .

# Run Go linter
lint:
    golangci-lint run

# Preview/run the application
preview:
    go run .

# Clean build artifacts
clean:
    rm -f freezer-contents
    rm -rf ui/dist
