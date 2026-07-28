import QtQuick
import Quickshell
import Quickshell.Hyprland
import Quickshell.I3
import qs.globals

Item {
    id: workspaces
    implicitWidth: wsRow.implicitWidth
    implicitHeight: Dimensions.barHeight

    property var wsModel: Quickshell.hyprland ? Hyprland.workspaces : I3.workspaces

    Row {
        id: wsRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 4

        Repeater {
            model: wsModel
            delegate: Rectangle {
                property var ws: modelData

                width: wsText.implicitWidth + 10
                height: 20
                color: "transparent"

                Text {
                    id: wsText
                    anchors.centerIn: parent
                    text: ws.name
                    color: ws.focused ? Theme.primary : Theme.subtext
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeHeading
                    font.weight: ws.focused ? Font.Bold : Font.Medium
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
