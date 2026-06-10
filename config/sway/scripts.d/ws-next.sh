current=$(swaymsg -t get_workspaces | jq '.[] | select(.focused).num')
if [ "$current" -lt 4 ]; then
    swaymsg workspace number $((current + 1))
fi