#!/bin/zsh

# Alias
alias ga="git add ."
alias gs="git status"
alias gp="git push"
gco() { git commit "$1" "$2" }

# Hooks
autoload -U add-zsh-hook
typeset -gA _ONEFETCH_DONE

chpwd_onefetch() {
  local git_root=$(git rev-parse --show-toplevel 2>/dev/null)
  if [[ -n "$git_root" && -z "${_ONEFETCH_DONE[$git_root]}" ]]; then
    onefetch 2>/dev/null
    _ONEFETCH_DONE[$git_root]=1
  fi
}

add-zsh-hook chpwd chpwd_onefetch