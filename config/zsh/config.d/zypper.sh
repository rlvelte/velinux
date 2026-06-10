#!/bin/zsh

# Wrapper
zypper() {
    case "$1" in
        ins)
            shift

            if [[ $# -eq 0 ]]; then
                echo "Usage: zypper ins <pattern>"
                return 1
            fi

            command zypper packages --installed-only \
                | grep -i --color=auto "$@"
            ;;
        *)
            command zypper "$@"
            ;;
    esac
}
