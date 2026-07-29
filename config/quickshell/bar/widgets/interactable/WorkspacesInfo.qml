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

    function getWorkspace(num: int): var {
        var values = wsValues;
        for (var i = 0; i < values.length; i++) {
            if (values[i].number === num)
                return values[i];
        }
        return null;
    }

    function isWorkspaceFocused(num: int): bool {
        var ws = getWorkspace(num);
        return ws ? ws.focused : false;
    }

    function hasWindows(num: int): bool {
        var ws = getWorkspace(num);
        if (!ws || !ws.lastIpcObject) return false;
        var obj = ws.lastIpcObject;
        return (obj.nodes && obj.nodes.length > 0) || (obj.floating_nodes && obj.floating_nodes.length > 0);
    }

    function activateWorkspace(num: int): void {
        var ws = getWorkspace(num);
        if (ws) {
            ws.activate();
        } else {
            I3.dispatch("workspace number " + num);
        }
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
                property bool populated: workspaces.hasWindows(num)

                width: wsText.implicitWidth + 10
                height: 20
                color: "transparent"

                Text {
                    id: wsText
                    anchors.centerIn: parent
                    text: num
                    color: focused ? Theme.primary : (populated ? Theme.text : Qt.darker(Theme.subtext, 2))
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeHeading
                    font.weight: focused ? Font.Bold : (populated ? Font.Medium : Font.Normal)
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
