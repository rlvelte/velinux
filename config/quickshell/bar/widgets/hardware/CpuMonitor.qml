import QtQuick
import qs.bar.widgets.hardware

HardwareMonitor {
    label: "CPU "
    command: ["sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{printf \"%.0f\", $2}'"]
    pollInterval: 5000
}
