import QtQuick
import qs.core
import qs.globals
import qs.windows

PickerWindow {
    id: picker
    ipcTarget: "grouppicker"

    buildModel: function(items, query) {
        var all = items.map(function(e) {
            return {
                icon: e.icon, name: e.header, comment: e.description || "",
                group: e.group || "General", showHeader: false, _source: e, keywords: []
            }
        })
        all.sort(function(a, b) {
            if (a.group !== b.group) return a.group < b.group ? -1 : 1
            return (a.name || "") < (b.name || "") ? -1 : 1
        })
        var q = query ? query.toLowerCase() : ""
        var list = q ? all.filter(function(e) {
            return (e.name || "").toLowerCase().indexOf(q) !== -1
                || (e.comment || "").toLowerCase().indexOf(q) !== -1
        }) : all
        for (var i = 0; i < list.length; i++)
            list[i].showHeader = i === 0 || list[i].group !== list[i - 1].group
        return list
    }

    itemDelegate: Rectangle {
        required property var modelData
        required property int index
        width: ListView.view.width
        height: (modelData.showHeader ? 34 : 0) + Dimensions.listItemHeight
        color: "transparent"

        Rectangle {
            visible: modelData.showHeader
            width: parent.width; height: 30
            color: "transparent"
            anchors.top: parent.top

            Rectangle {
                anchors { left: parent.left; right: parent.right; bottom: parent.bottom }
                height: 1; color: Theme.surface1
            }

            Text {
                text: modelData.group
                font.family: Theme.fontName; font.pixelSize: Theme.fontSize
                font.capitalization: Font.AllUppercase; font.bold: true
                color: Theme.muted
                anchors.left: parent.left; anchors.leftMargin: 8
                anchors.verticalCenter: parent.verticalCenter
            }
        }

        Rectangle {
            width: parent.width; height: Dimensions.listItemHeight; radius: Dimensions.inputRadius
            y: modelData.showHeader ? 34 : 0
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
}
