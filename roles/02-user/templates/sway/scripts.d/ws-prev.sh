current=$(swaymsg -t get_workspaces | jq '.[] | select(.focused).num')
if [ "$current" -gt 1 ]; then
    swaymsg workspace number $((current - 1))
fi