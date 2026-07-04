import QtQuick
import Quickshell
import Quickshell.Services.SystemTray
import qs.services

Item {
    id: trayInfo
    implicitWidth: trayRow.implicitWidth
    implicitHeight: 40
    visible: SystemTray.items.values.length > 0

    Row {
        id: trayRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 4

        Repeater {
            model: SystemTray.items.values
            delegate: Rectangle {
                required property var modelData
                property var trayItem: modelData

                width: 20
                height: 20
                radius: 4
                color: mouseArea.containsMouse ? Theme.surface1 : "transparent"

                Image {
                    anchors.centerIn: parent
                    width: 20
                    height: 20
                    source: trayItem.icon
                    asynchronous: true
                    smooth: true
                    visible: status === Image.Ready
                }

                MouseArea {
                    id: mouseArea
                    anchors.fill: parent
                    hoverEnabled: true
                    acceptedButtons: Qt.LeftButton | Qt.MiddleButton | Qt.RightButton
                    cursorShape: Qt.PointingHandCursor
                    onClicked: mouse => {
                        if (mouse.button === Qt.LeftButton) trayItem.activate()
                        else if (mouse.button === Qt.MiddleButton) trayItem.secondaryActivate()
                        else if (mouse.button === Qt.RightButton && trayItem.hasMenu) {
                            var pt = mapToItem(trayInfo, width / 2, height)
                            trayItem.display(trayInfo.Window.window, Math.round(pt.x), Math.round(pt.y))
                        }
                    }
                }
            }
        }
    }
}
