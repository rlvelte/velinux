import QtQuick
import qs.services

Item {
    id: seperator
    implicitWidth: seperatorText.implicitWidth
    implicitHeight: 40

    Text {
        id: seperatorText
        anchors.verticalCenter: parent.verticalCenter
        text: "|"
        color: Theme.surface1
        font.family: Theme.fontName
        font.pixelSize: 18
        rightPadding: 4
        leftPadding: 4
        verticalAlignment: Text.AlignVCenter
    }
}
