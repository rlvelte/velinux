import QtQuick
import qs.core
import qs.globals
import qs.windows

PickerWindow {
    id: picker
    ipcTarget: "singlepicker"

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
