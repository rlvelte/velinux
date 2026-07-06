import QtQuick
import QtQuick.Controls
import Quickshell
import Quickshell.Io
import Quickshell.Wayland
import Quickshell.Widgets
import qs.services

PanelWindow {
    id: picker
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

    property var pickerItems: []
    property string pickerResultPath: ""

    visible: shown || animatingOut

    IpcHandler {
        target: "singlepicker"
        function vlxOpen(filePath: string, resultPath: string): void {
            picker.startPicker(filePath, resultPath)
        }
    }

    onShownChanged: {
        if (shown) {
            searchField.text = ""
            selected = 0
            canClose = false
            contentRect.opacity = 0
            dropTimer.restart()
            closeGuard.restart()
        }
    }

    function startPicker(filePath, resultPath) {
        picker.pickerResultPath = resultPath
        picker.selected = 0
        itemsReader.command = ["cat", filePath]
        itemsReader.running = true
        picker.shown = true
    }

    function cancelPicker() {
        picker.writeResult(null)
        picker.resetPicker()
    }

    function resetPicker() {
        picker.pickerResultPath = ""
        picker.pickerItems = []
    }

    Timer {
        id: dropTimer
        interval: 50
        repeat: false
        onTriggered: {
            showAnim.start()
            searchField.forceActiveFocus()
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

    Process {
        id: itemsReader
        stdout: StdioCollector {
            onStreamFinished: {
                try {
                    picker.pickerItems = JSON.parse(text)
                } catch (e) {
                    console.error("Failed to parse picker items:", e)
                    picker.hide()
                }
            }
        }
    }

    Process {
        id: resultWriter
    }

    function writeResult(item) {
        let json = JSON.stringify(item);
        resultWriter.command = ["sh", "-c", "printf '%s' " + picker.escapeShell(json) + " > " + picker.escapeShell(picker.pickerResultPath)]
        resultWriter.running = true
    }

    function escapeShell(str) {
        return "'" + str.replace(/'/g, "'\\''") + "'"
    }

    function hide() {
        if (!shown || animatingOut) return
        animatingOut = true
        shown = false
        hideAnim.start()
    }

    function launch(entry) {
        if (!entry) return
        picker.writeResult(entry._source)
        picker.resetPicker()
        hide()
    }

    NumberAnimation {
        id: showAnim
        target: contentRect
        property: "opacity"
        from: 0
        to: 1
        duration: 200
        easing.type: Easing.OutCubic
    }

    NumberAnimation {
        id: hideAnim
        target: contentRect
        property: "opacity"
        from: 1
        to: 0
        duration: 200
        easing.type: Easing.InCubic
        onFinished: hideTimer.restart()
    }

    MouseArea {
        anchors.fill: parent
        z: -1
        enabled: picker.canClose
        onClicked: {
            picker.cancelPicker()
            picker.hide()
        }
    }

    ScriptModel {
        id: filtered
        values: {
            var all = picker.pickerItems.map(function(e) {
                return {
                    icon: e.icon,
                    name: e.header,
                    comment: e.description || "",
                    _source: e,
                    keywords: []
                }
            })

            var q = searchField.text.trim().toLowerCase()
            if (q === "") return all

            return all.filter(function(e) {
                var name = (e.name || "").toLowerCase()
                var comment = (e.comment || "").toLowerCase()
                return name.indexOf(q) !== -1 || comment.indexOf(q) !== -1
            })
        }
    }

    Rectangle {
        id: contentRect
        width: 640
        height: 560
        anchors.centerIn: parent
        radius: 12
        color: Theme.base
        border.color: Theme.surface1
        border.width: 1

        opacity: 0

        MouseArea {
            anchors.fill: parent
            propagateComposedEvents: false
            onClicked: {}
        }

        Column {
            anchors.fill: parent
            anchors.margins: 16
            spacing: 12

            TextField {
                id: searchField
                focus: true
                width: parent.width
                height: 52
                leftPadding: 16
                placeholderText: "Search..."
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSizeHeading
                color: Theme.text
                placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0
                    radius: 8
                    border.color: searchField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: searchField.activeFocus ? 2 : 1
                }

                Keys.onDownPressed: {
                    picker.selected = Math.min(picker.selected + 1, filtered.values.length - 1)
                }
                Keys.onUpPressed: {
                    picker.selected = Math.max(picker.selected - 1, 0)
                }
                Keys.onReturnPressed: picker.launch(filtered.values[picker.selected])
                Keys.onEscapePressed: {
                    picker.cancelPicker()
                    picker.hide()
                }
                onTextChanged: {
                    if (picker.selected >= filtered.values.length)
                        picker.selected = Math.max(0, filtered.values.length - 1)
                }
            }

            ListView {
                id: appList
                width: parent.width
                height: parent.height - searchField.height - 12
                clip: true
                model: filtered
                currentIndex: picker.selected
                highlightMoveDuration: 100
                highlightFollowsCurrentItem: true
                boundsBehavior: Flickable.StopAtBounds
                spacing: 4

                highlight: Rectangle {
                    width: appList.width
                    height: 64
                    radius: 8
                    color: Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)

                    Rectangle {
                        anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                        width: 3
                        radius: 1.5
                        color: Theme.primary
                    }
                }

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: appList.width
                    height: 64
                    radius: 8
                    color: "transparent"

                    Row {
                        anchors.fill: parent
                        anchors.margins: 8
                        spacing: 12

                        Rectangle {
                            id: iconRect
                            width: 48
                            height: 48
                            radius: 10
                            color: Theme.surface0
                            anchors.verticalCenter: parent.verticalCenter

                            property bool isFileIcon: modelData.icon && modelData.icon.startsWith("/")

                            Image {
                                id: iconImg
                                anchors.centerIn: parent
                                source: iconRect.isFileIcon ? modelData.icon : ""
                                width: 28
                                height: 28
                                fillMode: Image.PreserveAspectFit
                                visible: iconRect.isFileIcon && status === Image.Ready
                            }

                            Text {
                                anchors.centerIn: parent
                                text: iconRect.isFileIcon || !modelData.icon || /^[\w.-]+$/.test(modelData.icon)
                                    ? modelData.name.charAt(0).toUpperCase()
                                    : modelData.icon
                                font.family: Theme.fontNameMono
                                font.pixelSize: iconRect.isFileIcon ? 20 : 28
                                color: Theme.subtext
                                visible: !iconRect.isFileIcon || iconImg.status !== Image.Ready
                            }
                        }

                        Column {
                            anchors.verticalCenter: parent.verticalCenter
                            spacing: 2

                            Text {
                                text: modelData.name
                                font.family: Theme.fontName
                                font.pixelSize: Theme.fontSizeLarge
                                color: index === picker.selected ? Theme.text : Theme.subtext
                                elide: Text.ElideRight
                                width: appList.width - 110
                            }

                            Text {
                                visible: !!modelData.comment
                                text: modelData.comment
                                font.family: Theme.fontName
                                font.pixelSize: Theme.fontSize
                                color: Theme.muted
                                elide: Text.ElideRight
                                width: appList.width - 110
                            }
                        }
                    }

                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: picker.launch(modelData)
                    }
                }
            }
        }
    }
}
