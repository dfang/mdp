version := `git describe --tags --always --dirty 2>/dev/null || echo "dev"`
commit := `git rev-parse --short HEAD 2>/dev/null || echo "none"`
date := `date -u +%Y-%m-%d`
ldflags := "-s -w -X main.version=" + version + " -X main.commit=" + commit + " -X main.date=" + date

# List available recipes
default:
    @just --list

# Build mdp binary
build:
    go build -ldflags '{{ldflags}}' -o bin/mdp .

# Run mdp with arguments (e.g. `just run README.md`)
run *args:
    go run . {{args}}

# Run tests
test:
    go test -v ./...

# Format Go source files
fmt:
    gofmt -s -w .

# Run go vet
vet:
    go vet ./...

# Install binary to GOPATH/bin
install:
    echo "install to $GOBIN"
    go install -ldflags '{{ldflags}}' .

# Generate CHANGELOG.md using git-cliff
changelog:
    git-cliff -o CHANGELOG.md

# Check goreleaser configuration
goreleaser-check:
    goreleaser check

# Build release snapshot locally with goreleaser
snapshot:
    goreleaser release --snapshot --clean

# Clean build artifacts
clean:
    rm -rf bin/ mdp dist/

