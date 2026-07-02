import QtQuick
import Quickshell.Io
import qs.services

Item {
    id: memoryMonitor
    implicitWidth: memoryText.implicitWidth
    implicitHeight: 40
    property string usage: "--"

    Text {
        id: memoryText
        anchors.verticalCenter: parent.verticalCenter
        text: "RAM " + usage + "%"
        color: {
            var val = parseInt(usage)
            if (val >= 90) return Theme.error
            if (val >= 70) return Theme.warning
            return Theme.subtext
        }
        font.family: Theme.fontName
        font.pixelSize: 18
        font.weight: Font.Medium
        verticalAlignment: Text.AlignVCenter
    }

    Process {
        id: memoryProcess
        command: ["sh", "-c", "free | awk '/Mem:/ {printf \"%.0f\", $3/$2 * 100}'"]
        stdout: StdioCollector {
            onStreamFinished: {
                memoryMonitor.usage = text.trim()
            }
        }
    }

    Timer {
        interval: 5000
        running: true
        repeat: true
        triggeredOnStart: true
        onTriggered:                 memoryProcess.running = true
    }
}
