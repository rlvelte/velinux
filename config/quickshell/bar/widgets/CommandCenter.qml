import QtQuick
import Quickshell
import Quickshell.Io
import qs.services

Item {
    id: commandCenter
    implicitWidth: commandCenterIcon.implicitWidth + 0
    implicitHeight: 40

    Image {
        id: commandCenterIcon
        anchors.centerIn: parent
        height: 28
        fillMode: Image.PreserveAspectFit
        source: Theme.logoPath

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
