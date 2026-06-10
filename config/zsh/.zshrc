#            _ _                           _     
# __   _____| (_)_ __  _   ___  __ _______| |__  
# \ \ / / _ \ | | '_ \| | | \ \/ /|_  / __| '_ \ 
#  \ V /  __/ | | | | | |_| |>  <  / /\__ \ | | |
#   \_/ \___|_|_|_| |_|\__,_/_/\_\/___|___/_| |_|
#                                  Zsh by rlvelte

# Distrobox compatibility
if [ -f "/run/.containerenv" ] || [ -n "$CONTAINER_ID_ANY" ]; then
    export IN_DISTROBOX=1

    PS1='%F{green}distrobox as %n:%f %~ %# '
    alias ls='ls --color=auto'
    alias ll='ls -lah'
    alias c='cat'

    return
fi

# ZSH
export ZSH="$ZDOTDIR/oh-my-zsh"
export ZSH_CUSTOM="$ZDOTDIR/oh-my-zsh/custom"
export ZSH_THEME="powerlevel10k/powerlevel10k"

# Configuration 
for f in $ZDOTDIR/config.d/*.sh(N); do
    source "$f"
done

# Plugins
plugins=(
    sudo
    zsh-interactive-cd
    zsh-autosuggestions
    zsh-syntax-highlighting
    zsh-eza
)

# Source
source $ZSH/oh-my-zsh.sh
source $ZDOTDIR/.p10k.zsh
HISTFILE="$ZDOTDIR/.hist_zsh"