#!/bin/bash

# SSH
export SSH_AUTH_SOCK="$(gpgconf --list-dirs agent-ssh-socket)"

# GnuPG
if tty -s 2>/dev/null; then
  export GPG_TTY="$(tty)"
fi

# Golang
export GOPATH="$HOME/.go"
export GOMODCACHE="$GOPATH/pkg/mod"
export GOCACHE="$GOPATH/.cache/go-build"
export PATH="$PATH:$GOPATH/bin"

# Npm
export NPM_CONFIG_CACHE="$HOME/.cache/npm"

# Opencode
export PATH="$PATH:$HOME/.opencode/bin"
