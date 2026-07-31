import QtQuick
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import qs.globals

PanelWindow {
    id: progressPopup
    color: "transparent"

    anchors { top: true; bottom: true; left: true; right: true }
    exclusiveZone: 0

    focusable: false
    WlrLayershell.keyboardFocus: WlrKeyboardFocus.None
    WlrLayershell.layer: WlrLayer.Overlay

    mask: Region {
        item: contentRect
    }

    property bool shown: false
    property bool animatingOut: false

    property int screenWidth: screen ? screen.width : (Quickshell.screens.length > 0 ? Quickshell.screens[0].width : 1920)

    property string progressDir: ""
    property string progressLabel: ""
    property real progressValue: 0.0
    property bool isIndefinite: false

    property int spinnerFrame: 0
    readonly property var spinnerChars: ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"]

    visible: shown || animatingOut

    IpcHandler {
        target: "progress"
        function vlxOpen(dirPath: string): void {
            progressPopup.startProgress(dirPath)
        }
    }

    function startProgress(dirPath) {
        progressPopup.reset()
        progressPopup.progressDir = dirPath
        progressPopup.shown = true
        stateReader.command = ["cat", dirPath + "/state.json"]
        stateReader.running = true
    }

    function closeProgress() {
        progressPopup.progressDir = ""
        pollTimer.running = false
        progressPopup.shown = false
        progressPopup.animatingOut = true
        hideAnim.restart()
    }

    function reset() {
        progressPopup.progressLabel = ""
        progressPopup.progressValue = 0.0
        progressPopup.isIndefinite = false
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
            pollTimer.running = true
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

    Timer {
        id: hideTimer
        interval: Anims.duration
        repeat: false
        onTriggered: {
            progressPopup.animatingOut = false
            spinnerTimer.running = false
            progressPopup.reset()
        }
    }

    Timer {
        id: pollTimer
        interval: 100
        repeat: true
        running: false
        onTriggered: {
            if (progressPopup.progressDir === "") return
            stateReader.command = ["cat", progressPopup.progressDir + "/state.json"]
            stateReader.running = true
        }
    }

    Process {
        id: stateReader
        stdout: StdioCollector {
            onStreamFinished: {
                try {
                    var state = JSON.parse(text)
                    if (state.done === true) {
                        progressPopup.closeProgress()
                        return
                    }
                    progressPopup.progressLabel = state.label || ""
                    progressPopup.isIndefinite = state.indefinite === true
                    progressPopup.progressValue = state.progress || 0.0
                } catch (e) {
                    console.error("Failed to parse progress state:", e)
                }
            }
        }
    }

    Timer {
        id: spinnerTimer
        interval: 120
        repeat: true
        running: shown && isIndefinite
        onTriggered: spinnerFrame = (spinnerFrame + 1) % spinnerChars.length
    }

    Rectangle {
        id: contentRect
        width: Math.round(progressPopup.screenWidth / 3)
        height: Dimensions.listItemHeight
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        color: Theme.base
        radius: Dimensions.overlayRadius
        border.color: Theme.surface1; border.width: 1

        transform: Translate {
            id: contentTranslate
            y: -Dimensions.overlaySlideOffset
        }

        Row {
            anchors.fill: parent
            anchors.leftMargin: 20
            anchors.rightMargin: 20
            spacing: 12

            Item {
                width: 32
                height: parent.height

                Text {
                    visible: isIndefinite
                    anchors.centerIn: parent
                    text: spinnerChars[spinnerFrame]
                    font.pixelSize: 28
                    color: Theme.primary
                }

                Text {
                    visible: !isIndefinite
                    anchors.centerIn: parent
                    text: "●"
                    font.pixelSize: Theme.fontSizeHeading
                    color: progressValue >= 1.0 ? Theme.success : Theme.primary
                }
            }

            Text {
                id: labelText
                width: isIndefinite
                    ? parent.width - 32 - parent.spacing
                    : Math.max(80, parent.width * 0.30 - 32 - parent.spacing)
                anchors.verticalCenter: parent.verticalCenter
                text: progressPopup.progressLabel
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSizeLarge
                color: Theme.text
                elide: Text.ElideRight
                verticalAlignment: Text.AlignVCenter
            }

            Item {
                visible: !isIndefinite
                height: parent.height
                width: parent.width - labelText.width - 32 - parent.spacing * 2
                anchors.verticalCenter: parent.verticalCenter

                Rectangle {
                    width: parent.width
                    height: 20
                    anchors.verticalCenter: parent.verticalCenter
                        radius: Dimensions.overlayRadius
                    color: Theme.surface0
                    clip: true

                    Rectangle {
                        width: parent.width * progressPopup.progressValue
                        height: parent.height
                    radius: Dimensions.overlayRadius
                        color: progressPopup.progressValue >= 1.0 ? Theme.success : Theme.primary

                        Behavior on width {
                            NumberAnimation { duration: Anims.durationFast; easing.type: Easing.OutCubic }
                        }
                    }

                    Text {
                        anchors.centerIn: parent
                        text: Math.round(progressPopup.progressValue * 100) + "%"
                        font.family: Theme.fontNameMono
                        font.pixelSize: Theme.fontSizeSmall
                        font.bold: true
                        color: progressPopup.progressValue > 0.5 ? Theme.textOnPrimary : Theme.text
                    }
                }
            }
        }
    }
}
