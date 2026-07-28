import QtQuick
import Quickshell
import Quickshell.Io
import qs.globals

Item {
    id: commandCenter
    implicitWidth: commandCenterIcon.implicitWidth
    implicitHeight: Dimensions.barHeight

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
}
