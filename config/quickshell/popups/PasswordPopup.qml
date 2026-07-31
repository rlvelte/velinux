import QtQuick
import QtQuick.Controls
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import qs.globals
import "../util/utils.js" as Utils

PanelWindow {
    id: passwordPopup
    color: "transparent"

    anchors { top: true; bottom: true; left: true; right: true }
    exclusiveZone: 0

    focusable: true
    WlrLayershell.keyboardFocus: shown ? WlrKeyboardFocus.Exclusive : WlrKeyboardFocus.None
    WlrLayershell.layer: WlrLayer.Overlay

    mask: Region { item: contentRect }

    property bool shown: false
    property bool animatingOut: false
    property string promptDir: ""
    property string promptText: "Password:"
    property string resultPath: ""

    property int screenWidth: screen ? screen.width
        : (Quickshell.screens.length > 0 ? Quickshell.screens[0].width : 1920)

    visible: shown || animatingOut

    IpcHandler {
        target: "su"
        function vlxOpen(dirPath: string): void {
            passwordPopup.startPrompt(dirPath)
        }
    }

    function startPrompt(dirPath) {
        passwordPopup.reset()
        passwordPopup.promptDir = dirPath
        passwordPopup.resultPath = dirPath + "/result"
        promptReader.command = ["cat", dirPath + "/prompt.json"]
        promptReader.running = true
        passwordPopup.shown = true
    }

    function closePrompt() {
        passwordPopup.promptDir = ""
        passwordPopup.shown = false
        passwordPopup.animatingOut = true
        hideAnim.restart()
    }

    function reset() {
        passwordField.text = ""
        passwordPopup.promptText = "Password:"
        passwordPopup.resultPath = ""
    }

    onShownChanged: {
        if (shown) {
            animatingOut = false
            hideTimer.stop()
            hideAnim.stop()
            contentTranslate.y = -Dimensions.overlaySlideOffset
            dropTimer.restart()
        }
    }

    Timer {
        id: dropTimer
        interval: Anims.dropDelay
        repeat: false
        onTriggered: {
            showAnim.start()
            passwordField.forceActiveFocus()
        }
    }

    Timer {
        id: hideTimer
        interval: Anims.duration
        repeat: false
        onTriggered: {
            passwordPopup.animatingOut = false
            passwordPopup.reset()
        }
    }

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

    Process {
        id: promptReader
        stdout: StdioCollector {
            onStreamFinished: {
                try {
                    var data = JSON.parse(text)
                    passwordPopup.promptText = data.prompt || "Password:"
                } catch (e) {
                    passwordPopup.promptText = "Password:"
                }
            }
        }
    }

    Process { id: resultWriter }

    function submitPassword() {
        var pwd = passwordField.text
        if (pwd === "") return
        resultWriter.command = ["sh", "-c", "printf '%s' " + Utils.escapeShell(pwd) + " > " + Utils.escapeShell(passwordPopup.resultPath)]
        resultWriter.running = true
        passwordPopup.closePrompt()
    }

    function cancelPrompt() {
        resultWriter.command = ["sh", "-c", "printf '' > " + Utils.escapeShell(passwordPopup.resultPath)]
        resultWriter.running = true
        passwordPopup.closePrompt()
    }

    MouseArea {
        anchors.fill: parent
        z: -1
        onClicked: passwordPopup.cancelPrompt()
    }

    Rectangle {
        id: contentRect
        width: Math.round(passwordPopup.screenWidth / 3)
        height: 64
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        color: Theme.base
        radius: Dimensions.overlayRadius
        border.color: Theme.surface1; border.width: 1

        transform: Translate {
            id: contentTranslate
            y: -Dimensions.overlaySlideOffset
        }

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        Row {
            anchors.left: parent.left
            anchors.right: parent.right
            anchors.verticalCenter: parent.verticalCenter
            anchors.leftMargin: 20
            anchors.rightMargin: 20
            spacing: 12

            Text {
                id: promptLabel
                text: passwordPopup.promptText
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSizeLarge
                color: Theme.text
                elide: Text.ElideRight
                anchors.verticalCenter: parent.verticalCenter
            }

            TextField {
                id: passwordField
                width: parent.width - promptLabel.width - okButton.width
                    - parent.spacing
                height: 38
                padding: 10
                focus: true
                echoMode: TextInput.Password
                inputMethodHints: Qt.ImhSensitiveData
                placeholderText: "Password"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                color: Theme.text
                placeholderTextColor: Theme.muted
                anchors.verticalCenter: parent.verticalCenter

                background: Rectangle {
                    color: Qt.rgba(Theme.surface0.r, Theme.surface0.g, Theme.surface0.b, 0.6)
                    radius: Dimensions.inputRadius
                    border.color: passwordField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: passwordField.activeFocus ? 2 : 1

                    Behavior on border.color {
                        ColorAnimation { duration: Anims.duration; easing.type: Easing.OutCubic }
                    }
                    Behavior on border.width {
                        NumberAnimation { duration: Anims.duration; easing.type: Easing.OutCubic }
                    }
                }

                onAccepted: passwordPopup.submitPassword()
                Keys.onEscapePressed: passwordPopup.cancelPrompt()
            }

            Rectangle {
                id: okButton
                width: 80; height: 38
                radius: Dimensions.inputRadius
                color: okMouse.containsMouse
                    ? Qt.lighter(Theme.primary, 1.08)
                    : Theme.primary
                opacity: passwordField.text !== "" ? 1.0 : 0.35
                anchors.verticalCenter: parent.verticalCenter

                Behavior on color {
                    ColorAnimation { duration: Anims.durationFast; easing.type: Easing.OutCubic }
                }
                Behavior on opacity {
                    NumberAnimation { duration: Anims.durationFast; easing.type: Easing.OutCubic }
                }

                Text {
                    anchors.centerIn: parent
                    text: "OK"
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeSmall
                    font.bold: true
                    color: Theme.textOnPrimary
                }

                MouseArea {
                    id: okMouse
                    anchors.fill: parent
                    hoverEnabled: true
                    cursorShape: passwordField.text !== ""
                        ? Qt.PointingHandCursor
                        : Qt.ArrowCursor
                    onClicked: {
                        if (passwordField.text !== "")
                            passwordPopup.submitPassword()
                    }
                }
            }
        }
    }
}
