import QtQuick
import Quickshell.Wayland
import qs.services

Item {
    id: titleInfo
    implicitWidth: titleInfoText.implicitWidth
    implicitHeight: 40

    Text {
        id: titleInfoText
        anchors.verticalCenter: parent.verticalCenter
        text: ToplevelManager.activeToplevel ? ToplevelManager.activeToplevel.title.substring(0, 80) : "..."
        color: Theme.text
        font.family: Theme.fontName
        font.pixelSize: 18
        font.weight: Font.Medium
        elide: Text.ElideRight
        verticalAlignment: Text.AlignVCenter
        maximumLineCount: 1
    }
}
