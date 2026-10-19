#!/usr/bin/pwsh

# SSH
$env:SSH_AUTH_SOCK = "$(gpgconf --list-dirs agent-ssh-socket)"

# GnuPG
if (-not [Console]::IsInputRedirected) {
  $env:GPG_TTY = "$(tty)"
}

# Golang
$env:GOPATH = "$HOME/.go"
$env:GOMODCACHE = "$env:GOPATH/pkg/mod"
$env:GOCACHE = "$env:GOPATH/.cache/go-build"
$env:PATH += ":$env:GOPATH/bin"

# Npm
$env:NPM_CONFIG_CACHE = "$HOME/.cache/npm"

# Opencode
$env:PATH += ":$HOME/.opencode/bin"
