selection=$(printf ' Lock\n󰢩 Logout\n Shutdown\n Reboot' | rofi \
  -config $XDG_CONFIG_HOME/rofi/config.d/power.rasi \
  -dmenu \
  -p Power)

case "$selection" in
  *Lock) hyprlock ;;
  *Logout) loginctl terminate-session self ;;
  *Shutdown) systemctl poweroff ;;
  *Reboot) systemctl reboot ;;
esac
exit 0