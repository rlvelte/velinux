import QtQuick
import QtQuick.Layouts
import qs.services

Item {
    id: bracketLeft
    implicitWidth: bracketLeftText.implicitWidth
    implicitHeight: 40
    Layout.alignment: Qt.AlignVCenter

    Text {
        id: bracketLeftText
        anchors.verticalCenter: parent.verticalCenter
        text: "["
        color: Theme.subtext
        font.family: Theme.fontName
        font.pixelSize: 18
        font.weight: Font.Bold
        verticalAlignment: Text.AlignVCenter
        rightPadding: 8
    }
}
