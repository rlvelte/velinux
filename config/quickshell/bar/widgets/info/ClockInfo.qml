import QtQuick
import Quickshell
import Quickshell.Io
import qs.globals
import qs.bar.widgets.design

Item {
    id: clockWidget
    implicitWidth: clockRow.implicitWidth
    implicitHeight: 40

    SystemClock {
        id: clock
        precision: SystemClock.Seconds
    }

    Process {
        id: calendarProcess
        command: ["thunderbird"]
    }

    Row {
        id: clockRow
        anchors.verticalCenter: parent.verticalCenter
        spacing: 6

        Text {
            text: Qt.formatDate(clock.date, "dd.MM.yyyy")
            color: Theme.subtext
            font.family: Theme.fontName
            font.pixelSize: Theme.fontSizeHeading
            font.weight: Font.Medium
        }

        Separator { }

        Text {
            text: Qt.formatTime(clock.date, "HH:mm")
            color: Theme.primary
            font.family: Theme.fontName
            font.pixelSize: Theme.fontSizeHeading
            font.weight: Font.Bold
        }
    }

    MouseArea {
        anchors.fill: parent
        cursorShape: Qt.PointingHandCursor
        onClicked: {
            calendarProcess.running = true
        }
    }
}
