import QtQuick
import Quickshell
import qs.services

Item {
    id: clockWidget
    implicitWidth: clockRow.implicitWidth
    implicitHeight: 40

    SystemClock {
        id: clock
        precision: SystemClock.Seconds
    }

    Row {
        id: clockRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 6

        Text {
            text: Qt.formatDate(clock.date, "dd.MM.yyyy")
            color: Theme.subtext
            font.family: Theme.fontName
            font.pixelSize: 18
            font.weight: Font.Medium
        }

        Text {
            text: "|"
            color: Theme.surface1
            font.family: Theme.fontName
            font.pixelSize: 18
            rightPadding: 4
            leftPadding: 4
        }

        Text {
            text: Qt.formatTime(clock.date, "HH:mm")
            color: Theme.primary
            font.family: Theme.fontName
            font.pixelSize: 18
            font.weight: Font.Bold
        }
    }
}
