import QtQuick
import Quickshell
import Quickshell.Io
import qs.services

Item {
    id: root
    property string filterText: ""

    property var bl1Data: []
    property var bl2Data: []
    property int activeLeague: 0

    function refresh() {
        bl1Process.running = true
        bl2Process.running = true
    }

    function moveUp() {
        var data = activeLeague === 0 ? bl1Data : bl2Data
        if (list.model && list.model.length > 0)
            list.currentIndex = Math.max(list.currentIndex - 1, 0)
    }

    function moveDown() {
        var data = activeLeague === 0 ? bl1Data : bl2Data
        if (list.model && list.model.length > 0)
            list.currentIndex = Math.min(list.currentIndex + 1, list.model.length - 1)
    }

    function activate() {}

    function formatValue(val, suffix) {
        if (val === undefined || val === null) return "--"
        return String(val) + (suffix || "")
    }

    Process {
        id: bl1Process
        command: ["vlx", "cc", "bl", "table", "--json", "--league", "bl1"]
        stdout: StdioCollector {
            onStreamFinished: {
                try { root.bl1Data = JSON.parse(text) } catch(e) { root.bl1Data = []; console.error("cc bl1 parse error:", e) }
            }
        }
    }

    Process {
        id: bl2Process
        command: ["vlx", "cc", "bl", "table", "--json", "--league", "bl2"]
        stdout: StdioCollector {
            onStreamFinished: {
                try { root.bl2Data = JSON.parse(text) } catch(e) { root.bl2Data = []; console.error("cc bl2 parse error:", e) }
            }
        }
    }

    Column {
        anchors.fill: parent
        spacing: 6

        Row {
            width: parent.width
            height: 24

            Repeater {
                model: ["1. Bundesliga", "2. Bundesliga"]

                Rectangle {
                    width: parent.width / 2
                    height: parent.height
                    color: index === root.activeLeague ? Theme.crust : "transparent"
                    radius: 4

                    Text {
                        anchors.centerIn: parent
                        text: modelData
                        font.family: Theme.fontName
                        font.pixelSize: Theme.fontSizeSmall
                        font.weight: index === root.activeLeague ? Font.Bold : Font.Medium
                        color: index === root.activeLeague ? Theme.text : Theme.subtext
                    }

                    MouseArea {
                        anchors.fill: parent
                        cursorShape: Qt.PointingHandCursor
                        onClicked: {
                            root.activeLeague = index
                            list.currentIndex = -1
                        }
                    }
                }
            }
        }

        Row {
            width: parent.width
            height: 24

            Text {
                width: 30
                height: parent.height
                text: "#"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                font.weight: Font.Bold
                color: Theme.muted
                verticalAlignment: Text.AlignVCenter
            }
            Text {
                width: parent.width - 30 - 40 - 60 - 80
                height: parent.height
                text: "Team"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                font.weight: Font.Bold
                color: Theme.muted
                verticalAlignment: Text.AlignVCenter
            }
            Text {
                width: 40
                height: parent.height
                text: "P"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                font.weight: Font.Bold
                color: Theme.muted
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignHCenter
            }
            Text {
                width: 60
                height: parent.height
                text: "G"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                font.weight: Font.Bold
                color: Theme.muted
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignHCenter
            }
            Text {
                width: 80
                height: parent.height
                text: "W/D/L"
                font.family: Theme.fontName
                font.pixelSize: Theme.fontSize
                font.weight: Font.Bold
                color: Theme.muted
                verticalAlignment: Text.AlignVCenter
                horizontalAlignment: Text.AlignHCenter
            }
        }

        Rectangle {
            width: parent.width
            height: 1
            color: Theme.surface1
        }

        ListView {
            id: list
            width: parent.width
            height: parent.height - 24 - 24 - 1 - 6 - 6
            clip: true
            spacing: 2
            highlightFollowsCurrentItem: true
            boundsBehavior: Flickable.StopAtBounds

            model: {
                var data = root.activeLeague === 0 ? root.bl1Data : root.bl2Data
                if (!data.length) return []
                if (!root.filterText) return data
                var f = root.filterText.toLowerCase()
                return data.filter(function(row) {
                    return (row.team || "").toLowerCase().indexOf(f) !== -1
                })
            }

            delegate: Rectangle {
                required property var modelData
                required property int index
                width: list.width
                height: 28
                radius: 4
                color: index === list.currentIndex
                    ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                    : blMouse.containsMouse
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                        : "transparent"

                Row {
                    anchors.fill: parent
                    anchors.margins: 4

                    Text {
                        width: 30
                        height: parent.height
                        text: root.formatValue(modelData.position)
                        font.family: Theme.fontNameMono
                        font.pixelSize: Theme.fontSizeSmall
                        color: index === list.currentIndex ? Theme.text : Theme.subtext
                        verticalAlignment: Text.AlignVCenter
                        horizontalAlignment: Text.AlignHCenter
                    }
                    Text {
                        width: parent.width - 30 - 40 - 60 - 80
                        height: parent.height
                        text: root.formatValue(modelData.team)
                        font.family: Theme.fontName
                        font.pixelSize: Theme.fontSize
                        color: index === list.currentIndex ? Theme.text : Theme.subtext
                        elide: Text.ElideRight
                        verticalAlignment: Text.AlignVCenter
                    }
                    Text {
                        width: 40
                        height: parent.height
                        text: root.formatValue(modelData.points)
                        font.family: Theme.fontNameMono
                        font.pixelSize: Theme.fontSizeSmall
                        color: Theme.primaryMuted
                        verticalAlignment: Text.AlignVCenter
                        horizontalAlignment: Text.AlignHCenter
                    }
                    Text {
                        width: 60
                        height: parent.height
                        text: root.formatValue(modelData.goals_for) + ":" + root.formatValue(modelData.goals_against)
                        font.family: Theme.fontNameMono
                        font.pixelSize: Theme.fontSize
                        color: Theme.subtext
                        verticalAlignment: Text.AlignVCenter
                        horizontalAlignment: Text.AlignHCenter
                    }
                    Text {
                        width: 80
                        height: parent.height
                        text: root.formatValue(modelData.wins) + "/" + root.formatValue(modelData.draws) + "/" + root.formatValue(modelData.losses)
                        font.family: Theme.fontNameMono
                        font.pixelSize: Theme.fontSize
                        color: Theme.muted
                        verticalAlignment: Text.AlignVCenter
                        horizontalAlignment: Text.AlignHCenter
                    }
                }

                MouseArea {
                    id: blMouse
                    anchors.fill: parent
                    hoverEnabled: true
                    onEntered: list.currentIndex = index
                }
            }
        }
    }
}
