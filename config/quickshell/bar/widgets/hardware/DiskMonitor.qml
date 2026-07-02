import QtQuick
import Quickshell.Io
import qs.services

Item {
    id: diskMonitor
    implicitWidth: diskText.implicitWidth
    implicitHeight: 40
    property string usage: "--"

    Text {
        id: diskText
        anchors.verticalCenter: parent.verticalCenter
        text: "HOME " + usage
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
        id: diskProcess
        command: ["sh", "-c", "df -h /home --output=pcent | tail -1 | tr -d ' '"]
        stdout: StdioCollector {
            onStreamFinished: {
                diskMonitor.usage = text.trim()
            }
        }
    }

    Timer {
        interval: 30000
        running: true
        repeat: true
        triggeredOnStart: true
        onTriggered: diskProcess.running = true
    }
}
