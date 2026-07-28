import QtQuick
import qs.core

PickerWindow {
    id: picker
    ipcTarget: "singlepicker"

    itemDelegate: Rectangle {
        required property var modelData
        required property int index
        width: ListView.view.width
        height: Theme.listItemHeight
        radius: Theme.inputRadius
        color: mouseArea.containsMouse ? Theme.primaryHovered : "transparent"

        ItemRow {
            anchors.fill: parent
            anchors.margins: 8
            icon: modelData.icon
            name: modelData.name
            comment: modelData.comment || ""
            selected: index === picker.selected
            textWidth: ListView.view.width - 110
        }

        MouseArea {
            id: mouseArea
            anchors.fill: parent
            hoverEnabled: true
            cursorShape: Qt.PointingHandCursor
            onEntered: picker.selected = index
            onClicked: picker.launch(modelData)
        }
    }
}
