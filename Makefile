# Makefile - Main build file
# Includes all modular Makefile fragments

.PHONY: all shell local build test generate vmlinuxh help

# Include variable definitions
include make/Makefile.vars

# Include build tasks
include make/Makefile.build

# Include development environment tasks
include make/Makefile.dev

# Include test tasks
include make/Makefile.test

# Include documentation tasks
include make/Makefile.docs

# Default target
all: binary

# Help target - show available commands
help:
	@echo "Available targets:"
	@echo ""
	@echo "Build targets:"
	@echo "  make build          - Build binary files"
	@echo "  make binary         - Build binary using docker buildx"
	@echo "  make image          - Build development Docker image"
	@echo "  make install        - Install binary to /usr/local/bin"
	@echo ""
	@echo "Development targets:"
	@echo "  make shell          - Start interactive shell in development container"
	@echo ""
	@echo "Test targets:"
	@echo "  make unittest       - Run unit tests locally"
	@echo "  make test           - Run unit tests in Docker container"
	@echo "  make test.bin       - Test binary files (requires PKG variable)"
	@echo "  make e2e            - Run end-to-end tests (requires DIR variable)"
	@echo ""
	@echo "eBPF targets:"
	@echo "  make generate       - Generate eBPF .o files"
	@echo "  make vmlinuxh       - Generate vmlinux header file"
	@echo ""
	@echo "Documentation targets:"
	@echo "  make doc              - Update README.md and vulnerability table"
	@echo "  make update-vul-table - Update vulnerability table in README.md"
	@echo "  make update-readme    - Update README.md (command help and vulnerability table)"
	@echo ""
	@echo "Environment variables:"
	@echo "  SLIM_LDFLAGS=       - LDFLAGS for build (default: -s -w)"
	@echo "  CN=1                - Use Chinese mirrors for apt and Go proxy"
	@echo "  DEBUG=1             - Enable debug output (--progress=plain)"
	@echo "  TEST_ENV=           - Test environment variable; filters e2e test_envs by name"
	@echo "  PKG=                - Package name for test.bin"
	@echo "  DIR=                - Directory for e2e tests"
