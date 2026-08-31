import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import qs.globals

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

    readonly property var actions: [
        { label: "Logout",   command: ["swaymsg", "exit"] },
        { label: "Suspend",  command: elevated(["systemctl", "suspend"]) },
        { label: "Shutdown", command: elevated(["systemctl", "poweroff"]) },
        { label: "Reboot",   command: elevated(["systemctl", "reboot"]) },
        { label: "BIOS",     command: elevated(["systemctl", "reboot", "--firmware-setup"]) }
    ]

    visible: shown || animatingOut

    IpcHandler {
        target: "power"
        function toggle(): void { if (powerMenu.shown) powerMenu.hide(); else powerMenu.shown = true }
        function open(): void { powerMenu.shown = true }
        function close(): void { powerMenu.hide() }
    }

    onShownChanged: {
        if (shown) {
            animatingOut = false
            hideTimer.stop()
            hideAnim.stop()
            selected = 0
            canClose = false
            contentTranslate.y = -Dimensions.overlaySlideOffset
            dropTimer.restart()
            closeGuard.restart()
        }
    }

    Timer {
        id: dropTimer
        interval: Anims.dropDelay
        repeat: false
        onTriggered: {
            showAnim.start()
            contentRect.forceActiveFocus()
        }
    }

    Timer {
        id: closeGuard
        interval: Anims.duration
        repeat: false
        onTriggered: canClose = true
    }

    Timer {
        id: hideTimer
        interval: Anims.duration
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
        for (var i = 0; i < cmd.length; i++) {
            if (cmd[i] === "$XDG_SESSION_ID")
                cmd[i] = Quickshell.env("XDG_SESSION_ID")
        }
        processRunner.command = cmd
        processRunner.running = true
        powerMenu.hide()
    }

    function elevated(cmd) {
        var joined = cmd.map(function(arg) { return "'" + arg.replace(/'/g, "'\\''") + "'" }).join(" ")
        return ["sh", "-c", "SUDO_ASKPASS=$(command -v vlxpass) exec sudo -A " + joined]
    }

    Process { id: processRunner }

    NumberAnimation {
        id: showAnim
        target: contentTranslate; property: "y"
        from: -Dimensions.overlaySlideOffset; to: 0
        duration: Anims.duration; easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentTranslate; property: "y"
        from: 0; to: -Dimensions.overlaySlideOffset
        duration: Anims.duration; easing.type: Easing.InCubic
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
        color: Theme.base
        radius: Dimensions.overlayRadius
        border.color: Theme.surface1; border.width: 1
        focus: true

        transform: Translate {
            id: contentTranslate
            y: -Dimensions.overlaySlideOffset
        }

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        Keys.onLeftPressed:  selected = Math.max(selected - 1, 0)
        Keys.onRightPressed: selected = Math.min(selected + 1, powerMenu.actions.length - 1)
        Keys.onReturnPressed: executeCommand(powerMenu.actions[selected].command)
        Keys.onEscapePressed: powerMenu.hide()

        Row {
            anchors.fill: parent
            anchors.margins: 12
            spacing: 4

            Repeater {
                model: powerMenu.actions

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: (parent.width - parent.spacing * (powerMenu.actions.length - 1)) / powerMenu.actions.length
                    height: parent.height
                    radius: Dimensions.inputRadius
                    color: index === powerMenu.selected
                        ? Theme.primarySelected
                        : mouseArea.containsMouse ? Theme.primaryHovered : "transparent"

                    Rectangle {
                        visible: index === powerMenu.selected
                        anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                        width: 3; radius: 1.5; color: Theme.primary
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
                        onClicked: powerMenu.executeCommand(modelData.command)
                    }
                }
            }
        }
    }
}
