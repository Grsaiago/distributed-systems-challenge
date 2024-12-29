OUT_DIR = ./solutions


.PHONY: help
help:
	@echo "============================================="
	@echo "🚀  Welcome to the Distributed Systems Repo 🚀"
	@echo "============================================="
	@echo ""
	@echo "Available Commands:"
	@echo ""
	@echo "  make build       - Build all binaries into the ./solutions folder"
	@echo "  make run         - Run the distributed system locally"
	@echo "  make test        - Run all unit tests"
	@echo "  make clean       - Clean up build artifacts"
	@echo "  make fmt         - Format the code"
	@echo "  make lint        - Run the linter"
	@echo ""
	@echo "🔥  Pro Tip: Tweak the system config in 'config.yml'."
	@echo ""

	@echo kj
	@echo "build: 
	@echo "build/<echo | unique-id>: Build the specified binary to ./solutions folder"
	@echo "<echo | unique-id>: Run the specified binary"


# Run the echo binary
.PHONY: echo
echo:
	@go run ./cmd/echo/main.go

# Run the unique-id binary
.PHONY: unique-id
unique-id:
	@go run ./cmd/unique-id/unique-id.go

# Build all binaries into the ./solutions folder
build: build/echo build/unique-id build/broadcast-a

# Build the echo binary into the ./solutions folder
.PHONY: build/echo
build/echo:
	@go build -o $(OUT_DIR) ./cmd/echo/echo.go

# Build the unique-id binary into the ./solutions folder
.PHONY: build/unique-id
build/unique-id:
	@go build -o $(OUT_DIR) ./cmd/unique-id/unique-id.go

# Build the broadcast-a binary into the ./solutions folder
.PHONY: build/broadcast-a
build/broadcast-a:
	@go build -o $(OUT_DIR) ./cmd/broadcast-a/broadcast-a.go
