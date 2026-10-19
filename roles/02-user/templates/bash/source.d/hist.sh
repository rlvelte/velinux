#!/bin/bash

# History
HISTFILE="$BASHDOTDIR/.hist_bash"
HISTSIZE=50000
HISTFILESIZE=50000
HISTTIMEFORMAT='%F %T '
HISTCONTROL=ignoreboth:erasedups

shopt -s histappend
shopt -s checkwinsize

__share_history() {
  history -a
  history -c
  history -r
}

PROMPT_COMMAND="__share_history${PROMPT_COMMAND:+;$PROMPT_COMMAND}"
