import QtQuick
import Quickshell
import Quickshell.Io
import qs.globals

Item {
    id: scratchpadWidget
    implicitWidth: scratchpadRow.implicitWidth
    implicitHeight: Dimensions.barHeight

    property int count: 0
    readonly property bool active: count > 0

    Row {
        id: scratchpadRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 6

        Text {
            text: scratchpadWidget.count
            color: scratchpadWidget.active ? Theme.text : Qt.darker(Theme.subtext, 2)
            font.family: Theme.fontName
            font.pixelSize: Theme.fontSizeHeading
            font.weight: Font.Medium
            verticalAlignment: Text.AlignVCenter
        }
    }

    Process {
        id: scratchpadProcess
        command: ["sh", "-c", "swaymsg -t get_tree | jq '[.. | objects | select(.scratchpad_state? != null and .scratchpad_state != \"none\")] | length'"]
        stdout: StdioCollector {
            onStreamFinished: {
                var n = parseInt(text.trim())
                if (!isNaN(n))
                    scratchpadWidget.count = n
            }
        }
    }

    Timer {
        interval: 1000
        running: true
        repeat: true
        triggeredOnStart: true
        onTriggered: scratchpadProcess.running = true
    }
}
