# This Makefile is for local development purposes only
# Use at your own risk...

# Executables
GO := go

# Directories
INSTALL_DIR         := /usr/local/bin
BASH_COMPLETION_DIR := /usr/share/bash-completion/completions
ZSH_COMPLETION_DIR  := /usr/share/zsh/site-functions

.PHONY: build restore install

## install: Install binary and shell completions guard-wide
install: build
	sudo install -m 755 -t $(INSTALL_DIR) build/vlx
	sudo install -Dm644 build/completions/vlx.bash $(BASH_COMPLETION_DIR)/vlx
	sudo install -Dm644 build/completions/_vlx $(ZSH_COMPLETION_DIR)/_vlx

## build: Compile the vlx binary and generate shell completions
build: restore
	@mkdir -p build
	$(GO) build -o build/vlx .

	@mkdir -p build/completions
	./build/vlx completion bash > build/completions/vlx.bash
	./build/vlx completion zsh  > build/completions/_vlx

## restore: Download all Go module dependencies
restore:
	$(GO) mod download
	$(GO) mod tidy
	$(GO) vet ./...