import QtQuick
import QtQuick.Layouts
import Quickshell
import qs.globals
import qs.bar.widgets.interactable
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

    exclusiveZone: Dimensions.barHeight
    implicitHeight: Dimensions.barHeight
    color: "transparent"

    RowLayout {
        anchors.fill: parent
        anchors.leftMargin: 12
        anchors.rightMargin: 12
        spacing: 8

        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            CommandCenter { }
            Separator { }
            CpuMonitor { }
            Separator { }
            MemoryMonitor { }
            Separator { }
            DiskMonitor { }
        }

        Item { Layout.preferredWidth: 8 }
        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            WorkspacesInfo { }
            Separator { }
            ScratchpadInfo { }
        }

        Item { Layout.fillWidth: true }
        BarWidget {
            Layout.alignment: Qt.AlignVCenter
            TrayInfo { }
            Separator { }
            AudioInfo { }
            Separator { }
            ClockInfo { }
        }
    }

    Row {
        anchors.centerIn: parent
        spacing: 0

        Bracket { opening: true; anchors.verticalCenter: parent.verticalCenter }
        TitleInfo {
            width: Math.min(implicitWidth, (bar.width - 400))
        }
        Bracket { opening: false; anchors.verticalCenter: parent.verticalCenter }
        Item { width: 6; height: 1 }
        ResizeIndicator { anchors.verticalCenter: parent.verticalCenter }
    }
}
