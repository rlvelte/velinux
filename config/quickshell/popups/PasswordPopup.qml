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

    property bool shown: false
    property bool animatingOut: false
    property string promptDir: ""
    property string promptText: "Password:"
    property string resultPath: ""

    visible: shown || animatingOut

    IpcHandler {
        target: "su"
        function vlxOpen(dirPath: string): void {
            passwordPopup.startPrompt(dirPath)
        }
    }

    function startPrompt(dirPath) {
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
            contentRect.opacity = 0
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
        target: contentRect; property: "opacity"
        from: 0; to: 1; duration: Anims.duration; easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentRect; property: "opacity"
        from: 1; to: 0; duration: Anims.duration; easing.type: Easing.InCubic
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
        width: 400
        height: 200
        anchors.centerIn: parent
        radius: Dimensions.overlayRadius
        color: Theme.base
        border.color: Theme.surface1; border.width: 1
        opacity: 0

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        Column {
            anchors.fill: parent
            anchors.margins: 24
            spacing: 16

            Text {
                text: passwordPopup.promptText
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSizeHeading
                color: Theme.text
                width: parent.width
                horizontalAlignment: Text.AlignLeft
            }

            TextField {
                id: passwordField
                width: parent.width
                height: 48
                padding: 12
                focus: true
                echoMode: TextInput.Password
                inputMethodHints: Qt.ImhSensitiveData
                placeholderText: "Password"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                color: Theme.text
                placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0
                    radius: Dimensions.inputRadius
                    border.color: passwordField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: passwordField.activeFocus ? 2 : 1
                }

                onAccepted: passwordPopup.submitPassword()
                Keys.onEscapePressed: passwordPopup.cancelPrompt()
            }

            Row {
                width: parent.width
                spacing: 12
                layoutDirection: Qt.RightToLeft

                Button {
                    text: "OK"
                    height: 40
                    padding: 16
                    enabled: passwordField.text !== ""
                    onClicked: passwordPopup.submitPassword()
                }

                Button {
                    text: "Cancel"
                    height: 40
                    padding: 16
                    onClicked: passwordPopup.cancelPrompt()
                }
            }
        }
    }
}
