import QtQuick
import Quickshell
import Quickshell.Hyprland
import Quickshell.I3
import qs.services

Item {
    id: workspaces
    implicitWidth: wsRow.implicitWidth
    implicitHeight: 40

    property var wsModel: Quickshell.hyprland ? Hyprland.workspaces : I3.workspaces
    property var wsFocused: Quickshell.hyprland ? (Hyprland.focusedWorkspace ? Hyprland.focusedWorkspace.id : -1) : (I3.focusedWorkspace ? I3.focusedWorkspace.id : -1)

    Row {
        id: wsRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 4

        Repeater {
            model: wsModel
            delegate: Rectangle {
                property var ws: modelData
                property bool isFocused: ws.id === workspaces.wsFocused

                width: wsText.implicitWidth + 10
                height: 20
                color: isFocused ? "transparent" : "transparent"
                radius: 4

                Text {
                    id: wsText
                    anchors.centerIn: parent
                    text: ws.name
                    color: isFocused ? Theme.primary : Theme.subtext
                    font.family: Theme.fontName
                    font.pixelSize: 18
                    font.weight: isFocused ? Font.Bold : Font.Medium
                }

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: ws.activate()
                }
            }
        }
    }
}
