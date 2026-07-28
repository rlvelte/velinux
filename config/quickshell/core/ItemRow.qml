import QtQuick
import Quickshell
import qs.globals

Row {
    id: root

    required property string icon
    required property string name
    property string comment: ""
    property bool selected: false
    property color dimColor: Theme.subtext
    property real textWidth: 0
    property real glyphSize: 28

    spacing: 12

    IconWidget {
        anchors.verticalCenter: parent.verticalCenter
        icon: root.icon
        name: root.name
        glyphSize: root.glyphSize
    }

    Column {
        anchors.verticalCenter: parent.verticalCenter
        spacing: 2

        Text {
            text: root.name
            font.family: Theme.fontName
            font.pixelSize: Theme.fontSizeLarge
            color: root.selected ? Theme.text : root.dimColor
            elide: Text.ElideRight
            width: root.textWidth > 0 ? root.textWidth : undefined
        }

        Text {
            visible: !!root.comment
            text: root.comment
            font.family: Theme.fontName
            font.pixelSize: Theme.fontSize
            color: Theme.muted
            elide: Text.ElideRight
            width: root.textWidth > 0 ? root.textWidth : undefined
        }
    }
}
