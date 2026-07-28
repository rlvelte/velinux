import QtQuick
import QtQuick.Layouts
import qs.globals

Item {
    id: root

    property bool opening: true

    implicitWidth: bracketText.implicitWidth
    implicitHeight: Dimensions.barHeight
    Layout.alignment: Qt.AlignVCenter

    Text {
        id: bracketText
        anchors.verticalCenter: parent.verticalCenter
        text: root.opening ? "[" : "]"
        color: Theme.subtext
        font.family: Theme.fontName
        font.pixelSize: Theme.fontSizeHeading
        font.weight: Font.Medium
        verticalAlignment: Text.AlignVCenter
        rightPadding: root.opening ? 8 : 0
        leftPadding: root.opening ? 0 : 8
    }
}
