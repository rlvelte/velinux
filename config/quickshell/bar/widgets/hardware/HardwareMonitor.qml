import QtQuick
import Quickshell.Io
import qs.globals

Item {
    id: root

    required property string label
    required property var command
    required property int pollInterval

    implicitWidth: monitorText.implicitWidth
    implicitHeight: Dimensions.barHeight

    property string usage: "--"

    Text {
        id: monitorText
        anchors.verticalCenter: parent.verticalCenter
        text: root.label + root.usage
        color: {
            var val = parseInt(root.usage)
            if (val >= 90) return Theme.error
            if (val >= 70) return Theme.warning
            return Theme.subtext
        }
        font.family: Theme.fontName
        font.pixelSize: Theme.fontSizeHeading
        font.weight: Font.Medium
        verticalAlignment: Text.AlignVCenter
    }

    Process {
        id: monitorProcess
        command: root.command
        stdout: StdioCollector {
            onStreamFinished: root.usage = text.trim()
        }
    }

    Timer {
        interval: root.pollInterval
        running: true
        repeat: true
        triggeredOnStart: true
        onTriggered: monitorProcess.running = true
    }
}
