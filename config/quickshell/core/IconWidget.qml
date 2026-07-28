import QtQuick
import Quickshell
import qs.globals

Rectangle {
    id: root

    required property string icon
    required property string name

    property color bgColor: Theme.surface0
    property color fgColor: Theme.subtext
    property int glyphSize: 28
    property int imageSize: 28

    readonly property bool __isFile: icon && icon.startsWith("/")
    readonly property bool __isGlyph: !__isFile && !!icon && !/^[\w.-]+$/.test(icon)

    width: 48
    height: 48
    radius: 10
    color: bgColor

    Image {
        id: iconImg
        anchors.centerIn: parent
        source: root.__isFile ? root.icon : ""
        width: root.imageSize
        height: root.imageSize
        fillMode: Image.PreserveAspectFit
        visible: root.__isFile && status === Image.Ready
    }

    Text {
        anchors.centerIn: parent
        text: {
            if (root.__isFile) return ""
            if (!root.icon) return root.name.charAt(0).toUpperCase()
            if (root.__isGlyph) return root.icon
            return root.name.charAt(0).toUpperCase()
        }
        font.family: Theme.fontNameMono
        font.pixelSize: root.__isFile ? root.imageSize : root.glyphSize
        color: root.fgColor
        visible: !root.__isFile || iconImg.status !== Image.Ready
    }
}
