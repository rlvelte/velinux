#!/bin/zsh

# History
HISTSIZE=50000
SAVEHIST=50000

setopt EXTENDED_HISTORY          # record timestamps
setopt HIST_EXPIRE_DUPS_FIRST    # expire duplicates first when trimming
setopt HIST_IGNORE_DUPS          # don't record an entry that duplicates the previous
setopt HIST_IGNORE_SPACE         # don't record lines starting with a space
setopt HIST_FIND_NO_DUPS         # don't show duplicates in search results
setopt HIST_REDUCE_BLANKS        # trim leading/trailing whitespace
setopt INC_APPEND_HISTORY        # write to history file immediately, not on exit
setopt SHARE_HISTORY             # share history across sessions
