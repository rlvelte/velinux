#!/bin/zsh

# Wrapper
docker() {
    case "$1" in
        nuke)
            echo "Tactical nuke ready to deploy."
            printf "You're sure Mr. President? [y/N]: "
            read -r confirm

            case "$confirm" in
                y|Y|yes|YES)
                    ;;
                *)
                    echo "Aborted."
                    return 1
                    ;;
            esac

            projects=$(command docker compose ls -q 2>/dev/null)
            for project in $projects; do
                command docker compose -p "$project" down -v --remove-orphans 2>/dev/null
            done

            containers=$(command docker ps -q)
            [ -n "$containers" ] && command docker stop $containers 2>/dev/null

            allContainers=$(command docker ps -aq)
            [ -n "$allContainers" ] && command docker rm -f $allContainers 2>/dev/null

            images=$(command docker images -aq)
            [ -n "$images" ] && command docker rmi -f $images 2>/dev/null

            volumes=$(command docker volume ls -q)
            [ -n "$volumes" ] && command docker volume rm -f $volumes 2>/dev/null

            networks=$(command docker network ls --filter type=custom -q)
            [ -n "$networks" ] && command docker network rm $networks 2>/dev/null

            command docker builder prune -af
            command docker system prune -af --volumes

            echo "What have you done?!"
            ;;
        *)
            command docker "$@"
            ;;
    esac
}
