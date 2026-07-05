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

    // Set of selected item indices
    property var selectedIndices: new Set()

    visible: shown || animatingOut

    IpcHandler {
        target: "multipicker"
        function vlxOpen(filePath: string, resultPath: string): void {
            picker.startPicker(filePath, resultPath)
        }
    }

    onShownChanged: {
        if (shown) {
            searchField.text = ""
            selected = 0
            selectedIndices = new Set()
            canClose = false
            contentRect.opacity = 0
            dropTimer.restart()
            closeGuard.restart()
        }
    }

    function startPicker(filePath, resultPath) {
        picker.pickerResultPath = resultPath
        picker.selected = 0
        picker.selectedIndices = new Set()
        itemsReader.command = ["cat", filePath]
        itemsReader.running = true
        picker.shown = true
    }

    function cancelPicker() {
        picker.writeResult([])
        picker.resetPicker()
    }

    function resetPicker() {
        picker.pickerResultPath = ""
        picker.pickerItems = []
        picker.selectedIndices = new Set()
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

    function writeResult(items) {
        let json = JSON.stringify(items);
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

    function toggleSelection(entry) {
        if (!entry) return
        var src = entry._source
        var idx = picker.pickerItems.indexOf(src)
        if (idx === -1) return
        var s = new Set(picker.selectedIndices)
        if (s.has(idx)) {
            s.delete(idx)
        } else {
            s.add(idx)
        }
        picker.selectedIndices = s
    }

    function confirmSelection() {
        var items = []
        picker.selectedIndices.forEach(function(idx) {
            if (idx >= 0 && idx < picker.pickerItems.length) {
                items.push(picker.pickerItems[idx])
            }
        })
        picker.writeResult(items)
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
            var all = picker.pickerItems.map(function(e, i) {
                return {
                    icon: e.icon,
                    name: e.header,
                    comment: e.description || "",
                    _source: e,
                    _index: i,
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
                Keys.onReturnPressed: picker.toggleSelection(filtered.values[picker.selected])
                Keys.onEscapePressed: {
                    picker.confirmSelection()
                }
                onTextChanged: {
                    if (picker.selected >= filtered.values.length)
                        picker.selected = Math.max(0, filtered.values.length - 1)
                }
            }

            Row {
                width: parent.width
                height: 36
                spacing: 8

                Text {
                    text: picker.selectedIndices.size + " selected"
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSize
                    color: Theme.muted
                    anchors.verticalCenter: parent.verticalCenter
                }

                Item { width: parent.width - 200; height: 1 }

                Text {
                    text: "Return = toggle · Esc = confirm"
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSizeSmall
                    color: Theme.muted
                    anchors.verticalCenter: parent.verticalCenter
                }
            }

            ListView {
                id: appList
                width: parent.width
                height: parent.height - searchField.height - 12 - 36
                clip: true
                model: filtered
                currentIndex: picker.selected
                highlightFollowsCurrentItem: true
                boundsBehavior: Flickable.StopAtBounds
                spacing: 4

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: appList.width
                    height: 64
                    radius: 8
                    color: index === picker.selected
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                        : mouseArea.containsMouse
                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                            : "transparent"

                    Rectangle {
                        visible: index === picker.selected
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

                        Rectangle {
                            width: 48
                            height: 48
                            radius: 10
                            color: Theme.surface0
                            anchors.verticalCenter: parent.verticalCenter

                            Text {
                                anchors.centerIn: parent
                                text: modelData.icon
                                font.pixelSize: 22
                                color: Theme.subtext
                            }
                        }

                        Column {
                            anchors.verticalCenter: parent.verticalCenter
                            spacing: 2
                            width: appList.width - 110 - 32

                            Text {
                                text: modelData.name
                                font.family: Theme.fontName
                                font.pixelSize: Theme.fontSizeLarge
                                color: index === picker.selected ? Theme.text : Theme.subtext
                                elide: Text.ElideRight
                                width: parent.width
                            }

                            Text {
                                visible: !!modelData.comment
                                text: modelData.comment
                                font.family: Theme.fontName
                                font.pixelSize: Theme.fontSize
                                color: Theme.muted
                                elide: Text.ElideRight
                                width: parent.width
                            }
                        }

                        Rectangle {
                            width: 24
                            height: 24
                            radius: 4
                            color: picker.selectedIndices.has(modelData._index)
                                ? Theme.primary
                                : Theme.surface0
                            border.color: Theme.surface1
                            border.width: 1
                            anchors.verticalCenter: parent.verticalCenter

                            Text {
                                anchors.centerIn: parent
                                text: "✓"
                                font.pixelSize: 14
                                color: Theme.base
                                visible: picker.selectedIndices.has(modelData._index)
                            }
                        }
                    }

                    MouseArea {
                        id: mouseArea
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onEntered: picker.selected = index
                        onClicked: picker.toggleSelection(modelData)
                    }
                }
            }
        }
    }
}
