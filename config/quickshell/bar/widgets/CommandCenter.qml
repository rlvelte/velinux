import QtQuick
import Quickshell
import qs.services

Item {
    id: commandCenter
    implicitWidth: commandCenterIcon.implicitWidth + 8
    implicitHeight: 40

    Text {
        id: commandCenterIcon
        anchors.centerIn: parent
        text: {
            for (let i = 0; i < Vlx.themes.length; i++) {
                if (Vlx.themes[i].active) return Vlx.themes[i].icon
            }
            return ""
        }
        color: Theme.primary
        font.family: Theme.fontName
        font.pixelSize: 22
        verticalAlignment: Text.AlignVCenter
    }

    MouseArea {
        anchors.fill: parent
        cursorShape: Qt.PointingHandCursor
        onClicked: popup.visible = !popup.visible
    }

    FloatingWindow {
        id: popup
        visible: false
        color: "transparent"
        implicitWidth: 300
        implicitHeight: popupContent.implicitHeight + 32

        Rectangle {
            anchors.fill: parent
            anchors.margins: 8
            color: Theme.base
            radius: 12
            border.color: Theme.surface1
            border.width: 1
            focus: true
            Keys.onEscapePressed: popup.visible = false

            Column {
                id: popupContent
                anchors.fill: parent
                anchors.margins: 16
                spacing: 4

                Text {
                    text: "Themes"
                    font.family: Theme.fontName
                    font.pixelSize: 18
                    font.weight: Font.Bold
                    color: Theme.text
                }

                Repeater {
                    model: Vlx.themes

                    delegate: Rectangle {
                        width: popupContent.width
                        height: 48
                        radius: 8
                        color: mouseArea.containsMouse ? Theme.surface1 : (modelData.active ? Theme.surface0 : "transparent")

                        Row {
                            anchors.fill: parent
                            anchors.margins: 8
                            spacing: 12

                            Text {
                                text: modelData.icon
                                font.family: Theme.fontName
                                font.pixelSize: 18
                                color: Theme.primary
                                width: 32
                                horizontalAlignment: Text.AlignHCenter
                            }

                            Text {
                                text: modelData.name
                                font.family: Theme.fontName
                                font.pixelSize: 16
                                color: modelData.active ? Theme.primary : Theme.text
                                font.weight: modelData.active ? Font.Bold : Font.Normal
                                anchors.verticalCenter: parent.verticalCenter
                            }
                        }

                        MouseArea {
                            id: mouseArea
                            anchors.fill: parent
                            hoverEnabled: true
                            cursorShape: Qt.PointingHandCursor
                            onClicked: {
                                Vlx.applyTheme(modelData.id)
                                popup.visible = false
                            }
                        }
                    }
                }
            }
        }
    }
}
