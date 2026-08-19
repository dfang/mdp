# List available recipes
default:
    @just --list

# Build mdp binary
build:
    go build -o bin/mdp .

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
    go install .

# Generate CHANGELOG.md using git-cliff
changelog:
    git-cliff -o CHANGELOG.md

# Clean build artifacts
clean:
    rm -rf bin/ mdp

