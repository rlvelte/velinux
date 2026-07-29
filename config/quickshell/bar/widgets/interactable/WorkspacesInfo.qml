import QtQuick
import Quickshell
import Quickshell.I3
import qs.globals

Item {
    id: workspaces
    implicitWidth: wsRow.implicitWidth
    implicitHeight: Dimensions.barHeight

    property var wsModel: I3.workspaces
    property var wsValues: I3.workspaces.values

    readonly property int totalWorkspaces: 5

    function isWorkspaceFocused(num: int): bool {
        var values = wsValues;
        for (var i = 0; i < values.length; i++) {
            var ws = values[i];
            if (ws.number === num && ws.focused)
                return true;
        }
        return false;
    }

    function isWorkspaceActive(num: int): bool {
        var values = wsValues;
        for (var i = 0; i < values.length; i++) {
            var ws = values[i];
            if (ws.number === num && ws.active)
                return true;
        }
        return false;
    }

    function activateWorkspace(num: int): void {
        var values = wsValues;
        for (var i = 0; i < values.length; i++) {
            var ws = values[i];
            if (ws.number === num) {
                ws.activate();
                return;
            }
        }
        I3.dispatch("workspace number " + num);
    }

    Row {
        id: wsRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 4

        Repeater {
            model: totalWorkspaces
            delegate: Rectangle {
                property int num: index + 1
                property bool focused: workspaces.isWorkspaceFocused(num)
                property bool active: workspaces.isWorkspaceActive(num)

                width: wsText.implicitWidth + 10
                height: 20
                color: "transparent"

                Text {
                    id: wsText
                    anchors.centerIn: parent
                    text: num
                    color: focused ? Theme.primary : (active ? Theme.subtext : Qt.darker(Theme.subtext, 2))
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeHeading
                    font.weight: focused ? Font.Bold : Font.Medium
                }

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: workspaces.activateWorkspace(num)
                }
            }
        }
    }
}
