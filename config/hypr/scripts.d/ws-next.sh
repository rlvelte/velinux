current=$(hyprctl activeworkspace -j | jq '.id')
if [ "$current" -lt 4 ]; then
    hyprctl dispatch workspace $((current + 1))
fi