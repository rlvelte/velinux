pragma Singleton
import QtQuick
import Quickshell

Singleton {
    readonly property color primary: "#fd7803"
    readonly property color primaryDim: "#ac5323"
    readonly property color primarySubtle: "#d46202"
    readonly property color primaryMuted: "#fcd4a8"

    readonly property color primarySelected: Qt.rgba(primary.r, primary.g, primary.b, 0.14)
    readonly property color primaryHovered:  Qt.rgba(primary.r, primary.g, primary.b, 0.07)

    readonly property color secondary: "#003562"
    readonly property color secondaryDim: "#021b2f"
    readonly property color secondaryLight: "#0a4a7a"

    readonly property color accent: "#00aaff"
    readonly property color accentDim: "#0077cc"
    readonly property color accentLight: "#66ccff"

    readonly property color base: "#010d18"
    readonly property color mantle: "#021b2f"
    readonly property color crust: "#032238"
    readonly property color surface0: "#003562"
    readonly property color surface1: "#0a4a7a"
    readonly property color surface2: "#1a5a8a"

    readonly property color text: "#e0eaf0"
    readonly property color subtext: "#a0b8c8"
    readonly property color overlay: "#608090"
    readonly property color muted: "#405060"

    readonly property color success: "#00a651"
    readonly property color warning: "#fd7803"
    readonly property color warningSubtle: "#fcd4a8"
    readonly property color error: "#ed1c24"
    readonly property color errorSubtle: "#f8d7da"
    readonly property color info: "#0099ff"
    readonly property color infoSubtle: "#b3e0ff"

    readonly property color textOnPrimary: "#010d18"
    readonly property color textOnSecondary: "#e0eaf0"
    readonly property color textOnAccent: "#000000"
    readonly property color textOnSurface: "#e0eaf0"

    readonly property string fontName: "Montserrat"
    readonly property string fontNameHeading: "Montserrat"
    readonly property string fontNameMono: "JetBrainsMono Nerd Font"
    readonly property int fontSize: 16
    readonly property int fontSizeSmall: 14
    readonly property int fontSizeLarge: 18
    readonly property int fontSizeHeading: 20

    readonly property string logoPath: "file:///home/rvelte/.config/vlx/themes/logos/vt.png"
}
