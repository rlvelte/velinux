import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import qs.services

PanelWindow {
    id: powerMenu
    color: "transparent"

    anchors { top: true; bottom: true; left: true; right: true }
    exclusiveZone: 0

    focusable: true
    WlrLayershell.keyboardFocus: shown ? WlrKeyboardFocus.Exclusive : WlrKeyboardFocus.None
    WlrLayershell.layer: WlrLayer.Overlay

    property bool shown: false
    property bool animatingOut: false
    property bool canClose: false
    property int selected: 0

    property int screenWidth: screen ? screen.width : (Quickshell.screens.length > 0 ? Quickshell.screens[0].width : 1920)

    visible: shown || animatingOut

    IpcHandler {
        target: "power"
        function toggle(): void { if (powerMenu.shown) powerMenu.hide(); else powerMenu.shown = true }
        function open(): void { powerMenu.shown = true }
        function close(): void { powerMenu.hide() }
    }

    onShownChanged: {
        if (shown) {
            selected = 0
            canClose = false
            contentTranslate.y = -80
            dropTimer.restart()
            closeGuard.restart()
        }
    }

    Timer {
        id: dropTimer
        interval: 50
        repeat: false
        onTriggered: {
            showAnim.start()
            contentRect.forceActiveFocus()
        }
    }

    Timer {
        id: closeGuard
        interval: 200
        repeat: false
        onTriggered: canClose = true
    }

    Timer {
        id: hideTimer
        interval: 200
        repeat: false
        onTriggered: animatingOut = false
    }

    function hide() {
        if (!shown || animatingOut) return
        animatingOut = true
        shown = false
        hideAnim.start()
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
            { command: ["systemctl", "reboot"] },
            { command: ["systemctl", "reboot", "--firmware-setup"] }
        ]
        executeCommand(model[index].command)
    }

    Process {
        id: processRunner
    }

    NumberAnimation {
        id: showAnim
        target: contentTranslate
        property: "y"
        from: -80
        to: 0
        duration: 200
        easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentTranslate
        property: "y"
        from: 0
        to: -80
        duration: 200
        easing.type: Easing.InCubic
        onFinished: hideTimer.restart()
    }

    MouseArea {
        anchors.fill: parent
        z: -1
        enabled: powerMenu.canClose
        onClicked: powerMenu.hide()
    }

    Rectangle {
        id: contentRect
        width: Math.round(powerMenu.screenWidth / 3)
        height: 80
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        anchors.topMargin: 0
        color: Theme.base
        radius: 12
        border.color: Theme.surface1
        border.width: 1
        focus: true

        transform: Translate {
            id: contentTranslate
            y: -80
        }

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        Keys.onLeftPressed: {
            selected = Math.max(selected - 1, 0)
        }
        Keys.onRightPressed: {
            selected = Math.min(selected + 1, 4)
        }
        Keys.onReturnPressed: {
            launch(selected)
        }
        Keys.onEscapePressed: {
            powerMenu.hide()
        }

        Row {
            anchors.fill: parent
            anchors.margins: 12
            spacing: 4

            Repeater {
                model: [
                    { label: "Lock" },
                    { label: "Logout" },
                    { label: "Shutdown" },
                    { label: "Reboot" },
                    { label: "BIOS" }
                ]

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: (parent.width - parent.spacing * 4) / 5
                    height: parent.height
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

                    Text {
                        anchors.centerIn: parent
                        text: modelData.label
                        font.family: Theme.fontName
                        font.pixelSize: Theme.fontSizeHeading
                        color: index === powerMenu.selected ? Theme.text : Theme.subtext
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
    }
}