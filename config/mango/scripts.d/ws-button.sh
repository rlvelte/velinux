#!/bin/bash
# Outputs workspace number for waybar.
# Usage: ws-button.sh <workspace_id>

ws_id="$1"
active=$(mmsg get all-monitors | jq -r '.monitors[] | select(.active == true) | .active_tags[0]')

if [ "$active" = "$ws_id" ]; then
    echo "[$ws_id]"
else
    echo " $ws_id "
fi