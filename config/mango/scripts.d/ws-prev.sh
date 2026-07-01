current=$(mmsg get all-monitors | jq -r '.monitors[] | select(.active == true) | .active_tags[0]')
if [ "$current" -gt 1 ]; then
    mmsg dispatch view,$((current - 1))
fi
