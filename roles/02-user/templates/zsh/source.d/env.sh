#!/bin/zsh

# Golang
export GOPATH="$HOME/.go"
export GOMODCACHE="$GOPATH/pkg/mod"
export GOCACHE="$GOPATH/.cache/go-build"
export PATH="$PATH:$GOPATH/bin"

# Npm
export NPM_CONFIG_CACHE="$HOME/.cache/npm"

# Opencode
export PATH="$PATH:$HOME/.opencode/bin"