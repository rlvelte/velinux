swaymsg -t get_tree | \
  jq -r '.. | select(.type?) | select(.focused) | "\(.rect.x),\(.rect.y) \(.rect.width)x\(.rect.height)"' | \
  grim -g - && notify-send "Screenshot" "Window"
