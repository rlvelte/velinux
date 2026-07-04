import QtQuick
import Quickshell
import Quickshell.Io
import qs.services

Row {
    id: root
    property var sysData: ({})
    width: parent.width
    height: 32
    spacing: 6
    visible: sysData.hostname !== undefined

    function refresh() {
        hwProcess.running = true
    }

    Text {
        id: hostTxt
        text: sysData.hostname || ""
        font.family: Theme.fontNameMono
        font.pixelSize: Theme.fontSizeHeading
        font.weight: Font.Bold
        color: Theme.text
        verticalAlignment: Text.AlignVCenter
        height: parent.height
    }

    Item { width: parent.width - hostTxt.implicitWidth - rightRow.implicitWidth - 12; height: 1 }

    Row {
        id: rightRow
        height: parent.height
        spacing: 6

        Text {
            text: sysData.uptime || ""
            font.family: Theme.fontNameMono
            font.pixelSize: Theme.fontSizeLarge
            color: Theme.subtext
            verticalAlignment: Text.AlignVCenter
            height: parent.height
        }
        Text {
            text: "-"
            font.family: Theme.fontName
            font.pixelSize: Theme.fontSizeLarge
            color: Theme.muted
            verticalAlignment: Text.AlignVCenter
            height: parent.height
        }
        Text {
            text: (sysData.os || "") + " " + (sysData.kernel || "")
            font.family: Theme.fontNameMono
            font.pixelSize: Theme.fontSizeLarge
            color: Theme.subtext
            verticalAlignment: Text.AlignVCenter
            height: parent.height
            elide: Text.ElideRight
        }
    }

    Process {
        id: hwProcess
        command: ["vlx", "cc", "hw", "system", "--json"]
        stdout: StdioCollector {
            onStreamFinished: {
                try { root.sysData = JSON.parse(text) } catch(e) { console.error("cc hw parse error:", e) }
            }
        }
    }

    Timer {
        interval: 5000
        running: true
        repeat: true
        triggeredOnStart: true
        onTriggered: root.refresh()
    }
}
