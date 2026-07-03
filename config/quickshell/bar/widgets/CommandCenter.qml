import QtQuick
import Quickshell
import Quickshell.Io
import qs.services

Item {
    id: commandCenter
    implicitWidth: commandCenterIcon.implicitWidth + 8
    implicitHeight: 40

    Image {
        id: commandCenterIcon
        anchors.centerIn: parent
        width: 22
        height: 22
        fillMode: Image.PreserveAspectFit
        source: {
            for (let i = 0; i < Vlx.themes.length; i++) {
                if (Vlx.themes[i].active) return Vlx.themes[i].logo
            }
            return ""
        }
    }

    MouseArea {
        anchors.fill: parent
        cursorShape: Qt.PointingHandCursor
        onClicked: toggleProcess.running = true
    }

    Process {
        id: toggleProcess
        command: ["quickshell", "ipc", "call", "commandCenter", "toggle"]
    }
}
