import QtQuick
import QtQuick.Layouts
import Quickshell
import qs.services
import qs.bar.widgets
import qs.bar.widgets.hardware
import qs.bar.widgets.design
import qs.bar.widgets.info

PanelWindow {
    id: bar

    anchors {
        top: true
        left: true
        right: true
    }

    margins {
        top: 8
        left: 8
        right: 8
    }

    exclusiveZone: 40
    implicitHeight: 40
    color: "transparent"

    Rectangle {
        anchors.fill: parent
        color: Theme.base
        radius: 8
        opacity: 0.0
    }

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: 12
        anchors.rightMargin: 12
        spacing: 8

        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            CommandCenter { }
            Seperator { }
            CpuMonitor { }
            Seperator { }
            MemoryMonitor { }
            Seperator { }
            DiskMonitor { }
        }

        Item { Layout.preferredWidth: 8 }
        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            Workspaces { }
        }

        Item { Layout.fillWidth: true }
        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            TitleInfo { }
        }

        Item { Layout.fillWidth: true }
        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            TrayInfo { }
            Seperator { }
            AudioInfo { }
            Seperator { }
            ClockInfo { }
        }
    }
}
