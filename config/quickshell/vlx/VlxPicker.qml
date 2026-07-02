import QtQuick
import QtQuick.Controls
import Quickshell
import Quickshell.Io
import qs.services

FloatingWindow {
    id: picker
    title: "vlx-picker"
    color: "transparent"
    implicitWidth: 520
    implicitHeight: 420

    property bool shown: false
    property var items: []
    property string filter: ""
    property bool multiMode: false
    property var selectedItems: []
    property string resultFile: ""
    property int selected: 0

    visible: shown

    function open(filePath, resultPath, multi) {
        items = []
        filter = ""
        multiMode = multi || false
        selectedItems = []
        resultFile = resultPath || ""
        searchField.text = ""
        selected = 0
        reader.command = ["cat", filePath]
        reader.running = true
    }

    function writeResult() {
        if (!resultFile) return
        resultWriter.path = resultFile
        if (multiMode) {
            resultWriter.setText(selectedItems.join("\n") + "\n")
        } else {
            resultWriter.setText(selectedItems[0] + "\n")
        }
    }

    function cancel() {
        shown = false
        if (resultFile) {
            resultWriter.path = resultFile
            resultWriter.setText("")
        }
    }

    property var filteredItems: {
        if (!filter) return items
        var f = filter.toLowerCase()
        var out = []
        for (var i = 0; i < items.length; i++) {
            if (items[i].toLowerCase().indexOf(f) !== -1)
                out.push(items[i])
        }
        return out
    }

    onShownChanged: {
        if (shown) {
            Qt.callLater(function() { searchField.forceActiveFocus() })
        }
    }

    Process {
        id: reader
        stdout: StdioCollector {
            onStreamFinished: {
                var lines = text.trim().split("\n")
                picker.items = lines
                picker.shown = true
                Qt.callLater(function() { searchField.forceActiveFocus() })
            }
        }
    }

    FileView {
        id: resultWriter
        blockWrites: true
        atomicWrites: true
    }

    Rectangle {
        anchors.fill: parent
        anchors.margins: 16
        color: Theme.base
        radius: 12
        border.color: Theme.surface1
        border.width: 1
        focus: true

        Keys.onEscapePressed: picker.cancel()

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
                font.pixelSize: 22
                color: Theme.text
                placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0
                    radius: 8
                    border.color: searchField.activeFocus ? Theme.primary : Theme.surface1
                    border.width: searchField.activeFocus ? 2 : 1
                }

                onTextChanged: {
                    picker.filter = text
                    if (picker.selected >= filteredItems.length)
                        picker.selected = Math.max(0, filteredItems.length - 1)
                }
                onAccepted: {
                    if (filteredItems.length > 0)
                        picker.selectItem(filteredItems[selected], selected)
                }
                Keys.onDownPressed: {
                    selected = Math.min(selected + 1, filteredItems.length - 1)
                }
                Keys.onUpPressed: {
                    selected = Math.max(selected - 1, 0)
                }
                Keys.onEscapePressed: {
                    if (searchField.length === 0)
                        picker.cancel()
                    else
                        searchField.text = ""
                }
            }

            ListView {
                id: itemList
                width: parent.width
                height: parent.height - searchField.height - 12 - (multiMode ? 56 : 0)
                clip: true
                model: picker.filteredItems
                spacing: 4
                currentIndex: picker.selected
                highlightFollowsCurrentItem: true
                boundsBehavior: Flickable.StopAtBounds

                delegate: Rectangle {
                    required property var modelData
                    required property int index
                    width: itemList.width
                    height: 48
                    radius: 8
                    color: {
                        if (multiMode && picker.selectedItems.indexOf(modelData) !== -1)
                            return Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                        if (index === picker.selected)
                            return Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                        return mouseArea.containsMouse
                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.04)
                            : "transparent"
                    }

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

                        Text {
                            text: modelData
                            font.family: Theme.fontName
                            font.pixelSize: 20
                            color: Theme.text
                            anchors.verticalCenter: parent.verticalCenter
                            elide: Text.ElideRight
                            width: parent.width - (multiMode ? 40 : 8)
                        }

                        Text {
                            visible: multiMode && picker.selectedItems.indexOf(modelData) !== -1
                            text: "\u2713"
                            font.family: Theme.fontName
                            font.pixelSize: Theme.fontSizeLarge
                            color: Theme.primary
                            anchors.verticalCenter: parent.verticalCenter
                            width: 32
                            horizontalAlignment: Text.AlignHCenter
                        }
                    }

                    MouseArea {
                        id: mouseArea
                        anchors.fill: parent
                        hoverEnabled: true
                        cursorShape: Qt.PointingHandCursor
                        onEntered: picker.selected = index
                        onClicked: {
                            picker.selected = index
                            picker.selectItem(modelData, index)
                        }
                    }
                }
            }

            Rectangle {
                visible: multiMode
                width: parent.width
                height: 48
                radius: 8
                color: Theme.surface0

                Text {
                    text: "Done (" + picker.selectedItems.length + " selected)"
                    font.family: Theme.fontName
                    font.pixelSize: Theme.fontSize
                    color: Theme.primary
                    anchors.centerIn: parent
                }

                MouseArea {
                    anchors.fill: parent
                    cursorShape: Qt.PointingHandCursor
                    onClicked: {
                        picker.shown = false
                        picker.writeResult()
                    }
                }
            }
        }
    }

    function selectItem(item, index) {
        if (multiMode) {
            var idx = selectedItems.indexOf(item)
            if (idx === -1) {
                var arr = selectedItems.slice()
                arr.push(item)
                selectedItems = arr
            } else {
                var arr = selectedItems.slice()
                arr.splice(idx, 1)
                selectedItems = arr
            }
        } else {
            selectedItems = [item]
            writeResult()
            shown = false
        }
    }
}
