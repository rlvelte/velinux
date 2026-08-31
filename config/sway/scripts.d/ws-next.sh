current=$(swaymsg -t get_workspaces | jq '.[] | select(.focused).num')
if [ "$current" -lt 5 ]; then
    swaymsg workspace number $((current + 1))
fi