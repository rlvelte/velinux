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

# GnuPG
alias gpg-reset='gpgconf --kill gpg-agent && gpg-connect-agent /bye && ssh-add -l'

# Systemctl
alias shutdown='systemctl poweroff'
alias reboot='systemctl reboot'
alias bios='sudo systemctl reboot --firmware-setup'
