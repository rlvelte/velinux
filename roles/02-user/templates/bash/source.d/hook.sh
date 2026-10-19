#!/bin/bash

# Git
declare -A _ONEFETCH_DONE
_ONEFETCH_LAST_PWD=""

__onefetch_hook() {
  [ "$PWD" = "$_ONEFETCH_LAST_PWD" ] && return 0
  _ONEFETCH_LAST_PWD="$PWD"

  local git_root
  git_root="$(git rev-parse --show-toplevel 2>/dev/null)" || return 0
  [ -n "${_ONEFETCH_DONE[$git_root]}" ] && return 0

  onefetch 2>/dev/null
  _ONEFETCH_DONE[$git_root]=1
}
PROMPT_COMMAND="__onefetch_hook${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
