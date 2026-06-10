current=$(hyprctl activeworkspace -j | jq '.id')
if [ "$current" -gt 1 ]; then
    hyprctl dispatch workspace $((current - 1))
fi