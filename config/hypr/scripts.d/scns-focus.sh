hyprctl activewindow -j | \
  jq -r '"\(.at[0]),\(.at[1]) \(.size[0])x\(.size[1])"' | \
  grim -g - && notify-send "Screenshot" "Window"
