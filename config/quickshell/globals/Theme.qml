pragma Singleton
import QtQuick
import Quickshell

Singleton {
    readonly property color primary: "#fd7803"

    readonly property color base: "#010d18"
    readonly property color surface0: "#003562"
    readonly property color surface1: "#0a4a7a"

    readonly property color text: "#e0eaf0"
    readonly property color subtext: "#a0b8c8"
    readonly property color muted: "#405060"

    readonly property color success: "#00a651"
    readonly property color warning: "#fd7803"
    readonly property color error: "#ed1c24"

    readonly property color primarySelected: Qt.rgba(primary.r, primary.g, primary.b, 0.14)
    readonly property color primaryHovered:  Qt.rgba(primary.r, primary.g, primary.b, 0.07)
    readonly property color textOnPrimary: "#010d18"

    readonly property string fontName: "Montserrat"
    readonly property string fontNameMono: "JetBrainsMono Nerd Font"
    readonly property int fontSize: 16
    readonly property int fontSizeSmall: 14
    readonly property int fontSizeLarge: 18
    readonly property int fontSizeHeading: 20

    readonly property string logoPath: "file://" + Quickshell.env("HOME") + "/.config/vlx/themes/logos/vt.png"
}
