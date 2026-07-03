import QtQuick
import Quickshell
import Quickshell.Io
import qs.services

Item {
    id: root
    property string filterText: ""

    property var feedsData: []
    property string currentFeedCmd: "hackernews"

    property var feedSources: [
        { name: "Tagesschau", cmd: "tagesschau" },
        { name: "openSUSE News", cmd: "openSUSE News" }
    ]

    function refresh() {
        feedsProcess.running = true
    }

    function moveUp() {
        if (feedsList.model && feedsList.model.length > 0)
            feedsList.currentIndex = Math.max(feedsList.currentIndex - 1, 0)
    }

    function moveDown() {
        if (feedsList.model && feedsList.model.length > 0)
            feedsList.currentIndex = Math.min(feedsList.currentIndex + 1, feedsList.model.length - 1)
    }

    function activate() {
        if (feedsList.model && feedsList.model.length > feedsList.currentIndex) {
            var item = feedsList.model[feedsList.currentIndex]
            if (item && item.link) {
                openUrlProcess.command = ["xdg-open", item.link]
                openUrlProcess.running = true
            }
        }
    }

    function formatValue(val, suffix) {
        if (val === undefined || val === null) return "--"
        return String(val) + (suffix || "")
    }

    Process {
        id: feedsProcess
        command: ["vlx", "cc", "feeds", "poll", root.currentFeedCmd, "--json"]
        stdout: StdioCollector {
            onStreamFinished: {
                try { root.feedsData = JSON.parse(text) } catch(e) { root.feedsData = []; console.error("cc feeds parse error:", e) }
            }
        }
    }

    Process {
        id: openUrlProcess
    }

    Column {
        anchors.fill: parent
        spacing: 8

        Flickable {
            width: parent.width
            height: 28
            contentWidth: sourcesRow.width
            clip: true
            boundsBehavior: Flickable.StopAtBounds

            Row {
                id: sourcesRow
                height: 28
                spacing: 4

                Repeater {
                    model: root.feedSources

                    delegate: Rectangle {
                        required property var modelData
                        width: srcLabel.implicitWidth + 16
                        height: 28
                        radius: 14
                        color: modelData.cmd === root.currentFeedCmd
                            ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.2)
                            : srcMouse.containsMouse
                                ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                                : Theme.surface0
                        border.color: modelData.cmd === root.currentFeedCmd ? Theme.primary : "transparent"
                        border.width: 1

                        Text {
                            id: srcLabel
                            anchors.centerIn: parent
                            text: modelData.name
                            font.family: Theme.fontName
                            font.pixelSize: Theme.fontSize
                            color: modelData.cmd === root.currentFeedCmd ? Theme.primaryMuted : Theme.subtext
                        }

                        MouseArea {
                            id: srcMouse
                            anchors.fill: parent
                            hoverEnabled: true
                            cursorShape: Qt.PointingHandCursor
                            onClicked: {
                                root.currentFeedCmd = modelData.cmd
                                feedsProcess.running = true
                            }
                        }
                    }
                }
            }
        }

        ListView {
            id: feedsList
            width: parent.width
            height: parent.height - 28 - 8
            clip: true
            spacing: 4
            highlightFollowsCurrentItem: true
            boundsBehavior: Flickable.StopAtBounds

            model: {
                if (!root.feedsData.length || !root.feedsData[0].items) return []
                var items = root.feedsData[0].items
                if (!root.filterText) return items
                var f = root.filterText.toLowerCase()
                return items.filter(function(item) {
                    return (item.title || "").toLowerCase().indexOf(f) !== -1
                        || (item.published || "").toLowerCase().indexOf(f) !== -1
                })
            }

            delegate: Rectangle {
                required property var modelData
                required property int index
                width: feedsList.width
                height: 64
                radius: 8
                color: index === feedsList.currentIndex
                    ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.14)
                    : feedMouse.containsMouse
                        ? Qt.rgba(Theme.primary.r, Theme.primary.g, Theme.primary.b, 0.07)
                        : "transparent"

                Rectangle {
                    visible: index === feedsList.currentIndex
                    anchors { left: parent.left; top: parent.top; bottom: parent.bottom }
                    width: 3
                    radius: 1.5
                    color: Theme.primary
                }

                Column {
                    anchors.fill: parent
                    anchors.margins: 8
                    spacing: 2

                    Text {
                        text: modelData.published || ""
                        font.family: Theme.fontNameMono
                        font.pixelSize: Theme.fontSize
                        color: Theme.muted
                    }
                    Text {
                        text: modelData.title || ""
                        font.family: Theme.fontName
                        font.pixelSize: Theme.fontSizeLarge
                        color: index === feedsList.currentIndex ? Theme.text : Theme.subtext
                        elide: Text.ElideRight
                        width: parent.width
                    }
                }

                MouseArea {
                    id: feedMouse
                    anchors.fill: parent
                    hoverEnabled: true
                    cursorShape: Qt.PointingHandCursor
                    onEntered: feedsList.currentIndex = index
                    onClicked: {
                        if (modelData.link) {
                            openUrlProcess.command = ["xdg-open", modelData.link]
                            openUrlProcess.running = true
                        }
                    }
                }
            }
        }
    }
}
