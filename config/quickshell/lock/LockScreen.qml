import QtQuick
import QtQuick.Layouts
import QtQuick.Controls.Fusion
import Quickshell.Wayland
import qs.services

Rectangle {
    id: root
    required property var context

    color: Theme.base

    Image {
        id: bgImage
        anchors.fill: parent
        source: Quickshell.env("HOME") + "/.config/vlx/themes/current.png"
        fillMode: Image.PreserveAspectCrop
        opacity: 0.15
    }

    Rectangle {
        anchors.fill: parent
        color: Qt.rgba(Theme.base.r, Theme.base.g, Theme.base.b, 0.85)
    }

    Label {
        id: clock
        property var date: new Date()
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.top
        anchors.topMargin: 120
        renderType: Text.NativeRendering
        font.family: Theme.fontName
        font.pointSize: 72
        color: Theme.text

        Timer {
            running: true
            repeat: true
            interval: 1000
            onTriggered: clock.date = new Date();
        }

        text: {
            const hours = clock.date.getHours().toString().padStart(2, '0');
            const minutes = clock.date.getMinutes().toString().padStart(2, '0');
            return `${hours}:${minutes}`;
        }
    }

    ColumnLayout {
        anchors.horizontalCenter: parent.horizontalCenter
        anchors.top: parent.verticalCenter
        anchors.topMargin: 40
        spacing: 16

        Text {
            text: Quickshell.env("USER")
            font.family: Theme.fontName
            font.pointSize: Theme.fontSizeHeading
            color: Theme.text
            anchors.horizontalCenter: parent.horizontalCenter
        }

        RowLayout {
            anchors.horizontalCenter: parent.horizontalCenter
            spacing: 8

            TextField {
                id: passwordBox
                implicitWidth: 300
                height: 48
                padding: 12
                focus: true
                enabled: !root.context.unlockInProgress
                Component.onCompleted: passwordBox.forceActiveFocus()
                echoMode: TextInput.Password
                inputMethodHints: Qt.ImhSensitiveData
                placeholderText: "Password"
                font.family: Theme.fontName
                font.pointSize: Theme.fontSize
                color: Theme.text
                placeholderTextColor: Theme.muted

                background: Rectangle {
                    color: Theme.surface0
                    radius: 8
                    border.color: passwordBox.activeFocus ? Theme.primary : Theme.surface1
                    border.width: passwordBox.activeFocus ? 2 : 1
                }

                onTextChanged: root.context.currentText = this.text;
                onAccepted: root.context.tryUnlock();

                Connections {
                    target: root.context
                    function onCurrentTextChanged() {
                        passwordBox.text = root.context.currentText;
                    }
                }
            }

            Button {
                text: "Unlock"
                height: 48
                padding: 12
                focusPolicy: Qt.NoFocus
                enabled: !root.context.unlockInProgress && root.context.currentText !== "";
                onClicked: root.context.tryUnlock();
            }
        }

        Text {
            visible: root.context.showFailure
            text: "Incorrect password"
            font.family: Theme.fontName
            font.pointSize: Theme.fontSize
            color: Theme.error
            anchors.horizontalCenter: parent.horizontalCenter
        }
    }
}
