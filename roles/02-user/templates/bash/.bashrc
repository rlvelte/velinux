#            _ _                  _               _     
# __   _____| (_)_ __  _   ___  _| |__   __ _ ___| |__  
# \ \ / / _ \ | | '_ \| | | \ \/ / '_ \ / _` / __| '_ \ 
#  \ V /  __/ | | | | | |_| |>  <| |_) | (_| \__ \ | | |
#   \_/ \___|_|_|_| |_|\__,_/_/\_\_.__/ \__,_|___/_| |_|
#                                         bash by rvelte

# Interactivity
case $- in *i*) ;; *) return 0 ;; esac

# Distrobox compatibility
if [ -f "/run/.containerenv" ] || [ -n "$CONTAINER_ID_ANY" ]; then
    export IN_DISTROBOX=1

    PS1='\[\e[32m\]distrobox as \u:\[\e[0m\] \w \$ '
    alias ll='ls -lah'
    alias c='cat'

    return 0
fi

# Bash
export BASHDOTDIR="$HOME/.config/bash"
export OMB="$BASHDOTDIR/oh-my-bash"

# Configurations
for f in "$BASHDOTDIR"/source.d/*.sh; do
    [ -r "$f" ] && . "$f"
done

# Plugins
plugins=(
    sudo 
    fzf
)

# Source
[ -r "$OMB/oh-my-bash.sh" ] && . "$OMB/oh-my-bash.sh"
eval "$(starship init bash)"

HISTFILE="$BASHDOTDIR/.hist_bash"
