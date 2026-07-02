import QtQuick
import Quickshell
import Quickshell.Io
import qs.services

FloatingWindow {
    id: powerMenu
    title: "vlx-power"
    color: "transparent"
    implicitWidth: 640
    implicitHeight: 560

    property bool shown: false
    property int selected: 0

    visible: shown

    IpcHandler {
        target: "power"
        function toggle(): void { powerMenu.shown = !powerMenu.shown }
        function open(): void { powerMenu.shown = true }
        function close(): void { powerMenu.shown = false }
    }

    onShownChanged: {
        if (shown) {
            selected = 0
            powerMenu.focus = true
        }
    }

    function hide() {
        shown = false
    }

    function executeCommand(cmd) {
        if (cmd[1] === "$XDG_SESSION_ID") {
            cmd[1] = Quickshell.env("XDG_SESSION_ID")
        }
        processRunner.command = cmd
        processRunner.running = true
        powerMenu.hide()
    }

    function launch(index) {
        var model = [
            { command: ["quickshell", "ipc", "call", "lock", "lock"] },
            { command: ["loginctl", "terminate-session", "$XDG_SESSION_ID"] },
            { command: ["systemctl", "poweroff"] },
            { command: ["systemctl", "reboot"] }
        ]
        executeCommand(model[index].command)
    }

    Process {
        id: processRunner
    }

    Rectangle {
        anchors.fill: parent
        anchors.margins: 16
        color: Theme.base
        radius: 12
        border.color: Theme.surface1
        border.width: 1
        focus: true

        Keys.onEscapePressed: powerMenu.hide()
        Keys.onDownPressed: {
            powerMenu.selected = Math.min(powerMenu.selected + 1, 3)
        }
        Keys.onUpPressed: {
            powerMenu.selected = Math.max(powerMenu.selected - 1, 0)
        }
        Keys.onReturnPressed: powerMenu.launch(powerMenu.selected)

        Column {
            anchors.fill: parent
            anchors.margins: 16
            spacing: 12

            Text {
                text: "Power"
                font.family: Theme.fontName
                font.pixelSize: 22
                font.weight: Font.Bold
                color: Theme.text
                anchors.horizontalCenter: parent.horizontalCenter
            }

            ListView {
                id: itemList
                width: parent.width
                height: parent.height - 60
                clip: true
                model: [
                    { icon: "", label: "Lock" },
                    { icon: "", label: "Logout" },
                    { icon: "", label: "Shutdown" },
                    { icon: "", label: "Reboot" }
                ]
                currentIndex: powerMenu.selected
                highlightFollowsCurrentItem: true
                boundsBehavior: Flickable.StopAtBounds
                spacing: 4

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: itemList.width
                    height: 48
                    radius: 8
                    color: index === powerMenu.selected
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                        : mouseArea.containsMouse
                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                            : "transparent"

                    Rectangle {
                        visible: index === powerMenu.selected
                        anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                        width: 3
                        radius: 1.5
                        color: Theme.primary
                    }

                    Row {
                        anchors.fill: parent
                        anchors.margins: 8
                        spacing: 12
                        anchors.verticalCenter: parent.verticalCenter

                        Text {
                            text: modelData.icon
                            font.family: Theme.fontName
                            font.pixelSize: 18
                            color: Theme.primary
                            width: 24
                            horizontalAlignment: Text.AlignHCenter
                        }

                        Text {
                            text: modelData.label
                            font.family: Theme.fontName
                            font.pixelSize: 20
                            color: index === powerMenu.selected ? Theme.text : Theme.subtext
                        }
                    }

                    MouseArea {
                        id: mouseArea
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onEntered: powerMenu.selected = index
                        onClicked: powerMenu.launch(index)
                    }
                }
            }
        }

        MouseArea {
            anchors.fill: parent
            z: -1
            onClicked: powerMenu.hide()
        }
    }
}
