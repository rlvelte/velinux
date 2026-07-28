import QtQuick
import qs.bar.widgets.hardware

HardwareMonitor {
    label: "HOME "
    command: ["sh", "-c", "df -h /home --output=pcent | tail -1 | tr -d ' '"]
    pollInterval: 30000
}
