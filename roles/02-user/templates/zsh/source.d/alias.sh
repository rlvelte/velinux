#!/bin/zsh

# Cat/Bat
alias c='cat'
alias b='bat'

# Fastfetch
alias ff='fastfetch'
alias of='onefetch'

# Opencode
alias oc='opencode'

# Git
alias ga="git add ."
alias gs="git status"
alias gp="git push"
gco() { git commit "$@"; }

# Systemctl
alias shutdown='systemctl poweroff'
alias reboot='systemctl reboot'
alias bios='sudo systemctl reboot --firmware-setup'