import QtQuick
import Quickshell.Io
import qs.services

Item {
    id: cpuMonitor
    implicitWidth: cpuText.implicitWidth
    implicitHeight: 40
    property string usage: "--"

    Text {
        id: cpuText
        anchors.verticalCenter: parent.verticalCenter
        text: "CPU " + usage + "%"
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
        id: cpuProcess
        command: ["sh", "-c", "top -bn1 | grep 'Cpu(s)' | awk '{printf \"%.0f\", $2}'"]
        stdout: StdioCollector {
            onStreamFinished: {
                cpuMonitor.usage = text.trim()
            }
        }
    }

    Timer {
        interval: 5000
        running: true
        repeat: true
        triggeredOnStart: true
        onTriggered: cpuProcess.running = true
    }
}
