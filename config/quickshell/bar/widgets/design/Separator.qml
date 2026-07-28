import QtQuick
import qs.globals

Item {
    id: root
    implicitWidth: separatorText.implicitWidth
    implicitHeight: Dimensions.barHeight

    Text {
        id: separatorText
        anchors.verticalCenter: parent.verticalCenter
        text: "|"
        color: Theme.surface1
        font.family: Theme.fontName
        font.pixelSize: Theme.fontSizeHeading
        rightPadding: 4
        leftPadding: 4
        verticalAlignment: Text.AlignVCenter
    }
}
