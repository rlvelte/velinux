import QtQuick
import qs.bar.widgets.hardware

HardwareMonitor {
    label: "RAM "
    command: ["sh", "-c", "free | awk '/Mem:/ {printf \"%.0f\", $3/$2 * 100}'"]
    pollInterval: 5000
    suffix: "%"
}
