import QtQuick
import qs.core
import qs.globals
import qs.windows

PickerWindow {
    id: picker
    ipcTarget: "multipicker"
    cancelResult: []

    returnKeyHandler: function(entry) { picker.toggleSelection(entry) }
    escapeKeyHandler: function() { picker.confirmSelection() }

    extraHeader: Component {
        Row {
            width: parent.width
            height: 36
            spacing: 8

            Text {
                text: picker.selectedIndices.size + " selected"
                font.family: Theme.fontName; font.pixelSize: Theme.fontSize
                color: Theme.muted
                anchors.verticalCenter: parent.verticalCenter
            }

            Item { width: parent.width - 200; height: 1 }

            Text {
                text: "Return = toggle · Esc = confirm"
                font.family: Theme.fontName; font.pixelSize: Theme.fontSizeSmall
                color: Theme.muted
                anchors.verticalCenter: parent.verticalCenter
            }
        }
    }

    itemDelegate: Rectangle {
        required property var modelData
        required property int index
        width: ListView.view.width
        height: Dimensions.listItemHeight
        radius: Dimensions.inputRadius
        color: index === picker.selected
            ? Theme.primarySelected
            : mouseArea.containsMouse ? Theme.primaryHovered : "transparent"

        Rectangle {
            visible: index === picker.selected
            anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
            width: 3; radius: 1.5; color: Theme.primary
        }

        Row {
            anchors.fill: parent
            anchors.margins: 8
            spacing: 12

            ItemRow {
                anchors.verticalCenter: parent.verticalCenter
                icon: modelData.icon
                name: modelData.name
                comment: modelData.comment || ""
                selected: index === picker.selected
                textWidth: ListView.view.width - 110 - 32
                glyphSize: 22
            }

            Rectangle {
                width: 24; height: 24; radius: 4
                color: picker.selectedIndices.has(modelData._index)
                    ? Theme.primary : Theme.surface0
                border.color: Theme.surface1; border.width: 1
                anchors.verticalCenter: parent.verticalCenter

                Text {
                    anchors.centerIn: parent
                    text: "✓"; font.pixelSize: Theme.fontSizeSmall; color: Theme.base
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

    function toggleSelection(entry) {
        if (!entry) return
        var src = entry._source
        var idx = pickerItems.indexOf(src)
        if (idx === -1) return
        var s = new Set(picker.selectedIndices)
        if (s.has(idx)) s.delete(idx)
        else s.add(idx)
        picker.selectedIndices = s
    }

    function confirmSelection() {
        var items = []
        picker.selectedIndices.forEach(function(idx) {
            if (idx >= 0 && idx < pickerItems.length)
                items.push(pickerItems[idx])
        })
        writeResult(items)
        resetPicker()
        hide()
    }
}
